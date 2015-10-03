// Package services provides Github Integration.

package services

import (
    "encoding/json"
)

// Git structs

type User struct {
    Name         string
    Email        string
    Username     string
    Display_name string
}

type GitRepository struct {
    Id          int
    Name        string
    Full_name   string
    Url         string
    AbsoluteUrl string
    Owner       User
    Pusher      User
}

type GitCommit struct {
    Id        string
    Message   string
    Timestamp string
    Url       string
    Author    User
    Committer User
    Modified  []string
}

type GitUser struct {
    Login      string
    Id         int
    Avatar_url string
    Type       string
    Site_admin bool
}

type GitPullRequest struct {
    Url       string
    Html_url  string
    Id        int
    State     string
    Title     string
    User      GitUser
    Body      string
    Repo      GitRepository
    Merged    bool
    Merged_by GitUser
}

type GitPayload struct {
    Zen          string
    Ref          string
    Compare      string
    Repository   GitRepository
    Commits      []GitCommit
    Action       string
    Number       int
    Pull_request GitPullRequest
    Pusher       User
}

// Return github data.
func getGithubData(decoder *json.Decoder, header string) (string, string) {
    var gEvent GitPayload
    decoder.Decode(&gEvent)
    var event, desc string
    if header == "push" {
        event = gEvent.Repository.Name + " --> " + header + " event"
        repo := gEvent.Repository
        desc = repo.Name + ": \n" +
            "\nName: " + repo.Name +
            "\nUrl: " + repo.Url +
            "\nOwner: " + repo.Owner.Email +
            "\nCompare: " + gEvent.Compare +
            "\nRef: " + gEvent.Ref +
            "\nModified files\n"
        for i := 0; i < len(gEvent.Commits); i++ {
            commit := gEvent.Commits[i]
            desc += "\n* " + commit.Message + " (" + commit.Timestamp + ")"
            for j := 0; j < len(commit.Modified); j++ {
                desc += "\n * " + commit.Modified[j]
            }
        }
    } else if header == "pull_request" {
        pr := gEvent.Pull_request
        if gEvent.Action == "opened" {
            event = "New pull request for " + gEvent.Repository.Full_name +
                " from " + pr.User.Login
        } else if gEvent.Action == "closed" && pr.Merged {
            event = "Pull request merged by " + pr.Merged_by.Login
        }
        desc = "Title: " + pr.Title
        if pr.Body != "" {
            desc += "\nDescription: " + pr.Body
        }
        desc += "\nReview at " + pr.Html_url
    } else if gEvent.Zen != "" {
        event = "Ping! from " + gEvent.Repository.Name
    }
    return event, desc
}
