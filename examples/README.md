# Example projects to learn the Moov-go SDK

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys.

```bash
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```


## ACH
- [Debit a Bank Account with micro deposits](./ach/debit_bank_account/micro_deposits_test.go)
- [Link via Plaid Processor Tokens](./ach/debit_bank_account/plaid_processors_test.go)
- [Credit a Bank Account with micro deposits](./ach/credit_external_bank/micro_deposits_test.go)

## Card Acquiring
- [Checkout](./card_acquiring/checkout/checkout_example.go)

## Debit Card Push / Pull
- [Account Funding Transaction](./debit_card_pull/debit_pull_test.go)
- [Original Credit Transaction](./debit_card_push/debit_push_test.go)

## RTP
- [Send an RTP Transaction](./rtp/rtp_credit_test.go)
