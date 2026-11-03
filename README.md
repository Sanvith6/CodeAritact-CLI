# CodeArtifact CLI (`CA CLI`)

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Package Manager](https://img.shields.io/badge/uv-Astral-blueviolet?style=flat)](https://github.com/astral-sh/uv)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> **The Turnkey Developer Onboarding & Private Python Package Distribution Tool.**  
> Connect developers to enterprise AWS CodeArtifact repositories in seconds—with **zero** manual AWS configuration.

---

## 🚀 The Product at a Glance

Enterprises using private Python packages on **AWS CodeArtifact** face a constant onboarding struggle: distributing complex AWS Account IDs, SSO URLs, and managing 12-hour expiring tokens.

**CA CLI solves this in two steps:**
1. **Platform teams** configure [`configs/default.yaml`](configs/default.yaml) and build a single, branded executable (`ca.exe`).
2. **Developers** download `ca.exe` and run **`ca auth`**. No AWS knowledge required.

```
┌─────────────────────────────────┐       ┌─────────────────────────────────┐
│     DevOps / Platform Team      │       │            Developer            │
│                                 │       │                                 │
│  1. Edit configs/default.yaml   │       │  1. Put ca.exe in project root  │
│  2. Build: go build -o ca.exe   │ ────> │  2. Run: .\ca.exe auth          │
│                                 │       │                                 │
│  (Baked-in enterprise config)   │       │  (Dependencies synced in sec)   │
└─────────────────────────────────┘       └─────────────────────────────────┘
```

---

## ⚡ Minimum Setup Steps (Product Lifecycle)

### Phase 1: DevOps / Maintainer (Build the Product)

Configure your organization's cloud parameters and generate the self-contained binary:

#### Step 1: Edit [`configs/default.yaml`](configs/default.yaml)
Set your organization's AWS SSO and CodeArtifact repository details:

```yaml
version: 1

aws:
  sso:
    start_url: "https://your-company.awsapps.com/start"  # Your SSO Portal URL
    region: "us-east-1"                                  # AWS Region
    account_id: "080800845757"                           # Target AWS Account ID
    role_name: "Developer-permission"                    # IAM Role / Permission Set
  profile_name: "Developer"                              # Local AWS Profile Name
  region: "us-east-1"

artifact:
  domain: "grag-ai-factory"                              # CodeArtifact Domain
  domain_owner: "080800845757"                           # Domain Owner Account ID
  repository: "python-packages"                          # CodeArtifact Repository
  region: "us-east-1"

package:
  uv_index_name: "internal"                              # Local uv index name
  validation_package: "company-web"                      # Internal canary package to verify access
```

#### Step 2: Build the Standalone Binary
Run the code generator to bake the YAML config into Go, then compile:

```powershell
# 1. Embed the YAML into internal/config/embedded.go
go generate ./...

# 2. Build the optimized single executable
go build -ldflags="-s -w" -o ca.exe main.go

# 3. (Optional) Zip for distribution
Compress-Archive -Path ca.exe -DestinationPath ca.zip -Force
```

Distribute `ca.exe` (or `ca.zip`) to your developers via internal release, S3, or GitHub.

---

### Phase 2: Developer (Use the Product)

Developers need **zero AWS knowledge**—no Account IDs, no IAM roles, and no endpoints.

#### Step 1: Download `ca.exe`
Download **[`ca.zip`](https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip)** from this repo, or run in PowerShell:

```powershell
Invoke-WebRequest -Uri "https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip" -OutFile "ca.zip"; Expand-Archive ca.zip -DestinationPath . -Force; Remove-Item ca.zip
```

Place `ca.exe` in the root of your Python project alongside your `requirements.txt` or `pyproject.toml`.

#### Step 2: Run Setup
```powershell
.\ca.exe auth
```

The CLI executes the full onboarding workflow:
1. **Network Check**: Confirms internet & AWS endpoint access.
2. **Tool Check**: Verifies `aws`, `git`, `uv`, and `python`.
3. **AWS SSO Login**: Automatically launches the browser for login.
4. **Token Generation**: Fetches a 12-hour CodeArtifact authentication token.
5. **uv Configuration**: Injects credentials into `uv.toml` global index.
6. **Canary Verification**: Confirms private package repository access.
7. **Sync Dependencies**: Installs project packages directly from CodeArtifact.

---

## 📦 How Dependency Pulling Works

`ca auth` automatically inspects your project directory to determine how dependencies should be installed:

| If your project has | What `ca auth` executes |
| :--- | :--- |
| **`requirements.txt`** | Checks for `.venv` (creates one via `uv venv` if missing), then runs:<br>`uv pip install -r requirements.txt` |
| **`pyproject.toml`** | Resolves the project lockfile and installs via:<br>`uv sync` |

### Adding or Updating Packages
Simply add the package name to `requirements.txt` or `pyproject.toml` and re-run:
```powershell
.\ca.exe auth
```
Or install directly with `uv`:
```powershell
uv add <package-name>
# or
uv pip install <package-name>
```

---

## 🛠️ CLI Commands

| Command | Action |
| :--- | :--- |
| `.\ca.exe auth` | **Main command.** Logs in via SSO, generates tokens, configures `uv`, and syncs dependencies. |
| `.\ca.exe doctor` | Verifies your local environment and tool prerequisites. |
| `.\ca.exe doctor --fix` | Automatically fixes configuration issues and refreshes expired SSO sessions. |
| `.\ca.exe login` | Re-authenticates with AWS SSO and refreshes the CodeArtifact `uv` token. |
| `.\ca.exe logout` | Clears active AWS SSO credentials. |
| `.\ca.exe version` | Displays CLI version, architecture, and build metadata. |

---

## 🔒 Security & Governance

- **Zero Hardcoded Secrets**: `ca.exe` contains **no** AWS private keys, passwords, or permanent credentials. It only embeds non-secret routing metadata (SSO URL, Account ID, Domain).
- **Identity-First Security**: Actual repository access strictly requires an active corporate SSO identity (Okta, Azure AD, Ping) assigned to the designated IAM Role.
- **Ephemeral Credentials**: CodeArtifact tokens generated at runtime are temporary (12 hours maximum) and tied directly to the developer's session.
- **Reverse Engineering Note**: As with any compiled binary, embedded strings can be extracted using tools like `strings` or decompilers. Access control is enforced by cloud IAM and your Identity Provider, not binary obfuscation.

---

## 📄 License

[MIT License](LICENSE)
