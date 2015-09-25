package services

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
