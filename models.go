package webhook

import (
	"time"
)

type AccessTokens struct {
  Email string
  AccessToken string
}

type Webhook struct {
  Handler string
  Email string
  BoardId string
  ListId string
  Date time.Time
}
