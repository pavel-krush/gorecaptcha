gorecaptcha
===========

gorecaptcha package is a go library for verifying [reCaptcha](https://www.google.com/recaptcha)'s user response.

Works with both v2 and v3 versions.

Usage
-----

Install:

```
go get github.com/pavel-krush/gorecaptcha
```

Create Recaptcha Instance:

```
recaptcha := gorecaptcha.New(recaptchaSecret)
```

Verify:

```
response, _ = gorecaptcha.Verify(ip, recaptchaToken)
if response.Success {
    // captcha was completed successfully
} else {
    // captcha was not completed successfully
}
```

Context:

This package also has context support. Use `recaptcha.VerifyContext` instead of `recaptcha.Verify`.

Custom Transport:

It's possible to use custom transport for http requests.

```
recaptcha = gorecaptcha.New(...).WithHTTPClient(customHTTPClient)
```
