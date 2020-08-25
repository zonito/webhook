package services

import (
	"encoding/json"
	"strconv"
)

// Travis

type TRRepository struct {
	Id         int
	Name       string
	Owner_name string
	Url        string
}

type TRPayload struct {
	Author_email        string
	Author_name         string
	Branch              string
	Build_url           string
	Commit              string
	Committed_at        string
	Committer_email     string
	Committer_name      string
	Compare_url         string
	Duration            int
	Finished_at         string
	Id                  int
	Message             string
	Number              string
	Repository          TRRepository
	Started_at          string
	State               string
	Status              int
	Status_message      string
	Type                string
	Pull_request        bool
	Pull_request_number string
	Pull_request_type   string
	Tag                 string
}

// Return travis data.
func getTravisData(decoder *json.Decoder) (string, string) {
	var tEvent TRPayload
	decoder.Decode(&tEvent)
	if tEvent.Id > 0 {
		event := "Travis: " + tEvent.Status_message + " for " +
			tEvent.Repository.Name
		desc := "Status: " + tEvent.Status_message +
			"\n Duration: " + strconv.Itoa(tEvent.Duration) +
			"\n Message: " + tEvent.Message +
			"\n Build Number: " + tEvent.Number +
			"\n Type: " + tEvent.Type +
			"\n Compare URL: " + tEvent.Compare_url +
			"\n Committer Name: " + tEvent.Committer_name +
			"\n Build Url: " + tEvent.Build_url
		return event, desc
	}
	return "", ""
}
