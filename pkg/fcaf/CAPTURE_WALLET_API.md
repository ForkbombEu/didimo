<!--
SPDX-FileCopyrightText: 2026 Forkbomb BV

SPDX-License-Identifier: AGPL-3.0-or-later
-->

# Capture Wallet API capability reference

`https://beta-capture-wallet.credimi.io` is a stateful OpenID4VCI issuer and OpenID4VP verifier used to capture Wallet protocol evidence. This reference answers a narrow FCAF implementation question: can the public service create, deliver, and observe the protocol exchange needed by a test?

It is not a substitute for the OpenID4VCI, OpenID4VP, DCQL, or FCAF specifications. A service accepting a session-creation body does not prove that the same property reached the Wallet in its signed request, nor that the service captured the Wallet response.

## Sources and confidence

Published contract source: [Capture Wallet API documentation](https://beta-capture-wallet.credimi.io/docs) and its linked OpenAPI document. Entries marked **published** below come from that OpenAPI document. Entries marked **observed** come from the named local evidence record. Do not treat an unlisted field or behaviour as supported.

| Status | Meaning |
| --- | --- |
| Supported | Published as an input or endpoint. Still inspect the delivered request and captured response for the individual test. |
| Observed | Verified in a dated Credimi evidence record; it may differ across deployments. |
| Unknown | Not established by the published contract or local evidence. It requires a probe before it can justify a test implementation. |
| Blocked | The current service or reference wallet is known not to produce the required evidence. |

## Decision sequence for an FCAF test

Classify the test at all three boundaries, in order:

1. **Session input:** can `POST /sessions` or `POST /openid4vp/sessions` express the requested setup?
2. **Wallet delivery:** can `GET /openid4vp/sessions/{sessionId}/request` or the session/deeplink evidence prove that the exact resulting Authorization Request contains the required property?
3. **Evidence capture:** can the relevant session record or events prove the Wallet's protocol response, rather than only its UI state?

If a required property is rejected before a signed request is delivered, mark the case verifier-blocked; do not create synthetic Wallet evidence. If the Wallet is shown the request but does not submit the required response, record a reference-wallet failure or discontinuation only when the FCAF source permits it.

## Shared service and metadata

| Endpoint | Capability | Status |
| --- | --- | --- |
| `GET /healthz` | Readiness response `{ "status": "ok" }`. | Supported |
| `GET /issuers` | Lists the always-on issuer configurations, metadata URLs, warnings, and credential configuration IDs. | Supported |
| `GET /oid4vci/requests` | Bounded chronological OpenID4VCI request ledger. Sensitive values are redacted to presence/length metadata. | Supported |
| `GET /.well-known/openid-credential-issuer/issuers/{issuerConfigurationId}` | OpenID4VCI issuer metadata; send `Accept: application/jwt` for signed metadata. | Supported |
| `GET /.well-known/oauth-authorization-server/issuers/{issuerConfigurationId}` | OAuth authorization-server metadata. | Supported |
| `GET /.well-known/jwt-vc-issuer/issuers/{issuerConfigurationId}` | JWT VC issuer metadata. | Supported |
| `GET /issuers/{issuerConfigurationId}/jwks.json` | Authorization-server signing JSON Web Key Set. | Supported |
| `GET /issuers/{issuerConfigurationId}/credential-jwks.json` | Credential-signing JSON Web Key Set. | Supported |

Use the metadata endpoints rather than hard-coding credential configuration identifiers, authorization-server settings, or issuer keys in a test.

## Issuer: OpenID4VCI capture sessions

### Session surface

| Endpoint | Inputs / output | Status |
| --- | --- | --- |
| `POST /sessions` | Optional JSON: `issuer_configuration_id` (`eu-pid-device-bound` or `eu-pid-jwt-proof-only`), `flow` (`pre_authorized_code` or `authorization_code`), `credential_offer_mode` (`credential_offer` or `credential_offer_uri`), metadata-advertised `credential_configuration_id`, and `status_list_enabled` (boolean). Returns the selected issuer and authorization-server identifiers plus session and offer details. | Supported |
| `GET /sessions/{sessionId}` | Current issuance capture containing `observed`, `checks`, and `events` in addition to session state. | Supported |
| `GET /sessions/{sessionId}/offer` | Credential offer. | Supported |
| `GET /sessions/{sessionId}/deeplink` | `deeplink` and `credential_offer`. | Supported |
| `GET /sessions/{sessionId}/jwks` | Wallet holder-binding JWKS observed from the proof header; returns `409` until a proof header JWK exists. | Supported |
| `GET /sessions/{sessionId}/events` | Chronological events with timestamp, type, and arbitrary detail. | Supported |

The service has two always-on issuer configurations, `eu-pid-device-bound` and
`eu-pid-jwt-proof-only`. Obtain their credential configuration IDs and protocol
metadata from `GET /issuers` instead of hard-coding them.

| Issuer configuration ID | Reference EUDI Wallet version |
| --- | --- |
| `eu-pid-device-bound` | Version 39 and later |
| `eu-pid-jwt-proof-only` | Version 38 |

`status_list_enabled: true` allocates and embeds a Token Status List reference in each issued
credential. The former `broken` fixture toggle is not part of the published beta
contract.

### Issuer protocol surface

| Endpoint | Published requirements | Status |
| --- | --- | --- |
| `GET /issuers/{issuerConfigurationId}/offers/{credentialOfferId}` | Retrieves an offer referenced by a `credential_offer_uri` deeplink. | Supported |
| `POST /issuers/{issuerConfigurationId}/par` | DPoP header and form `response_type=code`, `client_id`, `redirect_uri`, `scope`, `code_challenge`, and `code_challenge_method=S256`; optional `issuer_state` and `state`. Returns `request_uri` and expiry. | Supported |
| `GET /issuers/{issuerConfigurationId}/authorize` | `client_id` and `request_uri`; begins the auto-approved authorization-code flow. | Supported |
| `GET /issuers/{issuerConfigurationId}/redirect` | Chained OAuth callback; redirects to the Wallet with the issuer authorization code. | Supported |
| `POST /issuers/{issuerConfigurationId}/token` | DPoP header plus either a pre-authorized-code or authorization-code form grant; returns a DPoP token and credential nonce. | Supported |
| `POST /issuers/{issuerConfigurationId}/nonce` | Returns `c_nonce` and its expiry. | Supported |
| `POST /issuers/{issuerConfigurationId}/credential` | DPoP access token, DPoP header, and a JSON Credential Request or compact-JWE `application/jwt` request; returns a credential response, optionally as a compact JWE. | Supported |

For the token endpoint, use either the pre-authorized-code grant (with optional
`tx_code`) or the authorization-code grant (with `code`, `code_verifier`, and
`redirect_uri`). The selected issuer metadata determines supported proof and
credential-response encryption capabilities. The service records issuance
capture evidence, but the exact `observed`, `checks`, and event-detail shapes
are intentionally open-ended in the public schema. Inspect an actual session
before a validator depends on a particular field.

### Test-only chained OAuth server

The authorization-code issuance flow uses an internal, auto-approving OAuth
server between Credo and the issuer. It is test-service infrastructure, not a
general-purpose identity provider.

| Endpoint | Capability | Status |
| --- | --- | --- |
| `GET /.well-known/oauth-authorization-server/authorization-servers/{issuerConfigurationId}` | Fake OAuth server metadata. | Supported |
| `GET /authorization-servers/{issuerConfigurationId}/authorize` | Validates Credo's authorization request and immediately redirects with a code. | Supported |
| `POST /authorization-servers/{issuerConfigurationId}/token` | Exchanges the chained code; requires configured client credentials, PKCE verifier, redirect URI, and client ID. | Supported |

## Verifier: OpenID4VP capture sessions

### Session creation

`POST /openid4vp/sessions` creates a session and returns `201` with `session_id`, delivery settings, `request_uri`, `response_uri`, `deeplink`, `authorization_request`, and `status: "created"`.

| Field | Published values / shape | Status |
| --- | --- | --- |
| `scheme` | URL-scheme prefix matching `scheme://`; defaults to `openid4vp://`. | Supported |
| `request_uri_method` | Any string; OpenID4VP defines the case-sensitive values `get` and `post`; defaults to `get`. | Supported by beta, including deliberate malformed values for Wallet negative tests. |
| `client_id_scheme` | `x509_hash`, `x509_san_dns`, `decentralized_identifier`, or `redirect_uri`; defaults to `x509_hash`. | Supported, subject to delivery constraints. |
| `request_delivery` | `by_reference`, `by_value`, or `plain`; defaults to `by_reference`. | Supported |
| `response_type` | `vp_token`, `vp_token id_token`, or `code`; default `vp_token`. | Supported |
| `response_mode` | `direct_post` or `direct_post.jwt`; default `direct_post.jwt`. | Supported |
| `presentation_request` | Open-ended JSON object. | Supported as input; delivery semantics must be inspected. |
| `dcql_query` | Open-ended JSON object. | Supported as input; delivery semantics must be inspected. |
| `scopes` | String or string array. | Supported as input. |
| `transaction_data` | Unconstrained JSON value. | Supported as input. Observed on 02/09/2026: when nested in `presentation_request`, it is preserved in the signed Request Object; supported Wallet types remain unknown. |
| `verifier_info` | Unconstrained JSON value. | Supported as input. Observed on 02/09/2026: when nested in `presentation_request`, it is preserved in the signed Request Object; attestation generation and Wallet support remain unknown. |
| `client_metadata` | Object replacing the generated verifier metadata, or `null` to omit it. | Supported, subject to response-mode constraints. |
| `redirect_uri` | Absolute URI for the Wallet after a successful presentation. | Supported; the service appends a fresh `response_code`. |

`request_uri_method` is valid only with `request_delivery: "by_reference"`.
The verifier preserves a supplied value other than `get` or `post` in the
deeplink for Wallet negative tests. A beta probe on 14/09/2026 accepted
`DELETE`, created session `697125c6-1c20-4e98-81b9-c3b9356a857c`, and returned
the value in the deeplink. Production still returned
`400 {"error":"unsupported_request_uri_method"}` at that time.
`by_value` delivers a signed Request Object in `request`; `plain` delivers
URL-encoded Authorization Request parameters in the deeplink and omits
`request`, `request_uri`, and `request_uri_method`. `client_id_scheme:
"redirect_uri"` requires unsigned `plain` delivery. `x509_san_dns` uses the
verifier certificate DNS Subject Alternative Name; `decentralized_identifier`
uses a separate `did:web` key published at `GET /openid4vp/did.json`.

When `dcql_query` is `null`, Capture omits it from the Wallet-facing request;
the service retains its normal query only as internal verification-session
state, so a response may not validate. With `client_metadata`, an absent field
uses generated metadata, an object replaces it, and `null` omits it. Omission
is supported only with `direct_post`; a `direct_post.jwt` replacement must
retain the generated verifier encryption JWK. After a successful response to a
session with `redirect_uri`, the endpoint returns `{ "redirect_uri": "..." }`
with `Cache-Control: no-store` for the Wallet to open.

Observed on 03/09/2026: a DCQL credential that omits `meta` is accepted at
session creation and preserved without `meta` in the signed Request Object.
An explicitly empty `meta: {}` object is likewise accepted and preserved.
An empty `trusted_authorities: []` array is also accepted and preserved.
On 03/09/2026, a `trusted_authorities` entry with `type: unsupported` and a
string `values` array was accepted and preserved in the signed Request Object.
On 03/09/2026, a `trusted_authorities` entry with a string `values` array but
no `type` was also accepted and preserved in the signed Request Object.
On 07/09/2026, a nested claim path containing a null selector,
`["address", null, "street_address"]`, was accepted and preserved unchanged in
the returned Authorization Request (Capture session `eb00f511-2fc2-4746-baaf-aea6961d3f71`).
On 07/09/2026, the same was observed for the integer-selector path
`["address", "street_address", 0]` (Capture session
`846a712e-19b2-4234-ad85-26bf51b07867`).
On 07/09/2026, Capture also preserved the unsupported Boolean component in
`["address", "street_address", false]` (Capture session
`d3c78b51-3d53-4e03-87f9-84b003e70384`).
On 07/09/2026, Capture preserved the nested missing-member path
`["address", "unavailable_address_member"]` (Capture session
`bb38af67-df85-4654-88b4-56b2699be424`).
On 07/09/2026, Capture accepted and preserved the mdoc path with an absent
namespace, `["org.iso.18013.5.1", "first_name"]`, under
`format: mso_mdoc` and `doctype_value: eu.europa.ec.eudi.pid.1` (Capture
session `13aa1df4-e5b8-432f-b208-5454d71bbea0`).

`additionalProperties` are accepted by the public request schema, but that does **not** establish that an unknown property appears in the signed Authorization Request. Retrieve and decode the request before using an unknown field as test evidence.

### Delivery and response endpoints

| Endpoint | Capability | Status |
| --- | --- | --- |
| `GET /openid4vp/sessions/{sessionId}` | Current presentation capture with `authorization_request`, `observed`, `checks`, `events`, and raw protocol evidence. | Supported |
| `GET /openid4vp/sessions/{sessionId}/deeplink` | Returned deeplink and decoded `authorization_request`. | Supported |
| `GET /openid4vp/sessions/{sessionId}/request` | Retrieves the signed request object as `application/oauth-authz-req+jwt` and marks it as retrieved. | Supported |
| `POST /openid4vp/sessions/{sessionId}/request` | Retrieves the signed request when `request_uri_method: post`; accepts form `wallet_nonce` and additional fields. | Supported |
| `POST /openid4vp/sessions/{sessionId}/response` | Captures a form-encoded Wallet response for that session. | Supported |
| `POST /openid4vp/response` | Alternative form-encoded direct-post endpoint; required `state` identifies the session. | Supported |
| `GET /openid4vp/sessions/{sessionId}/events` | Chronological protocol capture events. | Supported |
| `GET /openid4vp/did.json` | Verifier `did:web` Document used by `client_id_scheme: "decentralized_identifier"`. | Supported |

The session's `raw` object provides the protocol-evidence surface required for
assertions. `raw.authorization_request_jwt` is the exact signed Request Object
returned to the Wallet. `raw.request_uri_http` records the Wallet retrieval
method and redacted headers, and adds `body` only when a body was received.
The capture is attached to the session-specific `/request` endpoint; the
published schema does not expose a separate raw request-target or query field.

The direct-post endpoints return `200` only when the presentation was captured
and verified. A failed verifier check and a Wallet's decision to send no
response are distinct outcomes. `raw.presentation_response_http` and
`raw.presentation_response_verifier_http` provide machine-readable,
sensitive-value-redacted HTTP evidence for valid and invalid responses; the
former retains the exact received body. Inspect the session record and events;
never substitute a screenshot for missing callback evidence.

## Known local limitations

| Area | Finding | Status / source |
| --- | --- | --- |
| Empty `credential_sets[].options` | The reference Android wallet displayed an error but did not POST `error=invalid_request`; the beta session captured only request retrieval. | Blocked for the required protocol assertion. [RI-WALLET-001](REFERENCE-WALLET-ISSUES.md) |
| Positive PID verification | The beta verifier received a `vp_token` but rejected it because the PID issuer URI did not match the issuer certificate SAN. | Verifier-blocked acceptance, not a Wallet failure. [MOCK-VERIFIER-001](REFERENCE-WALLET-ISSUES.md) |
| Invalid `request_uri_method` | Beta preserves arbitrary values in the Wallet-facing deeplink; production still rejects `DELETE` at session creation. | Supported for `WS_RP_MS_ProtocolMessages__152` on beta; production deployment lag remains. |

`pkg/fcaf/MEMORY.md` additionally lists test-specific cases blocked because the public verifier validates malformed DCQL before it can create a signed request, or because it cannot expose the raw request/response feature required by the test. Treat that as coordination state and re-probe it when the service changes.

## Safe use in scenarios

- Use a source scenario under `config_templates/fcaf/wallet_solution/relying_party/scenarios/`; do not edit the generated aggregate pipeline directly.
- Persist the created `session_id` only as a pipeline output needed to fetch protocol evidence. Do not put live session URLs or tokens in fixtures.
- For verifier tests, bind validators to the exact scenario output containing the capture session/request/response. A different scenario's successful response is not fallback evidence.
- For issuer tests, use the session, events, and observed wallet JWKS to prove the relevant issuance exchange; inspect their actual shape first.
- Record a dated probe and update this document when a previously unknown capability becomes a test prerequisite.
