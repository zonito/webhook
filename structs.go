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
