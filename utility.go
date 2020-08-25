package main

import (
	"context"
	"google.golang.org/appengine/datastore"
)

// Return webhook datastore key.
func webhookKey(context context.Context, handler string) *datastore.Key {
	return datastore.NewKey(context, "Webhook", handler, 0, nil)
}

// Return AccessToken datastore key.
func accessTokenKey(context context.Context, email string) *datastore.Key {
	return datastore.NewKey(context, "AccessTokens", email, 0, nil)
}

// Return access token for provided email address.
func getAccessToken(context context.Context, email string) string {
	userAccessToken := datastore.NewQuery("AccessTokens").Ancestor(
		accessTokenKey(context, email)).Filter("Email =", email).Limit(1)
	aTokens := make([]AccessTokens, 0, 1)
	userAccessToken.GetAll(context, &aTokens)
	if len(aTokens) > 0 {
		return aTokens[0].AccessToken
	}
	return ""
}

// Return list of webhooks (datastore entities) for given email.
func getWebhooks(context context.Context, email string) []Webhook {
	query := datastore.NewQuery("Webhook").Filter("User =", email).Limit(50)
	webhooks := make([]Webhook, 0, 50)
	query.GetAll(context, &webhooks)
	return webhooks
}

// Return list of webhooks (datastore entities) from given handler.
func getWebhookFromHandler(
	context context.Context, handler string) *Webhook {
	query := datastore.NewQuery("Webhook").Ancestor(
		webhookKey(context, handler)).Limit(1)
	webhook := make([]Webhook, 0, 1)
	keys, _ := query.GetAll(context, &webhook)
	if len(webhook) > 0 {
		webhook[0].Count += 1
		datastore.Put(context, keys[0], &webhook[0])
		return &webhook[0]
	}
	return nil
}

// Delete handler.
func deleteWebhookFromHandler(
	context context.Context, handler string) *Webhook {
	query := datastore.NewQuery("Webhook").Ancestor(
		webhookKey(context, handler)).Limit(1)
	webhook := make([]Webhook, 0, 1)
	keys, _ := query.GetAll(context, &webhook)
	if len(webhook) > 0 {
		datastore.Delete(context, keys[0])
		return &webhook[0]
	}
	return nil
}
