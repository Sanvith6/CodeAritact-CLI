# RS CLI: Developer & DevOps Guide

This document outlines the commands available to developers using the RS CLI, as well as the build process for DevOps engineers distributing the tool.

---

## 👨‍💻 Developer Commands

Here is a list of all available commands for the `rs.exe` CLI and their corresponding uses:

### `.\rs.exe auth`
*(Previously named `bootstrap`)*
Configure your machine for development. This acts as a guided setup wizard that automatically runs through several important configuration steps:
- Verifies network connectivity.
- Runs system requirements diagnostics (like `rs doctor`).
- Authenticates your AWS SSO profile.
- Configures `uv` package manager and connects it to the private CodeArtifact repository.
- Verifies access to internal validation packages.

### `.\rs.exe doctor`
Check your environment for potential issues and prerequisites. It verifies if you have the required versions of tools such as AWS CLI, Git, uv, Python, and Docker installed.
- **Flag `--fix`:** Run `.\rs.exe doctor --fix` to automatically attempt to resolve configuration issues.

### `.\rs.exe login`
Manually authenticate with AWS SSO. It will open a browser window for you to log in, and then it configures your package manager (uv) with the new credentials.
- **Flag `--force` (or `-f`):** Force a new login flow even if currently authenticated.

### `.\rs.exe logout`
Log out of your current AWS SSO session.

### `.\rs.exe package`
Displays package management information. Note that the RS CLI prepares your environment but does not manage packages directly. This command simply prints instructions and examples on how to use `uv add` to install packages from your configured repositories.

### `.\rs.exe version`
Print the version number, build details, and organization information for the RS Platform CLI.

---

## 🛠️ DevOps / Platform Engineer Workflow

As a DevOps engineer, you are responsible for embedding the organization's configuration into the CLI before distributing it to developers.

### 1. Configuration
First, edit the `configs/default.yaml` file to match your organization's environment. This includes setting:
- AWS SSO URLs, Region, Account IDs, and Roles.
- CodeArtifact Domain, Domain Owner, and Repository names.
- Specific feature flags (e.g., enabling/disabling Docker, Java, Terraform).

### 2. Building the CLI
After making changes to the configuration, you must re-generate the embedded Go code and compile the new binary. Run the following commands from the repository root:

```powershell
# 1. Generate the embedded configuration code (crucial after editing default.yaml)
go generate ./...

# 2. Build the final executable
go build -o rs.exe main.go
```

*(Alternatively, if you have `make` installed on Windows via WSL, Git Bash, or MSYS2, you can run `make build` — but you still must run `go generate ./...` first.)*

### 3. Distribution
Once `rs.exe` is successfully built, you can distribute this single binary to your developers. 

Developers do not need to know the AWS Account IDs, SSO endpoints, or CodeArtifact details. They simply run `.\rs.exe auth`, and the embedded configuration handles the rest securely!
