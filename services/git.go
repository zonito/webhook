package services

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
