package services

import (
    "appengine"
    "appengine/urlfetch"
    "encoding/json"
    "io/ioutil"
    "math/rand"
    "net/http"
    "strconv"
    "strings"
    "time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Return url response
func getResponse(context appengine.Context, url string) string {
    client := urlfetch.Client(context)
    resp, err := client.Get(url)
    if err != nil {
        context.Infof("GetBoards client.Get: %v", err.Error())
        return ""
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        context.Infof("GetBoards ioutil.ReadAll: %v", err.Error())
        return ""
    }
    return string(body)
}

// Return event type and description to post.
func GetEventData(request *http.Request) (string, string) {
    hookType := getHookType(request)
    var decoder *json.Decoder
    if hookType != "travis" {
        decoder = json.NewDecoder(request.Body)
    } else {
        payload := request.FormValue("payload")
        decoder = json.NewDecoder(strings.NewReader(payload))
    }
    switch hookType {
    case "github":
        return getGithubData(
            decoder, request.Header.Get("X-Github-Event"))
    case "doorbell":
        return getDoorbellData(decoder)
    case "bitbucket":
        return getBitbucketData(decoder, request.Header.Get("X-Event-Key"))
    case "travis":
        return getTravisData(decoder)
    }
    return "", ""
}

// Return type of hook.
func getHookType(request *http.Request) string {
    if request.Header.Get("X-Github-Event") != "" {
        return "github"
    } else if request.Header.Get("X-Sender") == "Doorbell" {
        return "doorbell"
    } else if strings.Index(request.Header.Get("User-Agent"), "Bitbucket") > -1 {
        return "bitbucket"
    } else if request.Header.Get("Travis-Repo-Slug") != "" {
        return "travis"
    }
    return ""
}

// Prepare and return description for service.
func getGitDescription(payload GitPayload) string {
    repo := payload.Repository
    desc := repo.Name + "\n===========" +
        "\n**Name**: " + repo.Name +
        "\n**Url**: " + repo.Url +
        "\n**Owner**: " + repo.Owner.Email +
        "\n**Compare**: " + payload.Compare +
        "\n**Ref**: " + payload.Ref +
        "\n Modified files\n------------\n"
    for i := 0; i < len(payload.Commits); i++ {
        commit := payload.Commits[i]
        desc += "\n* " + commit.Message + " (" + commit.Timestamp + ")"
        for j := 0; j < len(commit.Modified); j++ {
            desc += "\n * " + commit.Modified[j]
        }
    }
    return desc
}

// Return github data.
func getGithubData(decoder *json.Decoder, header string) (string, string) {
    var gEvent GitPayload
    decoder.Decode(&gEvent)
    event := gEvent.Repository.Name + " --> " + header + " event"
    desc := getGitDescription(gEvent)
    return event, desc
}

// Return doorbell data.
func getDoorbellData(decoder *json.Decoder) (string, string) {
    var dEvent DBPayload
    decoder.Decode(&dEvent)
    data := dEvent.Data
    event := data.Application.Name + " --> " +
        data.Sentiment + " feedback - from " + data.Email
    desc := data.Message + "\n\n **User Agent**: " +
        data.User_Agent + "\n\n **Reply**: " + data.Url
    return event, desc
}

// Return bitbucket data.
func getBitbucketData(decoder *json.Decoder, eType string) (string, string) {
    var bEvent BBPayload
    decoder.Decode(&bEvent)
    action := strings.Split(eType, ":")
    event, desc := "", ""
    who := bEvent.Actor.Username + " (" + bEvent.Actor.Display_name + ")"
    if action[0] == "repo" {
        switch action[1] {
        case "push":
            event = bEvent.Repository.Name + ": Push Event"
            if len(bEvent.Push.Changes) > 0 {
                change := bEvent.Push.Changes[0]
                desc = "Who: " + who + "\nCommits\n-------"
                for i := 0; i < len(change.Commits); i++ {
                    desc += "\n* " + change.Commits[i].Message +
                        " (" + change.Commits[i].Hash + ")" +
                        "\n * " + change.Commits[i].Author.Raw
                }
            }
        case "fork":
            event = bEvent.Repository.Name + ": Fork Event"
            desc = "\n" + who + " Forked."
        case "commit_comment_created":
            event = bEvent.Repository.Name + ": Commit Comment Created"
            desc = who + " commented on " + bEvent.Commit.Hash
            desc += "\n Comment: " + bEvent.Comment.Content.Markup
            desc += "\n File: " + bEvent.Comment.Inline.Path
        }
    } else if action[0] == "pullrequest" {
        desc = "Description: " + bEvent.Pullrequest.Description
        desc += "\n From Repository: " + bEvent.Pullrequest.Source.Repository.Name
        switch action[1] {
        case "created":
            desc += "\n Created by: " + who
        case "updated":
            desc += "\n Updated by: " + who
        case "approved":
            desc += "\n Approved by: " + bEvent.Approval.User.Username
        case "unapproved":
            desc += "\n Unapproved by: " + bEvent.Approval.User.Username
        case "fulfilled":
            desc += "\n Merged by: " + who
            desc += "\n **Merged**"
        case "rejected":
            desc += "\n Rejected by: " + who
            desc += "\n **Rejected** because " + bEvent.Pullrequest.Reason
        case "comment_created":
            desc += "\n Commented by: " + who
            desc += "\n Comment: " + bEvent.Comment.Content.Markup
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        case "comment_updated":
            desc += "\n Comment updated by: " + who
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        case "comment_deleted":
            desc += "\n Comment deleted by: " + who
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        }
        event = bEvent.Repository.Name + ": Pull request " + action[1] + ": " +
            bEvent.Pullrequest.Title + "(" + bEvent.Pullrequest.State + ")"
    }
    return event, desc
}

// Return travis data.
func getTravisData(decoder *json.Decoder) (string, string) {
    var tEvent TRPayload
    decoder.Decode(&tEvent)
    if tEvent.Id > 0 {
        event := "Travis: " + tEvent.Status_message + " for " +
            tEvent.Repository.Name
        desc := "**Status**: " + tEvent.Status_message +
            "\n **Duration**: " + strconv.Itoa(tEvent.Duration) +
            "\n **Message**: " + tEvent.Message +
            "\n **Build Number**: " + tEvent.Number +
            "\n **Type**: " + tEvent.Type +
            "\n **Compare URL**: " + tEvent.Compare_url +
            "\n **Committer Name**: " + tEvent.Committer_name +
            "\n **Build Url**: " + tEvent.Build_url
        return event, desc
    }
    return "", ""
}

// Return random alphanumeric string
func GetAlphaNumberic(n int) string {
    var src = rand.NewSource(time.Now().UnixNano())
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return string(b)
}
