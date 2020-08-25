package services

import (
	"encoding/json"
)

type TCBuild struct {
	BuildStatus           string
	BuildResult           string
	BuildResultPrevious   string
	BuildResultDetalta    string
	NotifyType            string
	BuildFullName         string
	BuildName             string
	BuildId               string
	BuildTypeId           string
	BuildInternalTypeId   string
	BuildExternalTypeId   string
	BuildStatusUrl        string
	BuildStatusHtml       string
	RootUrl               string
	ProjectName           string
	ProjectId             string
	ProjectInternalId     string
	ProjectExternalId     string
	BuildNumber           string
	AgentName             string
	AgentOs               string
	AgentHostname         string
	TriggeredBy           string
	Message               string
	Text                  string
	BuildStateDescription string
}

type TCPayload struct {
	Build TCBuild
}

// Return teamcity data.
func getTeamcityData(decoder *json.Decoder) (string, string) {
	var tcEvent TCPayload
	decoder.Decode(&tcEvent)
	build := tcEvent.Build
	event := build.ProjectName + ": " +
		build.BuildResult + " (Previous: " + build.BuildResultPrevious + ") "
	desc := build.Message + "\nBuild Status Url: " +
		build.BuildStatusUrl
	return event, desc
}
