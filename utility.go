package webhook

import (
	"appengine"
	"appengine/datastore"
)

func webhookKey(context appengine.Context) *datastore.Key {
	return datastore.NewKey(context, "Webhook", "default_webhook", 0, nil)
}

func accessTokenKey(context appengine.Context) *datastore.Key {
	return datastore.NewKey(context, "AccessTokens", "default_at", 0, nil)
}

func getAccessToken(context appengine.Context, email string) string {
	userAccessToken := datastore.NewQuery("AccessTokens").Ancestor(
		accessTokenKey(context)).Filter("Email =", email).Limit(1)
	aTokens := make([]AccessTokens, 0, 1)
	userAccessToken.GetAll(context, &aTokens)
	if len(aTokens) > 0 {
		return aTokens[0].AccessToken
	}
	return ""
}

func getAccessTokenFromHandler(context appengine.Context, handler string) string {
	webhook := getWebhookFromHandler(context, handler)
	if webhook != nil {
		return getAccessToken(context, webhook.Email)
	}
	return ""
}

func getWebhooks(context appengine.Context, email string) []Webhook {
	query := datastore.NewQuery("Webhook").Ancestor(
		webhookKey(context)).Filter("Email =", email).Order("-Date").Limit(10)
	webhooks := make([]Webhook, 0, 10)
	query.GetAll(context, &webhooks)
	return webhooks
}

func getWebhookFromHandler(context appengine.Context, handler string) *Webhook {
	query := datastore.NewQuery("Webhook").Ancestor(
		webhookKey(context)).Filter("Handler =", handler).Limit(1)
	webhook := make([]Webhook, 0, 1)
	query.GetAll(context, &webhook)
	if len(webhook) > 0 {
		return &webhook[0]
	}
	return nil
}
