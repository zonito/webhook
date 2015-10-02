package services

import (
    "encoding/json"
    "strconv"
    "strings"
)

// Bitbucket

type BBActor struct {
    Username     string
    Display_name string
    Uuid         string
    Links        string
}

type BBRepository struct {
    Name       string
    Full_name  string
    Uuid       string
    Links      BBLinks
    Scm        string
    Is_private bool
}

type BBAuthor struct {
    User User
    Raw  string
}

type BBLink struct {
    Href string
}

type BBLinks struct {
    Html BBLink
    Diff BBLink
    Self BBLink
}

type BBCommits struct {
    Message string
    Hash    string
    Type    string
    Date    string
    Author  BBAuthor
}

type BBNew struct {
    Type    string
    Name    string
    Target  BBCommits
    Parents BBCommits
}

type BBChanges struct {
    Commits []BBCommits
    New     BBNew
    Old     BBNew
    Created bool
    Forced  bool
    Closed  bool
}

type BBPush struct {
    Changes []BBChanges
}

type BBContent struct {
    Raw    string
    Markup string
    Html   string
}

type BBInline struct {
    To   int
    From int
    Path string
}

type BBParent struct {
    Id int
}

type BBComment struct {
    Id         int
    Parent     BBParent
    Content    BBContent
    Inline     BBInline
    Created_on string
    Updated_on string
    Links      BBLinks
}

type BBBranch struct {
    Name string
}

type BBSource struct {
    Branch     BBBranch
    Commit     BBCommits
    Repository BBRepository
}

type BBPullrequest struct {
    Id                  int
    Title               string
    Description         string
    State               string
    Author              string
    Source              BBSource
    Destination         BBSource
    Merge_commit        BBCommits
    Participants        []BBActor
    Reviewers           []BBActor
    Close_source_branch bool
    Closed_by           BBActor
    Reason              string
    Created_on          string
    Updated_on          string
    Links               BBLinks
}

type BBApproval struct {
    Date string
    User BBActor
}

type BBPayload struct {
    Repository  BBRepository
    Push        BBPush
    Actor       BBActor
    Fork        BBRepository
    Comment     BBComment
    Commit      BBCommits
    Pullrequest BBPullrequest
    Approval    BBApproval
}

// Return bitbucket data.
func getBitbucketData(decoder *json.Decoder, eType string) (string, string) {
    var bEvent BBPayload
    decoder.Decode(&bEvent)
    action := strings.Split(eType, ":")
    event, desc := "", ""
    who := bEvent.Actor.Username + " (" + bEvent.Actor.Display_name + ")"
    if action[0] == "repo" {
        switch action[1] {
        case "push":
            event = bEvent.Repository.Name + ": Push Event"
            if len(bEvent.Push.Changes) > 0 {
                change := bEvent.Push.Changes[0]
                desc = "Who: " + who + "\nCommits:\n"
                for i := 0; i < len(change.Commits); i++ {
                    desc += "\n- " + change.Commits[i].Message +
                        " (" + change.Commits[i].Hash + ")" +
                        "\n * " + change.Commits[i].Author.Raw
                }
            }
        case "fork":
            event = bEvent.Repository.Name + ": Fork Event"
            desc = "\n" + who + " Forked."
        case "commit_comment_created":
            event = bEvent.Repository.Name + ": Commit Comment Created"
            desc = who + " commented on " + bEvent.Commit.Hash
            desc += "\n Comment: " + bEvent.Comment.Content.Raw
            desc += "\n File: " + bEvent.Comment.Inline.Path
        }
    } else if action[0] == "pullrequest" {
        desc = "Description: " + bEvent.Pullrequest.Description
        desc += "\n From Repository: " +
            bEvent.Pullrequest.Source.Repository.Name
        switch action[1] {
        case "created":
            desc += "\n Created by: " + who
        case "updated":
            desc += "\n Updated by: " + who
        case "approved":
            desc += "\n Approved by: " + bEvent.Approval.User.Username
        case "unapproved":
            desc += "\n Unapproved by: " + bEvent.Approval.User.Username
        case "fulfilled":
            desc += "\n Merged by: " + who
            desc += "\n **Merged**"
        case "rejected":
            desc += "\n Rejected by: " + who
            desc += "\n **Rejected** because " + bEvent.Pullrequest.Reason
        case "comment_created":
            desc += "\n Commented by: " + who
            desc += "\n Comment: " + bEvent.Comment.Content.Raw
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        case "comment_updated":
            desc += "\n Comment updated by: " + who
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        case "comment_deleted":
            desc += "\n Comment deleted by: " + who
            desc += "\n File: " + bEvent.Comment.Inline.Path + " at line " +
                strconv.Itoa(bEvent.Comment.Inline.To)
        }
        event = bEvent.Repository.Name + ": Pull request " + action[1] + ": " +
            bEvent.Pullrequest.Title + "(" + bEvent.Pullrequest.State + ")"
    }
    return event, desc
}
