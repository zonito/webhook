package services

import (
    "encoding/json"
    "strconv"
)

type JJNScm struct {
    Url    string
    Branch string
    Commit string
}

type JJNBuild struct {
    Full_url string
    Number   int
    Queue_id int
    Phase    string
    Status   string
    Url      string
    Scm      JJNScm
}

type JJNMessage struct {
    Name  string
    Url   string
    Build JJNBuild
    Log   string
}

func getJenkinsJobNoficationData(decoder *json.Decoder) (string, string) {
    var jjnEvent JJNMessage
    decoder.Decode(&jjnEvent)
    event := "Jenkins Job Notifier: " + jjnEvent.Name + ", Phase: " + jjnEvent.Build.Phase
    if jjnEvent.Build.Phase != "STARTED" {
        event += " (" + jjnEvent.Build.Status + ")"
    }
    desc := "URL: " + jjnEvent.Build.Full_url +
        "\nBuild Number: " + strconv.Itoa(jjnEvent.Build.Number)
    return event, desc
}
