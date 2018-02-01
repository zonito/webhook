package services

import (
    "encoding/json"
)

// AD Webhooks structs

type ADScope struct {
    application string
}

type ADHook struct {
    HookId      string
    Tenant      string
    CallbackUrl string
    CreatedBy   string
    EventType   string
    CreatedDate string
    SysType     string
    CreatorType string
    Status      string
    Scope       ADScope
    Urn         string
    __self__    string
}

type ADClientPayload struct {
    ClientId    string
}

type ADPayload struct {
    ResourceUrn string
    Payload     ADClientPayload
    Hook        ADHook
}

// Return AD data.
func getADData(decoder *json.Decoder) (string, string) {
    var dEvent ADPayload
    decoder.Decode(&dEvent)
    event := dEvent.Hook.CreatorType + " Delete: " + dEvent.Hook.Status
    desc := "`" + dEvent.ResourceUrn + "` created by `" +
        dEvent.Hook.CreatedBy + "`\nURN: `" + dEvent.Hook.Urn +
        "`,\nEventType: `" + dEvent.Hook.EventType + "`,\nTenant: `" +
        dEvent.Hook.Tenant + "`"
    return event, desc
}
