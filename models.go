package webhook

import "time"

type AccessTokens struct {
    Email       string
    AccessToken string // Access token of Trello
}

type Webhook struct {
    Handler   string    `json:"handler"`
    User      string    `json:"-"`
    BoardId   string    `json:"board_id"`
    BoardName string    `json:"board_name"`
    ListId    string    `json:"list_id"`
    ListName  string    `json:"list_name"`
    Date      time.Time `json:"date"`
    Count     int       `json:"count"`
}
