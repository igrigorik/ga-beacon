# Google Analytics for GitHub [![Analytics](https://ga-beacon.appspot.com/UA-71196-10/ga-beacon/readme?pixel)](https://github.com/igrigorik/ga-beacon)

[![GA Dashboard](https://www.evernote.com/shard/s1/sh/badbd5bb-228c-4304-8235-9bac60f6094e/a7c4a77d30584ada89b08fcd3b4c11d3/res/a6daea21-c886-497f-b906-f8bf6465a9b2/skitch.png)](https://www.evernote.com/shard/s1/sh/badbd5bb-228c-4304-8235-9bac60f6094e/a7c4a77d30584ada89b08fcd3b4c11d3/res/a6daea21-c886-497f-b906-f8bf6465a9b2/skitch.png)

**Curious which of your GitHub projects are getting all the traffic, or if anyone is reading your GitHub wiki pages?** Well, that's what Google Analytics is for! GitHub does not allow us to install arbitrary analytics, but we can still use a simple tracking image to log visits in real-time to Google Analytics - for full details, follow the instructions below. Once everything is setup, install [this custom dashboard](https://www.google.com/analytics/web/template?uid=MQS4cmZdSh2OWUVqRntqXQ) in your account for a nice real-time overview (as shown in above screenshot).


### Setup instructions

First, log in to your Google Analytics account and [set up a new property](https://support.google.com/analytics/answer/1042508?hl=en):

* Select "Website", use new "Universal Analytics" tracking
* **Website name:** anything you want (e.g. GitHub projects)
* **WebSite URL: https://ga-beacon.appspot.com/ (important)**
* Click "Get Tracking ID", copy the `UA-XXXXX-X` ID on next page.

Next, add a tracking image to the pages you want to track: 

* _https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page-name_
* `UA-XXXXX-X` should be your tracking ID
* `repo-name/page-name` is an arbitrary path. For best results specify the repository name and the page name - e.g. if you have multiple readme's or wiki pages you can use different paths to map them to same repo: `repo-name/readme`, `repo-name/other-page`, and so on!

Example tracker markup if you are using Markdown:

```markdown
[![Analytics](https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page)](https://github.com/igrigorik/ga-beacon)
```

Or RDoc:

```rdoc
{<img src="https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page" />}[https://github.com/igrigorik/ga-beacon]
```

If you prefer, you can skip the badge and use a transparent pixel. To do so, simply append `?pixel` to the image URL. 

And that's it, add the tracker image to the pages you want to track and then head to your Google Analytics account to see real-time and aggregated visit analytics for your projects!


### FAQ

- **How does this work?** GitHub does not allow arbitrary JavaScript to run on its pages. As a result, we can't use standard analytics snippets to track visitors and pageviews. However, Google Analytics provides a [measurement protocol](https://developers.google.com/analytics/devguides/collection/protocol/v1/devguide) which allows us to POST the visit data directly to Google servers, and that's exactly what GA Beacon does: we include an image request on our pages which hits the GA Beacon service, and GA Beacon POST's the visit to Google Analytics to record the visit.

- **Why do we need to proxy?** Google Analytics supports reporting of visit data [via GET requests](https://developers.google.com/analytics/devguides/collection/protocol/v1/reference#transport), but unfortunately we can't use that directly because we need to generate and report a unique visitor ID for each hit - GitHub does not allow us to run JS on the client to genrate the ID. To address this, we proxy the request through ga-beacon.appspot.com, which in turn is responsible for generating the unique visitor ID (hash of IP and User-Agent) and reporting the hit to Google Analytics.

- **What about referrals and other visitor information?** Unfortunately the static tracking pixel approach limits the information we can collect about the visit. For example, referral information is only available on the GitHub page itself and can't be passed to the tracking pixel. As a result, the available metrics are restricted to unique visitors, pageviews, and the User-Agent of the visitor.

- **Can I use this outside of GitHub?** Yep, you certainly can. It's a generic beacon service.

----

(MIT License) - Copyright (c) 2014 Ilya Grigorik
