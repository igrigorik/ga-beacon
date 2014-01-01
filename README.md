# Google Analytics collector-as-a-service

![HN Button](http://img.skitch.com/20120415-bp8igiq74w53f91swt6tcy9cx8.jpg)


### Setup instructions

First, log in to your Google Analytics account and [set up a new property](https://support.google.com/analytics/answer/1042508?hl=en):

* Select "Website", use new "Universal Analytics" tracking
* Website name: anything you want.
* WebSite URL: https://ga-beacon.appspot.com/ (*important*)
* Click "Get Tracking ID", copy the UA-XXXXX-X id on next page.

Now, just add a tracking image to your repository:

https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page

* UA-XXXXX-X should be your tracking ID
* repo-name/page is an arbitrary path. For best results specify the repository name and the page name - e.g. if you have multiple readme's or wiki pages you can use different paths to map them to same repo: repo-name/readme, repo-name/wiki/some-other-page, and so on!

If you are using Markdown:

```
[![Analytics](https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page)](https://github.com/igrigorik/ga-beacon)
```

Or RDoc:

```
{<img src="https://ga-beacon.appspot.com/UA-XXXXX-X/repo-name/page" />}[https://github.com/igrigorik/ga-beacon]
```

If you prefer, you can skip the badge and use a transparent pixel. To do so, simply append `?pixel` to the image URL. 

And that's it, add the badge/pixel to the pages you want to track, and then head to your Google Analytics account to see real-time and aggregated visit analytics for your projects!


### FAQ

. How does this work?

GitHub does not allow arbitrary JavaScript to run on its pages. As a result, we can't use standard analytics snippets to track visitors and pageviews. However, Google Analytics provides a measurement protocol which allows us to POST the visit data directly to Google servers, and that's exactly what GA Beacon does:

* We include an image request on our pages which hits the GA Beacon service.
* GA Beacon POST's the visit to Google Analytics and returns the pixel. Hence the reason why we embed the tracking ID and page name in the URL.

. What about referrals and other visitor information?

Unfortunately we can't get all the same information via a tracking pixel. For example, referral information is only available on the GitHub page itself and can't be passed to the tracking pixel. As a result, the available metrics are restricted to unique visitors, pageviews, and the User-Agent of the visitor.

. Can I use this outside of GitHub?

Yep, you certainly can. It's a generic beacon service.


### Misc

* (MIT License) - Copyright (c) 2014 Ilya Grigorik
