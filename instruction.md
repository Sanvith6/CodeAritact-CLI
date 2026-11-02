# RS CLI Commands

Here is a list of all available commands for the `rs.exe` CLI and their corresponding uses:

## `.\rs.exe auth`
*(Previously named `bootstrap`)*
Configure your machine for development. This acts as a guided setup wizard that automatically runs through several important configuration steps:
- Verifies network connectivity.
- Runs system requirements diagnostics (like `rs doctor`).
- Authenticates your AWS SSO profile.
- Configures `uv` package manager and connects it to the private CodeArtifact repository.
- Verifies access to internal validation packages.

## `.\rs.exe doctor`
Check your environment for potential issues and prerequisites. It verifies if you have the required versions of tools such as AWS CLI, Git, uv, Python, and Docker installed.
- **Flag `--fix`:** Run `.\rs.exe doctor --fix` to automatically attempt to resolve configuration issues.

## `.\rs.exe login`
Manually authenticate with AWS SSO. It will open a browser window for you to log in, and then it configures your package manager (uv) with the new credentials.
- **Flag `--force` (or `-f`):** Force a new login flow even if currently authenticated.

## `.\rs.exe logout`
Log out of your current AWS SSO session.

## `.\rs.exe package`
Displays package management information. Note that the RS CLI prepares your environment but does not manage packages directly. This command simply prints instructions and examples on how to use `uv add` to install packages from your configured repositories.

## `.\rs.exe version`
Print the version number, build details, and organization information for the RS Platform CLI.
