[![Build Status](https://storage.googleapis.com/pgwebhook/Webhook.png)](http://webhook.co)

# [webhook.co](http://webhook.co)

Webhook is a method of altering and augmenting the behavior of web applications, or a web pages, using some custom callbacks. These callbacks handled, modified, managed and maintained by some third-party developers and users who may not necessarily be affiliated with the origination web application or web page. 

Webhook.co offers the similar service for your web based application. It simply generates a URL which can be used in all your development application. Once your application performs any action such as git push, Travis pass or Travis fail etc. then Webhook.co will push the same information to your selected services. Isn’t is simple and amazing!

The popular applications which are already using web hooks include Assembla, CallMyApp, FreshBooks, Google Code, GitHub, Femtoo and PayPal etc. If you are also looking for a solution for your web hooks work then this is the right place to visit. 

## Important features of website.co:

Webhook.co provides a specific URL to users for various events, a web application will POST the data to those URLs when an action perform. It is very simple to use and it’s up to you and whatever you want to accomplish. Using our service, developers can create notification for themselves, real-time synchronization with other app, validate the data and also prevent it from being used by the application and process the data and repost using various APIs.

Webhook.co provide hooks for most popular platforms such as [Github](https://developer.github.com/webhooks/) (only push action), [Bitbucket](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html), [Travis](http://docs.travis-ci.com/user/notifications/#Webhook-notification), [Doorbell.co](https://doorbell.io/docs/webhooks), [TeamCity](https://www.jetbrains.com/teamcity/) and [Pingdom](https://help.pingdom.com/hc/en-us/articles/203611322-Setting-up-a-Webhook-and-an-Alerting-Endpoint). We also offer services for [Telegram](http://www.telegram.org), [Pushover](http://pushover.net), [Hipchat](http://hipchat.com) and [Trello](http://www.trello.com). The service includes three types of actions:

* **Push**: is used to receive data in real time.
* **Plugins**: is used to process data and return something.

Webhook.co let you offer a URL which is used to access all the APIs. The same URL works for accessing all the hooks and services offered. Please do create [new issue](https://github.com/PredictionGuru/webhook/issues) for any modification, features, queries in order to make our service more flexible, durable and easy to use.


## Webhook supported

* [Bitbucket](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html)
* [Github](https://developer.github.com/webhooks/) (Only Push)
* [Doorbell.co](https://doorbell.io/docs/webhooks)
* [Travis](http://docs.travis-ci.com/user/notifications/#Webhook-notification)
* [Pingdom](https://help.pingdom.com/hc/en-us/articles/203611322-Setting-up-a-Webhook-and-an-Alerting-Endpoint)
* [TeamCity](https://www.jetbrains.com/teamcity/)
* [Jenkins - Job Notification](https://wiki.jenkins-ci.org/display/JENKINS/Notification+Plugin)
* Anymore? Please contribute

## Connected Services

* [Trello](http://www.trello.com)
* [Telegram (@WebhookCo)](http://www.telegram.org)
* [Pushover](http://pushover.net)
* [Hipchat](http://hipchat.com)
* Pushbullet (Yet to start)
* Anymore? Please contribute

## Demo

* [webhook.co](http://webhook.co)

## Deploy to your Google App Engine

* Clone it `git clone`.
* instal and run `quack` [Quack](https://github.com/Autodesk/quack)
* Install and run `bower install`
* Add `services/keys.json` with respective details.
	```json
    {
      "pushoverKey": "",
      "trelloKey": "",
      "trelloSecret": "",
      "teleToken": ""
    }
    ```
* Deploy to your GAE application (Make sure you update `app.yaml`).

## Contributing

We <3 issue submissions, and address your problem as quickly as possible!

If you want to write code:

* Fork the repository
* Create your feature branch (`git checkout -b my-new-feature`)
* Commit your changes (`git commit -am 'add some feature'`)
* Push to your branch (`git push origin my-new-feature`)
* Create a new Pull Request


[![Build Status](http://38.media.tumblr.com/7d922f7b05a10891d00543c7a4acb79d/tumblr_inline_mk24hqGq6X1qz4rgp.jpg)](http://webhook.co)

**More about webhooks**: https://vimeo.com/4537957

![Analytics](https://ga-beacon.appspot.com/UA-68498210-1/webhook/repo)
