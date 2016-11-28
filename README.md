# Google Analytics Beacon [![Analytics](https://ga-beacon.appspot.com/UA-71196-10/ga-beacon/readme?pixel)](https://github.com/igrigorik/ga-beacon)

Sometimes it is impossible to embed the JavaScript tracking code provided by Google Analytics: the host page does not allow arbitrary JavaScript, and there is no Google Analytics integration. However, not all is lost! **If you can embed a simple image (pixel tracker), then you can beacon data to Google Analytics.** For a great, hands-on explanation of how this works, check out the following guides:

* [Using a Beacon Image for GitHub, Website and Email Analytics](http://www.sitepoint.com/using-beacon-image-github-website-email-analytics/)
* [Tracking Google Sheet views with Google Analytics using GA Beacon](http://mashe.hawksey.info/2014/02/tracking-google-sheet-views-with-google-analytics/)

### Can I use this production?

The ga-beacon.appspot.com instance is a **demo** instance, good for prototyping and proof of concepts. If you intend to use this in production for your application, you should deploy **your own instance** of this service, which will allow you to scale the service up and down to meet your capacity needs, introspect the logs, customize the code, and so on.

Deploying your own instance is trivial: fork this repo, modify the project name in app.yaml, and follow the [normal GAE deploy instructions](https://cloud.google.com/appengine/training/go-plus-appengine/deploy). You should be up and running in less than five minutes.

### Setup instructions

First, log in to your Google Analytics account and [set up a new property](https://support.google.com/analytics/answer/1042508?hl=en):

* Select "Website", use new "Universal Analytics" tracking
* **Website name:** anything you want (e.g. GitHub projects)
* **WebSite URL: https://ga-beacon.appspot.com/**
* Click "Get Tracking ID", copy the `UA-XXXXX-X` ID on next page

Next, add a tracking image to the pages you want to track:

* _https://ga-beacon.appspot.com/UA-XXXXX-X/insert/any/path_
* `UA-XXXXX-X` should be your tracking ID
* `insert/any/path` is an arbitrary path. For best results specify a meaningful and self-descriptive path. You have to do this manually, the beacon won't automatically record the page path it's embedded on.

Example tracker markup if you are using Markdown:

```markdown
[![Analytics](https://ga-beacon.appspot.com/UA-XXXXX-X/welcome-page)](https://github.com/igrigorik/ga-beacon)
```

Or RDoc:

```rdoc
{<img src="https://ga-beacon.appspot.com/UA-XXXXX-X/welcome-page" />}[https://github.com/igrigorik/ga-beacon]
```

If you prefer, you can skip the badge and use a transparent pixel. To do so, simply append `?pixel` to the image URL. There are also "flat" style variants available, which are available when appending `?flat` or `?flat-gif` to the image URL. And that's it, add the tracker image to the pages you want to track and then head to your Google Analytics account to see real-time and aggregated visit analytics for your projects!

You may also auto-calculate the tracking path based in the "referer" information of the image. To activate this simple add `?useReferer` to the image URL (or `&useReferer` if you need to combine this with the `?pixel`, `?flat` or `?flat-gif` parameter). Although they are some odd browsers that don't always send the referer header, the amount of traffic coming from those browsers is usually not relevant at all. Of course that if you need to measure the traffic from those odd browsers you should not use this method.

### FAQ

- **How does this work?** Google Analytics provides a [measurement protocol](https://developers.google.com/analytics/devguides/collection/protocol/v1/devguide) which allows us to POST arbitrary visit data directly to Google servers, and that's exactly what GA Beacon does: we include an image request on our pages which hits the GA Beacon service, and GA Beacon POSTs the visit data to Google Analytics to record the visit. As a result, if you can embed an image, you can beacon data to Google Analytics.

- **Why do we need to proxy?** Google Analytics supports reporting of visit data [via GET requests](https://developers.google.com/analytics/devguides/collection/protocol/v1/reference#transport), but unfortunately we can't use that directly because we need to generate and report a unique visitor ID for each hit - e.g. some pages do not allow us to run JS on the client to generate the ID. To address this, we proxy the request through ga-beacon.appspot.com, which in turn is responsible for generating the unique visitor ID (server generated UUID), setting the appropriate cookies for repeat hits, and reporting the hits to Google Analytics.

- **What about referrals and other visitor information?** Unfortunately the static tracking pixel approach limits the information we can collect about the visit. For example, referral information can't be passed to the tracking pixel because we can't execute JavaScript. As a result, the available metrics are restricted to unique visitors, pageviews, and the User-Agent and IP address of the visitor.

- **Do I have to use ga-beacon.appspot.com?** You can if you want to - it's free, but there are no capacity or availability promises. For best results, deploy your own instance directly on Google App Engine: clone this repository, change the project name, and deploy your own instance - easy as that. The project is under MIT license.
