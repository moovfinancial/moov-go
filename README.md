# Moov - Go Client
A [Go](http://golang.org) client for the [Moov payments API](https://docs.moov.io/api/). 


## Installation 

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys. 


```bash 
# Go Client 
go get github.com/moovfinancial/moov-go/
```

API Keys can be configured with environmental variables.

```bash 
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```

## Basic Usage 

Checkout flow for linking a card and sending the money to a merchants wallet.
[Checkout Example](./examples/checkout_example.go)
