# Git worktree utilities for parallel local and agent workflows

Research date: 15 September 2026. Repository activity and star counts are point-in-time
signals from GitHub's official API, not quality guarantees. Feature claims below cite official
repositories, versioned READMEs, or first-party documentation.

## Executive recommendation

1. **Prefer [Worktrunk](https://github.com/max-sixty/worktrunk) as the optional
   developer-facing worktree manager.** It is the strongest reviewed combination of
   create/switch/remove UX, parallel-agent support, blocking and background lifecycle hooks,
   approved project configuration, copy-on-write transfer of ignored files, deterministic port
   templates, per-branch state, and process teardown
   ([overview](https://worktrunk.dev/),
   [hooks](https://worktrunk.dev/hook/#hook-types),
   [copy-ignored](https://worktrunk.dev/step/#wt-step-copy-ignored),
   [per-worktree services](https://worktrunk.dev/tips-patterns/#per-worktree-services)).
   It also has the strongest adoption signal in this focused set: 7,747 stars and a push on
   15 September 2026 ([GitHub API](https://api.github.com/repos/max-sixty/worktrunk)), with
   v0.77.0 released on 8 September 2026
   ([release](https://github.com/max-sixty/worktrunk/releases/tag/v0.77.0)).

2. **Pilot [Workz](https://github.com/rohansx/workz) only as an optional runtime companion if the
   team wants a maintained port registry and provisioning layer instead of repository scripts.**
   Released v0.11.0 copies env files, syncs dependencies, allocates per-repository port ranges,
   derives database and Compose names, and exposes `post_start`/`pre_done` hooks
   ([v0.11.0 README](https://github.com/rohansx/workz/blob/v0.11.0/README.md)).
   It integrates with Worktrunk, Cursor, Claude Code, or native Git through
   `workz sync --isolated`
   ([hook usage](https://github.com/rohansx/workz/blob/v0.11.0/README.md#use-it-as-a-hook)).
   Its adoption is much smaller—87 stars, last push 5 September 2026
   ([GitHub API](https://api.github.com/repos/rohansx/workz))—and current `main` documents later
   named-service and reflink behavior absent from the latest v0.11.0 release
   ([current README](https://github.com/rohansx/workz#monorepos--named-services),
   [v0.11.0 source](https://github.com/rohansx/workz/tree/v0.11.0)). Pin and test a release rather
   than relying on unreleased documentation.

Credimi should not require either executable to build or test. The durable repository contract
should be tool-neutral setup/teardown scripts plus parameterized Makefile/Compose ports. Worktrunk,
Workz, Cursor, or native `git worktree` can call the same scripts. Developers who specifically want
a tmux/Zellij/Kitty/WezTerm agent cockpit should separately consider workmux
([configuration](https://workmux.raine.dev/guide/configuration),
[port recipe](https://workmux.raine.dev/guide/monorepos#port-isolation)).

## Comparison table

| Candidate | Problem solved | Lifecycle, env/cache, ports | Maturity on 15/09/2026 | Credimi fit |
|---|---|---|---|---|
| **[Worktrunk](https://github.com/max-sixty/worktrunk)** | Branch-like create/switch/list/merge/remove for parallel agents ([overview](https://worktrunk.dev/)). | Blocking `pre-*` and background `post-*` hooks for switch/create/commit/merge/remove; approved project hooks; `.worktreeinclude`; APFS/btrfs/XFS/ReFS reflink copy; `hash_port`; branch vars; tethered processes ([hooks](https://worktrunk.dev/hook/), [copy](https://worktrunk.dev/step/#wt-step-copy-ignored), [services](https://worktrunk.dev/tips-patterns/#per-worktree-services)). | 7,747 stars; pushed 15/09; v0.77.0 on 08/09 ([API](https://api.github.com/repos/max-sixty/worktrunk), [release](https://github.com/max-sixty/worktrunk/releases/tag/v0.77.0)). Best documentation in the set ([docs](https://worktrunk.dev/)). | **Best overall.** No bare-repo migration; exact hooks and cache-seeding features. A collision check must supplement hashed ports. |
| **[Workz](https://github.com/rohansx/workz)** | Makes worktrees created by any host runnable ([v0.11 README](https://github.com/rohansx/workz/blob/v0.11.0/README.md)). | Copies env, symlinks or copies dependencies, auto-installs Bun and peers, registers port ranges, derives DB/Compose names, setup/teardown, cleanup and doctor ([isolation](https://github.com/rohansx/workz/blob/v0.11.0/README.md#environment-isolation), [sync](https://github.com/rohansx/workz/blob/v0.11.0/README.md#what-gets-synced)). | 87 stars; pushed 05/09; v0.11.0 on 02/07 ([API](https://api.github.com/repos/rohansx/workz), [release](https://github.com/rohansx/workz/releases/tag/v0.11.0)). | **Best turnkey runtime pilot.** Port range and Compose intent fit; released v0.11 lacks current-main named service and reflink features, and PocketBase SQLite remains project-specific. |
| **[workmux](https://github.com/raine/workmux)** | Worktrees plus multiplexer panes and multiple agent sessions ([README](https://github.com/raine/workmux#readme)). | Copy/symlink globs; `post_create`, aborting `pre_merge` and `pre_remove`; nested monorepo config; first-party checked/incremented port script ([configuration](https://workmux.raine.dev/guide/configuration), [monorepos](https://workmux.raine.dev/guide/monorepos)). | 2,627 stars; pushed 14/09; v0.1.262 on 12/09 ([API](https://api.github.com/repos/raine/workmux), [release](https://github.com/raine/workmux/releases/tag/v0.1.262)). | **Strong conditional choice.** Excellent agent visibility, but terminal-multiplexer coupling is unnecessary for Cursor-only users. |
| **[`gtr`](https://github.com/coderabbitai/git-worktree-runner)** | Portable repository-scoped worktrees for PRs, editors, and agents ([README](https://github.com/coderabbitai/git-worktree-runner#readme)). | Selective file/dir copying, trusted `.gtrconfig`, `postCreate`, `preRemove`, `postRemove`, and shell `postCd`; no reflink cache or port allocator documented ([configuration](https://github.com/coderabbitai/git-worktree-runner/blob/main/docs/configuration.md#hooks)). | 1,776 stars; pushed 15/09; v2.11.1 on 14/09 ([API](https://api.github.com/repos/coderabbitai/git-worktree-runner), [release](https://github.com/coderabbitai/git-worktree-runner/releases/tag/v2.11.1)). | **Best mature conventional peer.** Credimi must still own ports, Compose isolation, and cache policy. |
| **[`wtp`](https://github.com/satococoa/wtp)** | Thin Go wrapper with automatic setup ([README](https://github.com/satococoa/wtp#readme)). | `post_create` has explicit copy, symlink, and command actions; copies ignored `.env` and can symlink `.bin`; no teardown or port primitive documented ([configuration](https://github.com/satococoa/wtp#configuration)). | 629 stars; pushed 30/03; v2.10.3 on 08/03 ([API](https://api.github.com/repos/satococoa/wtp), [release](https://github.com/satococoa/wtp/releases/tag/v2.10.3)). | **Good minimal Go-native fallback.** Easy env copies, weaker runtime cleanup. |
| **[`k1LoW/git-wt`](https://github.com/k1LoW/git-wt)** | Small Git subcommand for create/switch/delete/rename without a special repo layout ([README](https://github.com/k1LoW/git-wt#readme)). | Copies ignored/untracked/modified or selected paths; directory symlinks; create/delete hooks; no native ports or reflinks ([configuration](https://github.com/k1LoW/git-wt#configuration)). | 560 stars; pushed and released v0.29.3 on 11/09 ([API](https://api.github.com/repos/k1LoW/git-wt), [release](https://github.com/k1LoW/git-wt/releases/tag/v0.29.3)). | **Best small wrapper.** All runtime isolation remains project scripting. |
| **[`ahmedelgabri/git-wt`](https://github.com/ahmedelgabri/git-wt)** | Interactive management, diagnostics, migration, and cleanup around `.bare/` ([README](https://github.com/ahmedelgabri/git-wt#readme)). | Before/after add and remove hooks with context; copying, cache, and ports require custom hooks ([hooks](https://github.com/ahmedelgabri/git-wt#hooks)). | 78 stars; pushed 08/09; v2.1.0 on 28/07 ([API](https://api.github.com/repos/ahmedelgabri/git-wt), [release](https://github.com/ahmedelgabri/git-wt/releases/tag/v2.1.0)). | **Capable but wrong default shape.** Bare migration adds risk without solving ports. |
| **[`bkildow/wt-cli`](https://github.com/bkildow/wt-cli)** | Managed bare-repo project with shared assets and setup/teardown ([README](https://github.com/bkildow/wt-cli#readme)). | `shared/copy`, `shared/symlink`, APFS/btrfs/XFS reflinks, templates, serial/parallel setup and teardown, Claude hooks; no allocator ([configuration](https://github.com/bkildow/wt-cli#configuration)). | 9 stars; pushed 09/09; v0.10.0 on 08/09 ([API](https://api.github.com/repos/bkildow/wt-cli), [release](https://github.com/bkildow/wt-cli/releases/tag/v0.10.0)). | **Feature-rich experiment.** Exact cache/env mechanisms, but tiny adoption and mandatory managed layout. |
| **[`mvwi/wt`](https://github.com/mvwi/wt)** | End-to-end branch/PR workflow and automatic initialization ([README](https://github.com/mvwi/wt#readme)). | Auto-copies common env files and detects Bun/Go install commands; no teardown, reflink, or port primitive documented ([configuration](https://github.com/mvwi/wt#configuration)). | 2 stars; pushed and released v1.21.0 on 20/05 ([API](https://api.github.com/repos/mvwi/wt), [release](https://github.com/mvwi/wt/releases/tag/v1.21.0)). | **Do not standardize.** Tiny adoption and bundled submit/rebase/push opinions exceed scope. |
| **[sebasv/grove](https://github.com/sebasv/grove)** | TUI for worktrees, persistent terminals, diffs, PR/CI, and agent attention ([README](https://github.com/sebasv/grove#readme)). | Creates/removes worktrees but documents no configurable hooks, env copy, cache sharing, or port isolation ([config](https://github.com/sebasv/grove#power-user-config)). | 11 stars; pushed 04/05; v1.0 on 25/04 ([API](https://api.github.com/repos/sebasv/grove), [release](https://github.com/sebasv/grove/releases/tag/v1.0)). | **Viewer, not bootstrap infrastructure.** |
| **[`lost-in-the/grove`](https://github.com/lost-in-the/grove)**, a different project | Worktree + tmux manager with setup and Docker control ([v0.10 README](https://github.com/lost-in-the/grove/blob/v0.10.0/README.md)). | Copy/symlink/command lifecycle actions and slot-offset isolated Compose stacks ([hooks](https://github.com/lost-in-the/grove/blob/v0.10.0/README.md#lifecycle-hooks), [Docker](https://github.com/lost-in-the/grove/blob/v0.10.0/README.md#docker-integration)). | 3 stars; pushed 13/09; v0.10.0 on 22/07; officially beta ([API](https://api.github.com/repos/lost-in-the/grove), [release](https://github.com/lost-in-the/grove/releases/tag/v0.10.0)). | **Feature-relevant, too immature.** Do not confuse it with the requested Grove. |
| **[`johnlindquist/worktree-cli`](https://github.com/johnlindquist/worktree-cli)** | Cursor-focused create/open/PR/merge with trusted setup scripts ([README](https://github.com/johnlindquist/worktree-cli#readme)). | `worktrees.json` or `.cursor/worktrees.json`, `$ROOT_WORKTREE_PATH`, and dependency install flags; no remove hook, cache sharing, or ports documented ([setup](https://github.com/johnlindquist/worktree-cli#setup-worktree-configuration)). | 159 stars; pushed and released v2.14.0 on 09/12/2025 ([API](https://api.github.com/repos/johnlindquist/worktree-cli), [release](https://github.com/johnlindquist/worktree-cli/releases/tag/v2.14.0)). | **Cursor-friendly but stale relative to peers.** |
| **[Falq](https://github.com/sunduq-ai/falq)**, provisioner | Declarative runtime isolation for worktrees created by any manager ([README](https://github.com/sunduq-ai/falq#readme)). | Slot port ranges, nested env copy/upsert, config generation/patching, hooks, unique Compose projects, bind checks, Compose lint; no daemon ([manifest](https://github.com/sunduq-ai/falq#the-manifest-in-20-lines), [Compose](https://github.com/sunduq-ai/falq#docker-compose)). | 4 stars; pushed and released v0.6.1 on 21/06 ([API](https://api.github.com/repos/sunduq-ai/falq), [release](https://github.com/sunduq-ai/falq/releases/tag/v0.6.1)). | **Technically close, operationally premature.** Its linter flags Credimi's current fixed ports and explicit container names. |
| **[Coasts](https://github.com/coast-guard/coasts)**, runtime companion | Isolates complete local runtimes for parallel worktrees, especially Compose stacks ([docs](https://coasts.dev/docs)). | Dynamic ports, canonical-port checkout, isolated containers/volumes, Compose boot, daemon and UI; assumes worktrees already exist ([overview](https://coasts.dev/docs#why-coasts-for-worktrees)). | 429 stars; pushed 28/04; v0.1.53 on 20/04 ([API](https://api.github.com/repos/coast-guard/coasts), [release](https://github.com/coast-guard/coasts/releases/tag/v0.1.53)). | **Phase-two companion.** Comprehensive but substantially heavier than Compose parameterization. |

## Feature details with primary sources

### Worktrunk

- The hook matrix covers pre/post switch, start, commit, merge, and remove. Pre-hooks block and
  can abort; post-hooks run in the background with logs
  ([hook types](https://worktrunk.dev/hook/#hook-types)).
- Checked-in project hooks require first-run approval, and changed commands require approval again
  ([security](https://worktrunk.dev/hook/#security)).
- `wt step copy-ignored` uses Git's ignore rules, skips existing destination files unless forced,
  supports `.worktreeinclude`, and can require the allowlist before copying anything
  ([selection](https://worktrunk.dev/step/#what-gets-copied)).
- On APFS it reflinks per file, sharing disk blocks until either copy is written; unsupported
  filesystems receive a full copy
  ([copy-on-write](https://worktrunk.dev/step/#copy-on-write)).
- `hash_port` yields deterministic branch-derived candidates; labels can derive separate API, UI,
  Temporal, and Temporal UI candidates
  ([dev server](https://worktrunk.dev/tips-patterns/#dev-server-per-worktree),
  [database example](https://worktrunk.dev/tips-patterns/#database-per-worktree)).
  A finite hash range cannot guarantee uniqueness, so Credimi should check candidate availability
  and persist the result, as workmux's first-party script does
  ([checked allocation](https://workmux.raine.dev/guide/monorepos#port-isolation)).
- `wt step tether` terminates the process group when its worktree disappears, even if a normal
  remove hook is bypassed ([tether](https://worktrunk.dev/step/#wt-step-tether)).

### Workz and workmux

- Workz v0.11.0 tracks ranges in `~/.config/workz/ports.json`, writes managed isolation values
  without overwriting the rest of `.env.local`, derives a database/Compose project, and releases
  resources through `workz done`
  ([v0.11 isolation](https://github.com/rohansx/workz/blob/v0.11.0/README.md#environment-isolation)).
- Workz v0.11 symlinks Node dependency trees by default but permits a full-copy override
  ([v0.11 sync](https://github.com/rohansx/workz/blob/v0.11.0/README.md#what-gets-synced)).
  Current `main` adds a `clone` mode and named monorepo service ports, but those must not be assumed
  for v0.11
  ([current sync](https://github.com/rohansx/workz#what-gets-synced),
  [current services](https://github.com/rohansx/workz#monorepos--named-services)).
- Workmux copies or symlinks globs and exposes `post_create`, `pre_merge`, and `pre_remove`
  ([file operations](https://workmux.raine.dev/guide/configuration#file-operations),
  [lifecycle](https://workmux.raine.dev/guide/configuration#lifecycle-hooks)).
  Its main differentiator is persistent multiplexer panes and built-in agent command/status support
  ([panes](https://workmux.raine.dev/guide/configuration#panes)).

### Simpler wrappers

- `gtr` has the strongest maturity signal among conventional wrappers and a good trust model, but
  it only copies paths and invokes hooks; service identity and port allocation remain project code
  ([configuration](https://github.com/coderabbitai/git-worktree-runner/blob/main/docs/configuration.md),
  [trust](https://github.com/coderabbitai/git-worktree-runner#git-gtr-trust)).
- `wtp` has explicit copy/symlink/command actions, but its documented lifecycle is creation-only
  ([configuration](https://github.com/satococoa/wtp#configuration)).
- `k1LoW/git-wt` is the strongest minimal Git-subcommand design. Its docs explicitly warn that
  symlinked directories are shared, so an install in one worktree changes every worktree
  ([`wt.symlink`](https://github.com/k1LoW/git-wt#wtsymlink----symlink)).
- `ahmedelgabri/git-wt` has a complete four-edge add/remove lifecycle, but requires its bare
  topology and leaves copy/cache/ports to shell hooks
  ([layout and hooks](https://github.com/ahmedelgabri/git-wt#hooks)).
- `bkildow/wt-cli` has copied/symlinked shared trees, templates, reflinks, parallel
  setup/teardown, and Claude hooks, but only nine stars
  ([README](https://github.com/bkildow/wt-cli#readme),
  [API](https://api.github.com/repos/bkildow/wt-cli)).

### Runtime-specific companions

Falq is a daemonless declarative counterpart to Workz: a checked-in manifest assigns slots and
contiguous port ranges, copies nested env files, upserts values, runs provision/teardown hooks,
sets a unique Compose project, checks ports, and tears resources down
([v0.6.1 README](https://github.com/sunduq-ai/falq/blob/v0.6.1/README.md)).
Its Compose doctor flags `container_name`, fixed host ports, host networking, external networks,
and top-level names as isolation breakers
([Compose](https://github.com/sunduq-ai/falq/blob/v0.6.1/README.md#docker-compose)).

Coasts goes further: each worktree runtime gets dynamic ports, containers, networks, and data, and
one runtime can be checked out onto canonical ports
([Coasts docs](https://coasts.dev/docs)). It requires a daemon, Docker, Node.js, and `socat`
([requirements](https://coasts.dev/docs#requirements)), so it should be evaluated only if the
lighter Compose-native route remains painful.

## Anti-recommendations

- **Do not make a third-party wrapper a build/test dependency.** Native Git already defines
  worktree creation, movement, locking, repair, and pruning
  ([official reference](https://git-scm.com/docs/git-worktree)). Credimi should own only
  tool-neutral setup/runtime commands that wrappers may invoke.
- **Do not mandate a bare-repository conversion.** The bare-first layouts of
  `ahmedelgabri/git-wt`, `bkildow/wt-cli`, and `jarredkenny/worktree-manager` alter repository
  topology without solving Credimi's ports
  ([ahmedelgabri layout](https://github.com/ahmedelgabri/git-wt#clone-with-the-bare-worktree-layout),
  [bkildow layout](https://github.com/bkildow/wt-cli#how-it-works),
  [jarredkenny README](https://github.com/jarredkenny/worktree-manager#readme)).
- **Do not symlink `webapp/node_modules` across active branches.** Installs mutate the shared tree
  for every worktree
  ([k1LoW warning](https://github.com/k1LoW/git-wt#wtsymlink----symlink)).
  Prefer an isolated APFS reflink or plain `bun install`; Bun already uses a global cache and
  macOS `clonefile`
  ([Bun cache](https://bun.sh/docs/install/cache)).
- **Do not enable Workz's default shared `node_modules` behavior.** Configure
  `node_modules = "copy"` on v0.11 or disable that sync and use Bun
  ([v0.11 policy](https://github.com/rohansx/workz/blob/v0.11.0/README.md#what-gets-synced)).
  Do not configure unreleased `clone` behavior merely because it appears on `main`.
- **Do not copy every ignored path.** Worktrunk does so by default except for its built-in denylist
  ([copy defaults](https://worktrunk.dev/step/#what-gets-copied)); Credimi ignores secrets,
  `pb_data`, `.bin`, coverage, build output, and dependencies
  ([Credimi `.gitignore`](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/.gitignore)).
  Require a `.worktreeinclude` allowlist.
- **Do not copy `pb_data` automatically.** PocketBase uses SQLite-backed local data
  ([Credimi runtime](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/AGENTS.md#dev-runtime));
  SQLite documents locking/consistency requirements and its Online Backup API
  ([SQLite Backup API](https://www.sqlite.org/backup.html)).
- **Do not standardize Falq yet.** It has four stars and `0.6.x` maturity
  ([API](https://api.github.com/repos/sunduq-ai/falq)); its volume-removing teardown and rejection
  of explicit container names need deliberate reconciliation with Credimi
  ([Compose behavior](https://github.com/sunduq-ai/falq/blob/v0.6.1/README.md#docker-compose)).
- **Do not confuse the two Grove projects.** `sebasv/grove` is a viewer/terminal surface
  ([README](https://github.com/sebasv/grove#readme)); `lost-in-the/grove` adds lifecycle and
  isolated Compose but explicitly labels itself beta and has three stars
  ([v0.10 README](https://github.com/lost-in-the/grove/blob/v0.10.0/README.md),
  [API](https://api.github.com/repos/lost-in-the/grove)).
- **Do not adopt `copy-env` or `gitcrew` without an exact repository and renewed review.** Search
  found no material standalone `copy-env` worktree utility
  ([search](https://github.com/search?q=copy-env+worktree&type=repositories)); exact `gitcrew`
  results had zero stars
  ([search](https://github.com/search?q=gitcrew+in%3Aname&type=repositories),
  [agent-focused result](https://api.github.com/repos/dalenguyen/gitcrew)).
- **Do not direct agents to Worktrunk merge/commit or `mvwi/wt submit` by default.** Those commands
  bundle commit, merge, or push behavior
  ([Worktrunk workflow](https://worktrunk.dev/#quick-start),
  [mvwi commands](https://github.com/mvwi/wt#commands)), while Credimi requires explicit
  authorization before commit or push
  ([policy](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/AGENTS.md#git-and-commit-rules)).

## Credimi-specific implications

### Existing behavior to preserve

- The Makefile already derives `COMPOSE_PROJECT_NAME` from the checkout root basename and uses a
  project-specific override path
  ([source](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/Makefile#L3-L9)).
  Docker documents project names as the isolation key for multiple feature environments
  ([Docker docs](https://docs.docker.com/compose/how-tos/project-name/)).
- Compose still publishes fixed host ports 7233, 8280, and 8090
  ([source](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/docker-compose.yaml)),
  while API/UI processes use 7233, 8090, and 5100
  ([Procfile](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/Procfile.dev)).
  A unique Compose project does not remove host-port collisions.
- `.bin` is generated by `make tools` as symlinks to mise-managed binaries
  ([Makefile](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/Makefile#L252-L260)).
  Recreate it; do not copy it.
- `make dev` already runs `bun i`
  ([Makefile](https://github.com/ForkbombEu/credimi/blob/ecfbdb49cdbe09fca2a54ec1387934bc9e89acea/Makefile#L110-L115)).
  Bun's global cache and macOS copy-on-write installation avoid re-downloads and much duplication
  ([Bun docs](https://bun.sh/docs/install/cache)), so copying `node_modules` is an optimization to
  measure, not a prerequisite.

### Recommended tool-neutral contract

1. Add an idempotent repository-owned `scripts/worktree-env` command that:
   - derives labeled starting ports for API, UI, Temporal, and Temporal UI;
   - checks every candidate, increments on collision, and persists final selections, following the
     checked allocation pattern in workmux
     ([example](https://workmux.raine.dev/guide/monorepos#port-isolation));
   - writes a gitignored per-worktree env file containing selected ports, `ADDRESS_UI`,
     `TEMPORAL_ADDRESS`, and sanitized `COMPOSE_PROJECT_NAME`;
   - never copies secrets or data.
2. Parameterize `docker-compose.yaml`, `Procfile.dev`, and the Makefile to consume those values
   while preserving current ports as defaults. Compose supports `COMPOSE_PROJECT_NAME`
   ([project-name precedence](https://docs.docker.com/compose/how-tos/project-name/#set-a-project-name))
   and environment-driven published ports
   ([ports reference](https://docs.docker.com/reference/compose-file/services/#ports)).
3. Use a narrow `.worktreeinclude`:

   ```gitignore
   /.env
   /webapp/.env
   /webapp/node_modules/
   ```

   The `node_modules` entry is optional and should be enabled only after measuring it. Invoke
   `wt step copy-ignored --require-include` so `pb_data`, `.bin`, coverage, and unrelated output
   are excluded
   ([Worktrunk semantics](https://worktrunk.dev/step/#what-gets-copied)).
4. Make `pb_data` fresh by default. If real local data must be reused, provide a separate explicit
   snapshot command that requires the source to be stopped or uses SQLite's supported backup
   mechanism ([SQLite Backup API](https://www.sqlite.org/backup.html)).
5. Recreate `.bin` with `make tools`. Keep Go's module/build caches and Bun's package cache
   user-global. Seed `webapp/node_modules` by APFS reflink only if measurements show `bun i` remains
   a meaningful bottleneck
   ([Worktrunk copy-on-write](https://worktrunk.dev/step/#copy-on-write),
   [Bun cache](https://bun.sh/docs/install/cache)).
6. Scope teardown to the worktree's Compose project. Worktrunk exposes pre/post remove hooks
   ([hooks](https://worktrunk.dev/hook/#hook-types)); project-scoped
   `docker compose -p <name> down` is safer than broad Docker cleanup because Compose uses the
   project name for environment isolation
   ([Docker docs](https://docs.docker.com/compose/how-tos/project-name/)).

### Suggested rollout

1. Parameterize ports and add tool-neutral setup/teardown scripts.
2. Pilot Worktrunk in two worktrees with `--require-include`; leave commit/merge automation off.
3. Compare the repository-owned allocator with pinned Workz v0.11.0. Configure
   `node_modules = "copy"` and verify how one allocated range maps to Credimi's four services
   ([v0.11 isolation](https://github.com/rohansx/workz/blob/v0.11.0/README.md#environment-isolation)).
4. Measure fresh `bun i` against an APFS-reflinked `webapp/node_modules`.
5. Pilot workmux only with multiplexer-oriented maintainers; keep Falq as a design reference until
   adoption and compatibility improve.
6. Evaluate Coasts only if Compose volume/network/port switching remains materially painful after
   the lighter contract is in place.

