package services

import (
    "appengine"
    "encoding/json"
    "net/http"
)

type SDIncident struct {
    Incident_id    string
    Resource_id    string
    Resource_name  string
    State          string
    Started_at     int
    Ended_at       int
    Policy_name    string
    Condition_name string
    Url            string
    Summary        string
}

type SDMessage struct {
    Incident SDIncident
    Version  int
}

func getStackDriverData(decoder *json.Decoder, request *http.Request) (string, string) {
    var sdnEvent SDMessage
    decoder.Decode(&sdnEvent)
    context := appengine.NewContext(request)
    context.Infof("%s", decoder)
    context.Infof("%s", sdnEvent.Incident.Policy_name)
    event := "`StackDriver: " + sdnEvent.Incident.Policy_name +
        ", Condition: " + sdnEvent.Incident.Condition_name + "`"
    desc := "URL: " + sdnEvent.Incident.Url +
        "\nSummary: " + sdnEvent.Incident.Summary +
        "\nState: " + sdnEvent.Incident.State +
        "\nResource: " + sdnEvent.Incident.Resource_name
    return event, desc
}
