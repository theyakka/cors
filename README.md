CORS is a generic purpose library for controlling CORS preflight requests within Go server applications. It differs from other implementations in that it is not implemented as `http.Handler` directly. This allows you to control when / how the CORS preflight validation happens in your workflow.

CORS follows the  Cross Origin Resource Sharing W3 specification as described [here](http://www.w3.org/TR/cors).

CORS requires Go `1.13` (or higher) as it uses the error functionality that was introduced in Go `1.13`.

[![Go Version](https://img.shields.io/badge/Go-1.13+-lightgrey.svg)](https://golang.org/)

# Installing

To add to your Go project, just run:

```
go get -u github.com/theyakka/cors
```

# Getting started

TBD

See the [API documentation](http://godoc.org/github.com/theyakka/cors) for further details.

# Features

- Domain, Header and Method whitelisting
- Allows for wildcard Domains and Headers
- Allow credential option
- Max Age option

# Testing

To run all of the tests, execute the following in a terminal window: `go test`.

To run one of the tests individually, execute the following in a terminal window: `go test`.

# FAQ

## Why should I use this and not ____?

The short (and vague answer) is there is no reason. We find that supporting OSS isn't a competition. It's just "our way" of doing things that we thought others would find useful. 

In fact, CORS may / may not be the most appropriate choice for your server application. We built it for our specific purposes. If you're using a mostly vanilla http stack and looking for something more straightforward, you may wan to consider using [rs/cors](https://github.com/rs/cors).

We'd love for you to give it a try let us know if you like it! We love feedback either way.

## Has it been tested in production? Can I use it in production?

The code here has been written based on experiences with clients of all sizes. It has been production tested. That said, code is always evolving. We plan to keep on using it in production but we also plan to keep on improving it. If you find a bug, let us know!

## Who the f*ck is Yakka?

Yakka is an agency that makes better apps (be it mobile, web or whatever). We care about products, we work hard, we work fast and we write real code.

We use Flutter, Go and Vue (with some other stuff mixed in) and we like to think that we know how to use them to make awesome products.

Check us out at [https://theyakka.com](https://theyakka.com).

# Outro

## Credits

CORS is sponsored, owned and maintained by [Yakka LLC](http://theyakka.com). Feel free to reach out with suggestions, ideas or to say hey / hire us.

CORS draws significantly from the work done on the **rs/cors** library [here](https://github.com/rs/cors). In itself, the rs/cors library will probably be a better fit for most scenarios, we just needed something more "pluggable" for our situation.

## Specification issue

If you find an spot that we've missed the specification, please log an issue. Please use the spec issue template and we'll fix it ASAP.

## Security

If you believe you have identified a serious security vulnerability or issue with CORS, please report it as soon as possible to apps@theyakka.com. Please refrain from posting it to the public issue tracker so that we have a chance to address it and notify everyone accordingly.

# License

CORS is released under a modified MIT license. See LICENSE for details.