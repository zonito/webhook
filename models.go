package webhook

import (
	"time"
)

type AccessTokens struct {
  Email string
  AccessToken string		// Access token of Trello
}

type Webhook struct {
  Handler string		// Hook handler to receive events from services.
  Email string			// Email of user to whom its mapped with.
  BoardId string		// Board Id of trello
  ListId string			// List Id of trello
  Date time.Time		// Created Date
}
