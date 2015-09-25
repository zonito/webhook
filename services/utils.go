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
        return getBitbucketData(decoder)
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
func getBitbucketData(decoder *json.Decoder) (string, string) {
    var bEvent BBPayload
    decoder.Decode(&bEvent)
    if &bEvent.Push != nil {
        event := bEvent.Repository.Name + ": Push Event "
        var desc string
        if len(bEvent.Push.Changes) > 0 {
            change := bEvent.Push.Changes[0]
            desc = "\nCommits\n-------"
            for i := 0; i < len(change.Commits); i++ {
                desc += "\n\n* " + change.Commits[i].Message +
                    " (" + change.Commits[i].Hash + ")" +
                    "\n * " + change.Commits[i].Author.Raw
            }
        }
        return event, desc
    }
    return "", ""
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
