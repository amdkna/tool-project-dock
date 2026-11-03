# Product concept and roadmap

Dock by Webigram is a lightweight, portable Windows developer tool for moving Git projects to a local test/development environment. Product tagline: **Pull. Run. Preview.**

## Original requirements
- Desktop Windows GUI, portable and very lightweight.
- Git PULL and PUSH buttons.
- Local web server, initially suggested Nginx, for browser-based previews.
- Configuration file.
- Multiple independently running instances.
- Documentation and release changelog.

## Chosen approach
- Go for the compiled application.
- Walk for initial native Windows GUI.
- TOML for human-editable project configuration.
- Installed Git CLI initially; possible PortableGit bundle later.
- Go `net/http` as a future static web server instead of requiring Nginx.
- Separate workspace, local port and process lifecycle for each instance.
- Future project detection for Vite, Node, Python, Go, PHP and Docker; do not automatically execute unknown scripts.

## Milestones
1. **v0.1.1 delivered as source:** basic GUI, File/Info menus, config launcher, Git PULL/PUSH and command result log.
2. Git Clone, project selection, branch/status UX and authentication guidance.
3. Build/run/stop, preview URL, local HTTP server, auto port allocation.
4. Reliable isolation for multiple app instances and process clean-up.
5. Auto-detection of project stack and reusable project templates.

The tool currently manages an **existing** Git clone. Preview web hosting has not been implemented yet.
