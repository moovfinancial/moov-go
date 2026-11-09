# Example projects to learn the Moov-go SDK

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys.

```bash
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```

## Optional settings

- `MOOV_HOST` — API host. The default is `api.moov.io`.
- `PLAID_CLIENT_ID` and `PLAID_SECRET` — only for the Plaid examples.

The live tests use a shared sandbox account. These optional variables select the seeded resources. An absent or empty value uses the default:

- `MOOV_PARTNER_ID`
- `MOOV_PARTNER_WALLET_PM_ID`
- `MOOV_MERCHANT_ID`
- `MOOV_MERCHANT_WALLET_PM_ID`
- `MOOV_MERCHANT_WALLET_ID`

The live tests generate a random foreign ID for each process. CI and a local run can therefore share the seeded account without identifier collisions.


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
