# Dock by Webigram — Version 2.0 technical architecture

**Status:** Proposed, not implemented. **Baseline:** 1.0b. **Target:** 2.0.

## Product definition

Dock is a **lightweight Windows control plane for applications written in PHP, Python, JavaScript/TypeScript, Node.js, Go, and other languages**. Dock itself is written in Go, but applications are launched through their own installed runtimes and managed via processes, HTTP proxy routes, Docker Engine, and local pipelines. **Dock is not the runtime for those applications.**

## Proposed modules

| Module | Recommended design | External requirements |
|---|---|---|
| Core Agent | Small local Go service with loopback control API, independently of GUI | Windows Service integration optional; per-user agent first |
| GUI | Existing Walk application as client of local agent | None |
| Web Gateway | Go net/http reverse proxy + static files + WebSocket upgrades | None for proxy/static content |
| Certificates | Local CA and short-lived leaf certificates using Go crypto/x509 | Trust-store approval for locally trusted HTTPS |
| App Manager | Explicit command templates, process lifecycle, stdout/stderr logs, health probes | App's own runtime (PHP/Python/Node/Go, etc.) |
| Docker Manager | Docker Engine API client and optional Compose CLI | Existing Docker daemon, local or remote |
| Mini CI/CD | Go job scheduler/executor + SQLite history + TOML pipeline specs | Language build/deploy tools used in steps |
| Config | Versioned TOML with migration/backup | None |

## Proposed deployment topology

```text
                    localhost:80 and localhost:443
                                |
                Dock Web Gateway (reverse proxy)
                 /        |         |        \
         api.localhost  app.localhost  shop.localhost static.localhost
              |              |            |             |
       Python :8000     Node/TS :3000   PHP :8001      static files
              |              |            |
          runtime installed / external Docker daemon when needed

   Dock GUI  <--loopback control API--> Dock Agent
                                    |       |       |
                                 Process  Docker  CI/CD
                                 manager  client  + SQLite
```

The gateway routes by **Host/SNI and explicit route configuration**, including HTTP->HTTPS redirects and WebSocket forwarding. Route mapping MUST NOT be tied to the implementation language. Only an adapter/command determines how an application starts.

## HTTP and HTTPS: precise scope

- Default bind **127.0.0.1:80 and 127.0.0.1:443**; optional IPv6 loopback; no LAN/public binding without explicit consent and access controls.
- Ports can be busy or permission-restricted (e.g. IIS/HTTP.sys); detect conflicts and display owning process where possible. Support fallback ports without silently changing configured URLs.
- `*.localhost` is a suitable local-development naming convention. Do not promise arbitrary `.local` hosts resolve automatically; document DNS/hosts needs.
- Forward Host, request path, query, WebSocket upgrades and streaming responses safely; configurable upstream timeouts, request-size limits and health checks.
- Provide **static file hosting** for prebuilt frontend bundles; prevent path traversal and accidental directory listing.
- Process manager supplies upstream servers. PHP support relies on a PHP server or FastCGI bridge; proxy alone does not execute PHP. Python/Node/Go apps need installed runtimes or containers.
- Gateway must stay usable when the GUI closes; therefore separate lifecycle from the visible application.

## Certificate issuer

Support distinct modes:

1. **Self-signed leaf**: low setup cost, browser warning expected.
2. **Local CA**: generate local root CA, issue development certificates containing appropriate SAN entries; user must explicitly approve trusting the root in Windows. Display CA fingerprint and provide revoke/remove-trust instructions.
3. **Future ACME** (optional, not a core requirement): public-domain certificates only where DNS/HTTP challenges can be completed; never describe localhost self-signed certificates as Certbot/Let's Encrypt certificates.

Use Go `crypto/x509`, `crypto/ecdsa` or `crypto/ed25519`, `crypto/tls`, CSPRNG and secure key storage with user-profile ACLs/OS protection. Default to limited local CA scope and short-lived leaf certs; implement rotation, expiry checks and atomic hot reload. Never log private keys or distribute an unprotected CA private key. Import/export should default to certificates only, not secrets.

## Docker Manager (not a Docker Engine)

- Connect to existing local Docker Desktop/WSL2 daemon or **remote Linux Engine** using authenticated TLS or an SSH transport.
- Display containers, images, lifecycle state, published ports, CPU/RAM and logs.
- Actions: Start, Stop, Restart, with explicit confirmation for destructive or disruptive actions. Compose integration can use installed `docker compose` CLI.
- Docker Engine API access is privileged. Do NOT expose unauthenticated TCP 2375 or proxy the daemon to the web gateway.
- No daemon installation bundled; clear disconnected/unavailable state and diagnostic hints.

## Lightweight CI/CD

- No Jenkins/GitHub Actions runner service required: **Dock Agent is the worker**; it executes steps locally with the app's own installed build tools.
- Pipeline definitions in versioned `dock.pipeline.toml` per repo, with bounded, explicit steps for `build -> test -> deploy -> verify`. Make `rollback` a defined optional recovery step, not unconditional.
- Triggers: manual first, then scheduled and Git polling. Webhook receivers must be opt-in with secrets, replay protection and narrow network exposure.
- Persist runs, stages, results, timestamps and indexed log locations to SQLite (single writer, WAL, schema migrations); store large logs in rotated disk files, **not** unlimited database blobs.
- Isolate each job with per-repository locks, working directory and configured environment; never silently run untrusted pipeline definitions from an unknown clone. Require approval for commands, redact secrets in logs, enforce timeouts and cancellation and prevent overlapping deploys.
- Support branch allowlists, retries, concurrency limits, artifact paths, bounded retention and health-check verification. Recovery actions require explicit configuration and clear reporting.
- Closing Dock GUI does not terminate accepted jobs; stopping the agent does. Local polling only works while the agent is available. Remote autonomous CI is out of scope.

### Example (design only, not yet supported)

```toml
version = 1
[project]
name = "example-node-app"
branch = "main"

[triggers]
manual = true
poll_interval_seconds = 300

[[steps]]
name = "build"
command = "npm"
args = ["run", "build"]
timeout_seconds = 600

[[steps]]
name = "test"
command = "npm"
args = ["test"]
timeout_seconds = 600

[[steps]]
name = "verify"
command = "curl"
args = ["--fail", "http://127.0.0.1:3000/health"]
timeout_seconds = 30
```

This example intentionally does not include automatic deployment or a secret token.

## Design principles and quality gates

- Lightweight startup, bounded RAM and CPU, minimal binary count and no unnecessary external services.
- Windows 10/11 compatibility must be tested. Windows 8/8.1 require a separate feasibility/security decision because modern runtimes/toolchains have dropped support.
- Separation of privileges: per-user GUI, loopback authenticated agent control API, explicit elevation only for privileged actions (root CA trust, service setup, firewall/port binding where required).
- Never implicitly bind management APIs beyond loopback; authenticate control requests and protect from cross-origin requests.
- Never store GitHub PATs, ACME account keys, Docker secrets or pipeline environment secrets in plaintext TOML.
- Atomic config migrations, backups, rollbacks; old 1.0b TOML remains supported.
- Network and Git operations must be cancellable; preserve local changes on failure.
- A minimal end-to-end Windows test matrix is required before calling any of these modules delivered.

## Delivery sequence

1. **Foundation:** service architecture, config migrations, security model, tests.
2. **Web:** gateway, routes, ports, static assets, TLS issuer and CA management.
3. **Apps:** runtime/process adapters and health/preview UX.
4. **Docker:** daemon connection and safe container operations.
5. **CI/CD:** pipeline schema and executor, SQLite, triggers, logs, recovery.
6. **Workspace:** Git status/preview/history, editor integration, settings and backup.
7. **Release:** end-to-end testing, documentation, packaging and supported OS matrix.

See [NEXT-VERSION-2.0.md](NEXT-VERSION-2.0.md) and [V2-ISSUES.md](V2-ISSUES.md).
