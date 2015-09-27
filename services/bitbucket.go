package services

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
