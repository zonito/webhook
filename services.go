package webhook

import (
  "encoding/json"
)


/***
  * Return type of hook.
  */
func getHookType(request *http.Request) string {
  if request.Header.Get("X-Github-Event") != "" {
    return "github"
  } else if (request.Header.Get("X-Sender") == "Doorbell") {
    return "doorbell"
  }
  return ""
}


/***
  * Prepare and return description for service.
  */
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


/***
  * Return github data.
  */
func getGithubData(decoder *json.Decoder, header string) (string, string) {
  var gEvent GitPayload
  decoder.Decode(&gEvent)
  event := gEvent.Repository.Name + " --> " + header + " event"
  desc := getGitDescription(gEvent)
  return event, desc
}


/***
  * Return doorbell data.
  */
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
