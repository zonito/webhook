package webhook

// To send as trello api payload.
type TrelloPayLoad struct {
    Name string
    Desc string
}

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

type GitPayload struct {
    Ref        string
    Compare    string
    Repository GitRepository
    Commits    []GitCommit
}

// Doorbell structs

type DBApplication struct {
    Name string
}

type DBData struct {
    Email       string
    Url         string
    member      User
    User_Agent  string
    Message     string
    Sentiment   string
    Application DBApplication
    Created     string
}

type DBPayload struct {
    Event string
    Data  DBData
}

// Bitbucket

type BBRepository struct {
    Name       string
    Is_private bool
}

type BBAuthor struct {
    User User
    Raw  string
}

type BBCommits struct {
    Message string
    Hash    string
    Author  BBAuthor
}

type BBChanges struct {
    Commits []BBCommits
}

type BBPush struct {
    Changes []BBChanges
}

type BBPayload struct {
    Repository BBRepository
    Push       BBPush
}
