package webhook

// To send as trello api payload.
type TrelloPayLoad struct {
  Name string
  Desc string
}

// Git structs

type GitUser struct {
  Name string
  Email string
  Username string
}

type GitRepository struct {
  Id int
  Name string
  Full_name string
  Url string
  AbsoluteUrl string
  Owner GitUser
  Pusher GitUser
}

type GitCommit struct {
  Id string
  Message string
  Timestamp string
  Url string
  Author GitUser
  Committer GitUser
  Modified []string
}

type GitPayload struct {
  Ref string
  Compare string
  Repository GitRepository
  Commits []GitCommit
}

// Doorbell structs

type DBApplication struct {
  Name string
}

type DBData struct {
  Email string
  Url string
  member GitUser
  User_Agent string
  Message string
  Sentiment string
  Application DBApplication
  Created string
}

type DBPayload struct {
  Event string
  Data DBData
}
