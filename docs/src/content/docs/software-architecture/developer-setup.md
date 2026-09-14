---
title: "🤖 Developer Setup"
description: ""
sidebar:
  order: 7
---
## **Prerequisites**

Before you begin, ensure you have the following tools installed:

-   [Git](https://git-scm.com/downloads)
-   [Make](https://www.gnu.org/software/make)
-   [Mise](https://mise.jdx.dev/getting-started.html)
-   [Temporal](https://docs.temporal.io/cli)
-   [Tmux](https://github.com/tmux/tmux/wiki/Installing)
-   [Pre-commit](https://pre-commit.com/)
-   [Golang](https://go.dev/doc/install)

### **Install `slangroom-exec`**

Download the appropriate `slangroom-exec` binary for your OS from the [releases page](https://github.com/dyne/slangroom-exec/releases).

Add `slangroom-exec` to PATH and make it executable:

```bash
wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-Linux-x86_64 -O slangroom-exec
chmod +x slangroom-exec
sudo cp slangroom-exec /usr/local/bin/
```
### **Install `Pre-commit`**
For most Linux distrubution, just do: 
```bash
sudo apt install pre-commit
```


### **Install `slangroom-exec`**

## **Setup Workspace**

### **Clone the repository**

```bash
git clone https://github.com/ForkbombEu/credimi
```

### **Install dependencies**

```bash
cd credimi
mise trust
make credimi
```

## Edit your .env file 

Copy .env.example to .env

```bash
cp ./webapp/.env.example ./webapp/.env  
```

Then edit the .env file, particularly: 

1. set the absolute path in ROOT_DIR
1. Get a token from https://www.certification.openid.net and add it in TOKEN

Copy ./webapp/env.example to 

## **Start Development Server**

```bash
make dev
```

> [!TIP]
> Use `make help` to see all the commands available.

## Temporal Visibility Search Attributes

Pipeline listings rely on a Temporal visibility search attribute named `PipelineIdentifier` (type `Keyword`).
Register it once per Temporal cluster:

```bash
temporal operator search-attributes create --name PipelineIdentifier --type Keyword
```

If the attribute is added after workflows already exist, trigger a Temporal visibility reindex to backfill
historical data (see Temporal admin tooling docs for your deployment).

## Emulating a Cloudflare-proxied deployment

Production Credimi is served through a Cloudflare proxy that challenges non-browser clients. Temporal
activities use a plain HTTP client, so any callback sent to the public `app_url` is answered with a
managed-challenge interstitial instead of the API response, and the activity fails.

`scripts/waf-emulator` reproduces that split locally on the host:

- browser-like requests (an `Accept: text/html` document navigation with a real browser user agent) are proxied to the origin,
- every other client — Go's `net/http`, curl, server SDKs — receives `403` with `cf-mitigated: challenge`.

```bash
make waf-emulator   # listens on :8091, proxies to http://localhost:8090
```

Use `http://localhost:8091` as the public App URL when testing browser links from the host. The
production Compose file intentionally has no host-gateway aliases; add those only in a local override
if a containerized worker must reach the emulator.

To exercise the callback fix, keep the internal URL pointed at the origin:

```text
Settings → App URL          http://localhost:8091
CREDIMI_INTERNAL_APP_URL    http://credimi:8090
```

Callbacks then use the origin and avoid the emulated challenge. Clearing `CREDIMI_INTERNAL_APP_URL`
sends them through the emulated proxy and the run fails with `Unexpected HTTP status code: expected
200, got 403`, matching the production failure mode.

## mobile device semaphore Ops (Internal)

### Defaults and knobs

- Default acquire wait timeout: 45m.
- Override timeout: `MOBILE_DEVICE_SEMAPHORE_WAIT_TIMEOUT=30m` (or any valid `time.ParseDuration` value).
- Disable semaphore (no-op acquire/release): `MOBILE_DEVICE_SEMAPHORE_DISABLED=1`.
- Internal Temporal API auth key: `CREDIMI_INTERNAL_ADMIN_KEY=<plaintext key>` (required for internal workflow HTTP calls).

### Internal admin API key rollout

1. Run DB migrations so `api_keys` supports `key_type`, `superuser`, `revoked`, `expires_at`.
2. Provision an internal admin key (hash stored in DB, plaintext returned once).
3. Set `CREDIMI_INTERNAL_ADMIN_KEY` in backend and worker runtime secrets.
4. Deploy backend/workers with middleware + internal HTTP activity enabled.
5. Smoke-check one authenticated user route and one internal Temporal route.

### Emergency procedures

- Semaphore workflows live in the Temporal `default` namespace with IDs: `mobile-device-semaphore/<device_id>`.
- Query current state via Temporal UI (`GetState`) or `GET /api/mobile-runner/semaphore?device_identifier=...`.
- To unstick a runner, terminate the semaphore workflow in Temporal; it will be recreated on the next acquire.
