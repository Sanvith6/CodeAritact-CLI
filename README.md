# CodeArtifact CLI (`CA CLI`)

> **Zero-setup AWS CodeArtifact & SSO authentication for Python developers.**  
> No AWS Account IDs, Regions, or SSO URLs required. Just run `ca auth`.

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Package Manager](https://img.shields.io/badge/uv-Astral-blueviolet?style=flat)](https://github.com/astral-sh/uv)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## ⚡ Quickstart for Developers

### 1. Download `ca.exe`

Download **[`ca.zip`](https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip)** directly from this repo, or run in PowerShell:

```powershell
Invoke-WebRequest -Uri "https://github.com/Sanvith6/CodeAritact-CLI/raw/main/ca.zip" -OutFile "ca.zip"; Expand-Archive ca.zip -DestinationPath . -Force; Remove-Item ca.zip
```

Place `ca.exe` in the root of your Python project.

### 2. Make sure you have your dependencies listed

In your project folder, have either:
- **`requirements.txt`** (e.g. `company-web==1.2.3`), OR
- **`pyproject.toml`**

### 3. Run Setup

In your terminal, run:

```powershell
.\ca.exe auth
```

**That's it!** The CLI will:
1. Open your browser to log into AWS SSO.
2. Fetch a temporary CodeArtifact token.
3. Automatically configure `uv` package manager.
4. Create a virtual environment (`.venv`) and install all your dependencies.

---

## 📦 How Dependency Syncing Works

When you run `.\ca.exe auth`:

- **If you use `requirements.txt`**:  
  It checks for a `.venv` folder (creates one if missing), then runs:  
  `uv pip install -r requirements.txt`

- **If you use `pyproject.toml`**:  
  It runs:  
  `uv sync`

**Need to add or update packages?**  
Add the package name to `requirements.txt` or `pyproject.toml` and re-run:
```powershell
.\ca.exe auth
```
Or use `uv` directly:
```powershell
uv add <package-name>
```

---

## 🛠️ Commands

| Command | What it does |
| :--- | :--- |
| `.\ca.exe auth` | **Main command.** Logs in and installs all dependencies. |
| `.\ca.exe doctor` | Checks if required tools (`aws`, `git`, `uv`, `python`) are installed. |
| `.\ca.exe doctor --fix` | Automatically fixes configuration and refreshes expired login sessions. |
| `.\ca.exe login` | Manually re-authenticates with AWS SSO and refreshes repository token. |
| `.\ca.exe logout` | Logs out of AWS SSO. |
| `.\ca.exe version` | Shows CLI version and build details. |

---

## 👨‍💻 For DevOps / Maintainers: Building `ca.exe`

To build a pre-configured binary for your team so developers don't have to know any AWS details:

### 1. Set your company's AWS settings in `configs/default.yaml`

Edit `configs/default.yaml` with your SSO and CodeArtifact information:

```yaml
aws:
  sso:
    start_url: "https://your-company.awsapps.com/start"
    region: "us-east-1"
    account_id: "080800845757"
    role_name: "Developer-permission"
  profile_name: "Developer"
  region: "us-east-1"

artifact:
  domain: "grag-ai-factory"
  domain_owner: "080800845757"
  repository: "python-packages"
  region: "us-east-1"

package:
  uv_index_name: "internal"
  validation_package: "company-web"
```

### 2. Generate and Build the `.exe`

Run `scripts/generate_config.go` to bake the YAML settings into Go code, then compile:

```powershell
# 1. Embed the YAML config into internal/config/embedded.go
go generate ./...
# (or: go run scripts/generate_config.go)

# 2. Build the ca.exe binary
go build -o ca.exe main.go

# 3. Create zip for distribution
Compress-Archive -Path ca.exe -DestinationPath ca.zip -Force
```

Give `ca.zip` to your developers. They just extract it and run `.\ca.exe auth`.

---

## 🔒 Security Note

- **No Secrets Inside**: `ca.exe` contains **no** AWS passwords, private keys, or credentials. It only embeds routing identifiers (Account ID, Domain name, SSO URL).
- **Access Control**: Actual repository access strictly requires an active, approved corporate SSO identity.
- *Disclaimer*: As with any compiled binary, embedded strings can be viewed with reverse-engineering tools. Security is enforced by AWS IAM and your enterprise Identity Provider, not binary obfuscation.

---

## 📄 License

[MIT License](LICENSE)
