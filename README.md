![Moov Banner Logo](https://github.com/moovfinancial/moov-go/assets/120951/3632d9ea-0c64-40e5-8f9e-b13b28b5e197)

[![GoDoc](https://godoc.org/github.com/moovfinancial/moovgo?status.svg)](https://godoc./github.com/moovfinancial/moov-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/moovfinancial/moov-go)](https://goreportcard.com/report/github.com/moovfinancial/moov-go)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moovfinancial/moov-go/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)


# Moov - Go Client
A [Go](http://golang.org) client for the [Moov payments API](https://docs.moov.io/api/). 


This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys. 

```bash 
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```

## Basic Usage

Example flow for linking a card and sending the money to a merchant's wallet.
[Checkout Example](./examples/checkout_example.go)

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
