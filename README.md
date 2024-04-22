![Moov Logo](https://github.com/moovfinancial/moov-go/assets/120951/3632d9ea-0c64-40e5-8f9e-b13b28b5e197)

[![GoDoc](https://pkg.go.dev/badge/github.com/moovfinancial/moovgo?utm_source=godoc)](https://pkg.go.dev/github.com/moovfinancial/moov-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/moovfinancial/moov-go)](https://goreportcard.com/report/github.com/moovfinancial/moov-go)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moovfinancial/moov-go/master/LICENSE)
[![X](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)


# Moov - Go Client
The official [Go](http://golang.org) client for the [Moov payments API](https://docs.moov.io/api/).

## Installation

```go
import (
	"github.com/moovfinancial/moov-go/pkg/moov"
)
```

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account.

> [!NOTE]
> Note that API Keys for Sandbox and Production are different keys.

```bash
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```

In your Go program, create a new Moov client initiated with your public and secret keys.

```go
mc, err := moov.NewClient(
  moov.WithCredentials(moov.CredentialsFromEnv()), // optional, default is to read from environment
)
```

## Examples Usage

Checkout the [examples](./examples/README.md) exist for [ach](./examples/ach/), [card acquiring](./examples/card_acquiring/), debit [push](./examples/debit_card_push/)/[pull](./examples/debit_card_pull/), and [rtp](./examples/rtp/).

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
