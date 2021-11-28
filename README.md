# eater

_Not affiliated with Eater or Bitesnap_

Takes a Bitesnap JSON food log and makes it into something you could easily share with whoever you need to share it with. Currently only does photos and descriptions, because that's what I use it for.

## How it Runs
You can go to (eater.tech)[https://eater.tech] and see for yourself. Upload a JSON and go for it.

The lambda has a 1sec timeout and 128mb of RAM.

## Areas for Improvement
* This could almost certainly be a single html file with an embedded JS file, but sometimes you end up using the tools you're familiar with.
* The design is... obviously not my area of expertise. Any help with `foodlog.gohtml` and `style.css` would be appreciated
* A step to parse the time into something more readable than `YYYYMMDDHHMMSS` would be nice.