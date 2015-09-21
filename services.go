package webhook

import (
    "encoding/json"
    "net/http"
    "strings"
)

// Return type of hook.
func getHookType(request *http.Request) string {
    if request.Header.Get("X-Github-Event") != "" {
        return "github"
    } else if request.Header.Get("X-Sender") == "Doorbell" {
        return "doorbell"
    } else if strings.Index(request.Header.Get("User-Agent"), "Bitbucket") > -1 {
        return "bitbucket"
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
