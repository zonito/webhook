package webhook

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "math/rand"
    "net/url"
    "strconv"
    "strings"
    "time"
)

// Return webhook datastore key.
func webhookKey(context appengine.Context, handler string) *datastore.Key {
    return datastore.NewKey(context, "Webhook", handler, 0, nil)
}

// Return AccessToken datastore key.
func accessTokenKey(context appengine.Context, email string) *datastore.Key {
    return datastore.NewKey(context, "AccessTokens", email, 0, nil)
}

// Return TeleVerify datastore key.
func teleVerifyKey(context appengine.Context, code string) *datastore.Key {
    return datastore.NewKey(context, "TeleVerify", code, 0, nil)
}

// Return access token for provided email address.
func getAccessToken(context appengine.Context, email string) string {
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
func getWebhooks(context appengine.Context, email string) []Webhook {
    query := datastore.NewQuery("Webhook").Filter("User =", email).Limit(50)
    webhooks := make([]Webhook, 0, 50)
    query.GetAll(context, &webhooks)
    return webhooks
}

// Return list of webhooks (datastore entities) from given handler.
func getWebhookFromHandler(
    context appengine.Context, handler string) *Webhook {
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

// Return Chat id from Code
func getChatIdFromCode(context appengine.Context, code string) (int, string) {
    query := datastore.NewQuery("TeleVerify").Ancestor(
        teleVerifyKey(context, code)).Limit(1)
    teleVerify := make([]TeleVerify, 0, 1)
    keys, _ := query.GetAll(context, &teleVerify)
    if len(teleVerify) > 0 {
        chatId, name := teleVerify[0].ChatId, teleVerify[0].Name
        datastore.Delete(context, keys[0])
        return chatId, name
    }
    return 0, ""
}

// Return random alphanumeric string
