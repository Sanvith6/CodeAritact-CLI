# CA CLI (CodeArtifact CLI)

> **Zero-configuration AWS CodeArtifact & SSO authentication for Python teams.**  
> One binary. One command. All dependencies installed.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![uv](https://img.shields.io/badge/uv-Package%20Manager-blueviolet)](https://github.com/astral-sh/uv)

---

## 🎯 Why This is Important

Setting up internal Python packages on AWS CodeArtifact is one of the biggest friction points in engineering teams:

| Traditional Setup (Without CA CLI) | With CA CLI (This Product) |
| :--- | :--- |
| ❌ **10-step manual runbook** for every new engineer | ✅ **One command**: `.\ca.exe auth` |
| ❌ Developers must ask DevOps for **Account IDs, SSO URLs & Regions** | ✅ **Zero AWS knowledge**: Endpoints are baked into the tool |
| ❌ **Tokens expire every 12 hours**, breaking daily builds | ✅ **Automatic token generation & renewal** |
| ❌ Fragile 600-line shell scripts that fail on Windows/Mac | ✅ **Single compiled executable** (`ca.exe`) |
| ❌ Complex manual configuration for package managers | ✅ **Automated `uv` integration** with private repo index |

---

## ⚡ For Developers (2 Steps)

### Step 1: Download `ca.exe`
Run in PowerShell inside your project folder:
```powershell
Invoke-WebRequest -Uri "https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip" -OutFile ca.zip; Expand-Archive ca.zip -DestinationPath . -Force; Remove-Item ca.zip
```
*(Or manually download and extract [**`ca.zip`**](https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip))*

### Step 2: Run Setup
```powershell
.\ca.exe auth
```

**That is all.** The CLI automatically:
1. Opens your browser to authenticate with AWS SSO.
2. Retrieves a temporary CodeArtifact token.
3. Configures Astral's `uv` package manager.
4. Creates a virtual environment (`.venv`) and installs all project dependencies.

---

## 🛠️ For DevOps & Maintainers (2 Steps)

### Step 1: Configure `configs/default.yaml`
Add your organization's cloud settings to [`configs/default.yaml`](configs/default.yaml):

```yaml
aws:
  sso:
    start_url: "https://my-company.awsapps.com/start"  # Your SSO Portal URL
    region: "us-east-1"
    account_id: "123456789012"                         # Target AWS Account ID
    role_name: "DeveloperAccess"                       # IAM Role
artifact:
  domain: "my-company-domain"                          # CodeArtifact Domain
  domain_owner: "123456789012"                         # Account ID
  repository: "python-packages"                        # CodeArtifact Repository
package:
  uv_index_name: "internal"
  validation_package: "my-package"                     # Internal canary package
```

### Step 2: Build & Distribute
```powershell
go generate ./...
go build -ldflags="-s -w" -o ca.exe main.go
Compress-Archive -Path ca.exe -DestinationPath ca.zip -Force
```
Distribute `ca.exe` (or `ca.zip`) to your developers.

---

## 📦 How Dependency Pulling Works

When you run `.\ca.exe auth`, it automatically checks your project root:

- **`requirements.txt` present?**  
  Checks for `.venv` (creates one if missing), then runs:  
  `uv pip install -r requirements.txt`

- **`pyproject.toml` present?**  
  Resolves the lockfile and installs packages via:  
  `uv sync`

- **Adding a new package?**  
  Add the package name to `requirements.txt` or `pyproject.toml` and re-run `.\ca.exe auth`, or use `uv add <package>` directly.

---

## 🩺 Helpful Commands

| Command | Description |
| :--- | :--- |
| `.\ca.exe auth` | **Main command.** Logs in via SSO, generates tokens, and installs dependencies. |
| `.\ca.exe doctor` | Verifies required tools (`aws`, `git`, `uv`, `python`). |
| `.\ca.exe doctor --fix` | Automatically fixes configuration and refreshes expired login sessions. |
| `.\ca.exe login` | Re-authenticates with AWS SSO and refreshes the CodeArtifact token. |
| `.\ca.exe logout` | Logs out of active AWS SSO session. |

---

## 🔒 Security

- **No Secrets Inside**: `ca.exe` contains **no** AWS passwords, private keys, or credentials. It only embeds routing identifiers (Account ID, Domain name, SSO URL).
- **Access Control**: Actual repository access strictly requires an active corporate SSO identity. Cloud permissions are enforced by AWS IAM, not client-side binary obfuscation.

---

## 📄 License

[MIT License](LICENSE)
