# Example projects to learn the Moov-go SDK 

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys. 

```bash 
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```


## ACH
[Debit a Bank Account](./ach/debit_bank_account/debit_bank_account_example.go) 

## Card Acquiring 
[Checkout](./card_acquiring/checkout/checkout_example.go)

## Debit Card Push / Pull 
[Original Credit Transaction](./debit_card_push/debit_push_test.go)
[Account Funding Transaction](./debit_card_pull/debit_pull_test.go)
