## Blogger Migration Temporary Service

This was solving  my blog article migration problem, move articles to another without losing google-rank. 

This service redirect old article requests to new location by sending HTTP 301, those suggest by Google SearchConsole post, [Change page URLs with 301 redirects](https://support.google.com/webmasters/answer/93633).

More detail, please see [here](http://dev.twsiyuan.com/2017/01/blogger-move-without-losing-google-rank.html) (in Traditional Chinese).

## Install

First you need Go environment, then install package dependency.

```
go get github.com/codegangsta/negroni
go get github.com/gorilla/mux
```

## License

CC0, No Rights Reserved.