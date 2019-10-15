# gos3headersetter

[![Build Status](https://travis-ci.org/cariad/gos3headersetter.svg?branch=master)](https://travis-ci.org/cariad/gos3headersetter) [![Go Report Card](https://goreportcard.com/badge/github.com/cariad/gos3headersetter)](https://goreportcard.com/report/github.com/cariad/gos3headersetter) [![](https://godoc.org/github.com/cariad/gos3headersetter?status.svg)](http://godoc.org/github.com/cariad/gos3headersetter) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/cariad/gos3headersetter/blob/master/LICENSE)

## Introduction

**gos3headersetter** is a Golang module for setting the `Cache-Control` and `Content-Type` HTTP headers on S3 objects.

This is useful for static websites hosted out of S3:

- S3 can sometimes guess the `Content-Type` of objects for you, but not always.
- Setting `Cache-Control` to enable CloudFront and browser caching can save you bandwith (and money).

## Limitations

Only `Cache-Control` and `Content-Type` headers are currently supported.

## Command-line application

This project is just a Golang module. A command-line application version is available at [cariad/s3headersetter](https://github.com/cariad/s3headersetter).

I use that command-line application in my [Hugo website deployment pipeline](https://github.com/cariad/aws-hugo).

## Rules

The `Rule` struct describes a header and the value to set under certain circumstances:

- `Header (string)` describes the name of the header which the rule affects.
- `When ([]gos3headersetter.When)` describes the values to set for specific key (filename) extensions.
- `Else (string, optional)` describes the value to set for objects which were not matched by a `When` statement.

For example, this rule will:

- Set `Cache-Control` to `max-age=3600, public` on `.html` objects.
- Set `Cache-Control` to `max-age=604800, public` on `.css` objects.
- Set `Cache-Control` to `max-age=31536000, public` on all other objects.

```golang
rule := gos3headersetter.Rule {
    Header: "Cache-Control",
    When:   []gos3headersetter.When {
        gos3headersetter.When{
            Extension: ".html",
            Then:      "max-age=3600, public",
        },
        gos3headersetter.When{
            Extension: ".css",
            Then:      "max-age=604800, public",
        },
    },
    Else:   "max-age=31536000, public"
}
```

## Usage

To update a specific object, use the `Object` struct:

```golang
package main

import (
    "fmt"
    "github.com/cariad/gos3headersetter"
)

func main() {
    object := gos3headersetter.NewObject("my-bucket", "public/index.html")
    object.Apply(rules)
}
```

Note that you should use the `NewObject()` constructor rather than creating a new `Object` instance directly.

To update all the objects in a bucket, use the `Bucket` struct. The `KeyPrefix` is optional.

```golang
package main

import (
    "fmt"
    "github.com/cariad/gos3headersetter"
)

func main() {
    bucket := gos3headersetter.Bucket{
        Bucket:    "my-bucket",
        KeyPrefix: "public",
    )

    bucket.Apply(rules)
}
```

## Licence, credit & sponsorship

This project is published under the MIT Licence.

You don't owe me anything in return, but as an indie freelance coder there are two things I'd appreciate:

- **Credit.** If your app or documentation has a credits page, please consider mentioning the projects you use.
- **Cash.** If you want *and are able to* support future development, please consider [becoming a patron](https://www.patreon.com/cariad) or [buying me a coffee](https://ko-fi.com/cariad). Thank you!
