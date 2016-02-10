package services

import (
    "encoding/json"
)

type SDInstance struct {
    Incident_id    string
    Resource_id    string
    Resource_name  string
    State          string
    Started_at     string
    Ended_at       string
    Policy_name    string
    Condition_name string
    Url            string
    Summary        string
}

type SDMessage struct {
    Instance SDInstance
    Version  string
}

func getStackDriverData(decoder *json.Decoder) (string, string) {
    var sdnEvent SDMessage
    decoder.Decode(&sdnEvent)
    event := "StackDriver: " + sdnEvent.Instance.Policy_name +
        ", Condition: " + sdnEvent.Instance.Condition_name
    desc := "URL: " + sdnEvent.Instance.Url +
        "\nSummary: " + sdnEvent.Instance.Summary
    return event, desc
}
