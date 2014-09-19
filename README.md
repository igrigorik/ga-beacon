# Google Analytics Beacon [![Analytics](https://ga-beacon.appspot.com/UA-71196-10/ga-beacon/readme?pixel)](https://github.com/igrigorik/ga-beacon)

Sometimes it is impossible to embed the JavaScript tracking code provided by Google Analytics: the host page does not allow arbitrary JavaScript, and there is no Google Analytics integration. However, not all is lost! **If you can embed a simple image (pixel tracker), then you can beacon data to Google Analytics.** For a great, hands-on explanation of how this works, check out the following guides:

* [Using a Beacon Image for GitHub, Website and Email Analytics](http://www.sitepoint.com/using-beacon-image-github-website-email-analytics/)
* [Tracking Google Sheet views with Google Analytics using GA Beacon](http://mashe.hawksey.info/2014/02/tracking-google-sheet-views-with-google-analytics/)


### Hands-on example: Google Analytics for GitHub 

_Note: GitHub [released traffic analytics](https://github.com/blog/1672-introducing-github-traffic-analytics) on Jan 7, 2014 and provides deeper analytics (e.g. referrer tracking) than what is possible via GA-Beacon. If you're interested in analytics for your GitHub project, use that... That said, if you want real-time analytics, or want to see a demo of what's possible with a tracking pixel, see below._

[![GA Dashboard](https://lh5.googleusercontent.com/-Zu9r9m7Uv0c/UsSQlJ5OoeI/AAAAAAAAHwo/fvH_lrVUV0w/w1007-h467-no/skitch.png)](https://lh5.googleusercontent.com/-Zu9r9m7Uv0c/UsSQlJ5OoeI/AAAAAAAAHwo/fvH_lrVUV0w/w1007-h467-no/skitch.png)

**Curious which of your GitHub projects are getting all the traffic, or if anyone is reading your GitHub wiki pages?** Well, that's what Google Analytics is for! GitHub does not allow us to install arbitrary analytics, but we can still use a simple tracking image to log visits in real-time to Google Analytics - for full details, follow the instructions below. Once everything is setup, install [this custom dashboard](https://www.google.com/analytics/web/template?uid=MQS4cmZdSh2OWUVqRntqXQ) in your account for a nice real-time overview (as shown in above screenshot).

_Note: GitHub proxies all assets through their own fetch service. As a result, while you can still use the tracking pixel, the fetch service acts as an intermediary that (intentionally, in this case) hides a lot of information: visitors IP address, does not store cookies, etc. As a result, each hit will be counted and reported as a unique visitor in GA._


### Setup instructions

First, log in to your Google Analytics account and [set up a new property](https://support.google.com/analytics/answer/1042508?hl=en):

* Select "Website", use new "Universal Analytics" tracking
* **Website name:** anything you want (e.g. GitHub projects)
* **WebSite URL: https://ga-beacon.appspot.com/**
* Click "Get Tracking ID", copy the `UA-XXXXX-X` ID on next page

Next, add a tracking image to the pages you want to track:

* _https://ga-beacon.appspot.com/UA-XXXXX-X/your-repo/page-name_
* `UA-XXXXX-X` should be your tracking ID
* `your-repo/page-name` is an arbitrary path. For best results specify the repository name and the page name - e.g. if you have multiple readme's or wiki pages you can use different paths to map them to same repo: `your-repo/readme`, `your-repo/other-page`, and so on!

Example tracker markup if you are using Markdown:

```markdown
[![Analytics](https://ga-beacon.appspot.com/UA-XXXXX-X/your-repo/page-name)](https://github.com/igrigorik/ga-beacon)
```

Or RDoc:

```rdoc
{<img src="https://ga-beacon.appspot.com/UA-XXXXX-X/your-repo/page-name" />}[https://github.com/igrigorik/ga-beacon]
```

If you prefer, you can skip the badge and use a transparent pixel. To do so, simply append `?pixel` to the image URL.

And that's it, add the tracker image to the pages you want to track and then head to your Google Analytics account to see real-time and aggregated visit analytics for your projects!


### FAQ

- **How does this work?** Google Analytics provides a [measurement protocol](https://developers.google.com/analytics/devguides/collection/protocol/v1/devguide) which allows us to POST arbitrary visit data directly to Google servers, and that's exactly what GA Beacon does: we include an image request on our pages which hits the GA Beacon service, and GA Beacon POSTs the visit data to Google Analytics to record the visit. As a result, if you can embed an image, you can beacon data to Google Analytics.

- **Why do we need to proxy?** Google Analytics supports reporting of visit data [via GET requests](https://developers.google.com/analytics/devguides/collection/protocol/v1/reference#transport), but unfortunately we can't use that directly because we need to generate and report a unique visitor ID for each hit - e.g. GitHub does not allow us to run JS on the client to generate the ID. To address this, we proxy the request through ga-beacon.appspot.com, which in turn is responsible for generating the unique visitor ID (server generated UUID), setting the appropriate cookies for repeat hits, and reporting the hits to Google Analytics.

- **What about referrals and other visitor information?** Unfortunately the static tracking pixel approach limits the information we can collect about the visit. For example, referral information can't be passed to the tracking pixel because we can't execute JavaScript. As a result, the available metrics are restricted to unique visitors, pageviews, and the User-Agent and IP address of the visitor.

- **Can I use this outside of GitHub?** Yep, [you certainly can](http://www.sitepoint.com/using-beacon-image-github-website-email-analytics/). It's a generic beacon service.

- **Do I have to use ga-beacon.appspot.com?** You can if you want to - it's free. Alternatively, feel free to deploy your own instance directly on Google App Engine. The project is under MIT license.
