package services

import (
    "encoding/json"
    "strconv"
)

type PDMessage struct {
    Action      string
    Check       string
    Checkname   string
    Description string
    Host        string
    IncidentId  int
}

type PDPayload struct {
    Message PDMessage
}

func getPingdomData(decoder *json.Decoder) (string, string) {
    var pdEvent PDPayload
    decoder.Decode(&pdEvent)
    message := pdEvent.Message
    event := "Pingdom: " + message.Description +
        " for " + message.Host
    desc := "Checkname: " + message.Checkname +
        "\nCheck: " + message.Check +
        "\nIncident Id: " + strconv.Itoa(message.IncidentId) +
        "\nAction: " + message.Action
    return event, desc
}
