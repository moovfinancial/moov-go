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
