# Dock by Webigram

**Tagline:** Pull. Run. Preview.  
**Product type:** Portable Windows developer utility  
**Status:** Product concept / requirements (implementation not started)  
**Developed by:** Webigram team  
**Repository:** https://github.com/amdkna/tool-project-dock

## Overview

Project Dock (short name: **Dock**) is a lightweight, portable Windows tool for bringing a project from a Git repository into a local workspace, running it in a development or test environment, and viewing its result in the browser.

The primary goal is to simplify the loop: **Clone/Pull → Build/Run → Preview → Modify → Commit/Push**.

It is not intended to replace a full IDE, Git client, CI/CD system, or production hosting platform. It provides a focused interface for previewing and testing projects locally.

## Core requirements

### 1. Repository management

- Add a project by Git repository URL and selected branch.
- Clone into an isolated workspace if it does not exist locally.
- **Pull** changes from the configured remote branch.
- **Push** local commits to the configured branch.
- Show current branch, working-tree status, and Git operation logs.
- Detect local modifications, merge conflicts, divergent histories, and authentication errors.
- Never silently discard changes, force-push, or overwrite a dirty worktree.
- A push requires commits: show uncommitted changes and provide a safe commit workflow or explain why a push has nothing to send.
- Confirm before potentially destructive or remote-changing operations.

### 2. Local development and preview

- Start, stop, and restart an application locally.
- Serve static build output with Go's built-in HTTP server.
- For development servers or backend projects, launch the configured command as a child process and optionally reverse-proxy it.
- Open the local preview in the default browser.
- Show listening URL, assigned port, process status, and stdout/stderr logs.
- Support build commands and configure the output folder.
- Do not assume every repository is a static site: Node.js, Python, PHP, databases, etc. may require separate runtimes and services.
- Docker is optional, not a mandatory dependency.

### 3. Portable Windows operation

- Target Windows with a compact distribution and no mandatory installer.
- Keep settings and project workspace paths relative to the executable where practical.
- Run multiple instances independently, potentially with different repositories.
- Avoid fixed ports: allocate available loopback ports and pass chosen ports to child processes correctly.
- Isolate per-instance workspace, process, ports, settings, and logs.
- Bind preview servers to `127.0.0.1` by default.
- On exit, gracefully terminate managed child processes without terminating unrelated applications.

### 4. User interface

Suggested basic screen:

- Project name, repository, branch, and local workspace.
- Buttons: **Clone / Pull / Push / Start / Stop / Restart / Open Preview**.
- Status: Git sync status, process state, local preview URL and port.
- Scrollable live console for Git, build and runtime logs.
- Configuration editing and actionable error messages.

## Proposed technology stack

| Component | Suggested choice | Notes |
|---|---|---|
| Application language | **Go** | Lightweight native compilation, good process and HTTP support |
| Windows GUI | **Walk** (initial candidate) | Native Windows widgets; evaluate maintenance and compatibility before implementation |
| Alternative GUI | Go + WebView | Consider if native GUI limitations become significant |
| Static HTTP server | Go `net/http` | No separate Nginx needed for basic static preview |
| Reverse proxy | Go `httputil.ReverseProxy` | Optional for backend/dev-server projects |
| Git | Git CLI | Use installed Git or optionally bundle PortableGit |
| Configuration | YAML | Human-readable `config.yaml` |
| Runtime management | Go `os/exec` | Capture output; lifecycle management on Windows |

**Technical direction:** Go is the preferred language. Walk, WebView, and Git distribution strategy remain implementation decisions to validate. Nginx is not required for the MVP.

## Draft configuration

The following illustrates the intended configuration model, not an implemented schema.

```yaml
project:
  name: "My Project"
  repository: "https://github.com/example/my-project.git"
  branch: "main"
  directory: "./workspace"

git:
  auto_pull: false
  confirm_before_push: true

server:
  host: "127.0.0.1"
  port: 0
  mode: "dev"
  command: "npm run dev"
  build_command: "npm run build"
  output: "./dist"

app:
  auto_open_browser: true
  show_logs: true
```

**Port note:** `port: 0` directly selects an available port for a Go listener. A child development server must be explicitly configured with a selected port and checked for readiness; merely setting `port: 0` in YAML does not automatically configure external runtimes.

**Security note:** Do not store Git passwords, access tokens, or other secrets in plaintext YAML. Use Git Credential Manager, SSH authentication, or other appropriate secure mechanisms.

## Suggested architecture

```text
Dock.exe
  ├── Windows UI
  ├── Configuration Manager (YAML)
  ├── Git Manager (Git CLI)
  ├── Project Manager
  │    ├── Clone / Pull / Push
  │    ├── Detect project type
  │    └── Build / Start / Stop
  ├── Preview Manager
  │    ├── Go static HTTP server
  │    └── Reverse proxy (when required)
  ├── Port & Process Manager
  └── Live Logs
         │
         └── Isolated local project workspace
```

## Automatic project detection (proposed enhancement)

On import, inspect files to suggest execution settings:

- `package.json` / `vite.config.*`: Node.js / Vite scripts.
- `go.mod`: Go project.
- `requirements.txt` / `pyproject.toml`: Python project.
- `composer.json`: PHP project.
- `Dockerfile` / `compose.yaml`: optional container-based runtime.
- `index.html` or configured `dist` folder: static-site preview.

Detection should **suggest commands, not execute untrusted scripts automatically**. Users should review and approve commands, particularly when importing unknown repositories.

## MVP scope (recommended)

1. Portable Windows application with a basic GUI.
2. Configure a repository and clone it into a local workspace.
3. Pull, inspect changes and push existing commits safely.
4. Start/stop a configured local command and show its live logs.
5. Serve a static site or open an externally managed dev-server URL.
6. Open preview in the browser.
7. Choose non-conflicting local ports and support multiple instances.
8. Read/write a portable YAML configuration.

## Later enhancements

- Automatic project detection and one-click approved setup.
- Built-in Git diff, staging and commit UI.
- Multi-project management within one Dock window.
- Reverse proxy, custom hostnames and optional local TLS.
- Runtime dependency checks and guided installation.
- Optional Docker Compose integration.
- Saved templates for common web frameworks.
- Improved build diagnostics and health checks.

## Naming and brand

Product name: **Dock**  
Full name: **Dock by Webigram**  
Repository name: `tool-project-dock`  
Primary tagline: **Pull. Run. Preview.**  
Alternative positioning: **Your local development playground.**

Other names considered during concept development: Webigram Stage, Webigram Lab, Webigram Launchpad, Webigram Preview, Webigram Port, Webigram Forge, Webigram Nest.

## Decisions vs proposals

**User requirements:** Windows, extremely lightweight, portable, multiple simultaneous instances, configuration file, Git pull/push, local web preview for development/testing, Webigram team branding.

**Recommended but not yet verified/implemented:** Go, Walk GUI, YAML schema, Go HTTP server, Git CLI/PortableGit, automatic project detection, phased MVP features.

This README records the product discussion and initial technical proposal. It does not claim that the application or any listed feature has already been built.
