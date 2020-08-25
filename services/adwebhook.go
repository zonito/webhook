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
	ClientId string
	AppId    string
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
	event := dEvent.Hook.CreatorType + " Delete: `" + dEvent.Hook.Status + "`"
	desc := "Delete Client: `" + dEvent.ResourceUrn +
		"`\nURN: `" + dEvent.Hook.Urn +
		"`\nTenant: `" + dEvent.Hook.Tenant +
		"`\nHookID: `" + dEvent.Hook.HookId +
		"`\nAppID: `" + dEvent.Payload.AppId + "`"
	return event, desc
}
