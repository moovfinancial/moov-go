
## v0.4.4 (Released 2024-10-11)

ADDITIONS

- feat: adding in Scheduled Transfers to the client.
- feat: adding in Wallet Sweeps to the client.

## v0.4.3 (Released 2024-08-05)

IMPROVEMENTS

- feat: add support for Instant Micro Deposits (Bank Account Verification)

## v0.4.2 (Released 2024-08-02)

IMPROVEMENTS

- feat: add Commercial and Regulated field on Card
- feat: add end to end encryption endpoints and example

## v0.4.1 (Released 2024-07-23)

IMPROVEMENTS

- fix: specify http method in patch card endpoint

## v0.4.0 (Released 2024-07-22)

IMPROVEMENTS

- feat: allow custom content type on multipart file
- fix: Ensure cards patch only updates intended fields
- test: always cleanup created card
- test: use known account for pull from card test

## v0.3.0 (Released 2024-06-20)

BREAKING CHANGES

- transfers: update CreateTransfer to constrain return params based on async vs sync

ADDITIONS

- feat: add primary regulator to business profile
- examples: add business profile

IMPROVEMENTS

- feat: support webhooks ping event type and payload
- tests: enforce exact JSON decoding
- examples: retry list payment methods after bank account verification

BUG FIXES

- webhooks: fix JSON parsing error when decoding non-base64 encoded data

## v0.2.0 (Released 2024-05-16)

ADDITIONS

- examples: add webhooks receiver
- feat: add missing fields for Disputes, add endpoints
- feat: create CRUD commands for Representatives
- feat: create CRUD commands for Underwriting
- webhooks: add mhooks package, new models, and establish pattern

IMPROVEMENTS

- examples/ach: add wallet-to-bank example
- examples/rtp: add RTP fallback to ACH
- feat: add RTPRejectionCode on BankAccount ExceptionDetails
- meta: send moov-go version on User-Agent
- moov: raise query params to their noun level types

Milestone: https://github.com/moovfinancial/moov-go/issues?q=is%3Aissue+milestone%3Av0.2.0+

## v0.1.0 (Released 2024-04-23)

INITIAL RELEASE

We are excited to announce the initial release of Moov Financial's Go SDK! This library is already being
utilized in various production environments. We highly value your input and experiences as we continue to
enhance the SDK. Please don't hesitate to share your feedback, report issues, or discuss your use cases
with us.
