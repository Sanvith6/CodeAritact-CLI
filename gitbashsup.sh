#!/usr/bin/env bash

# ==============================================================================
# AWS CodeArtifact Developer Setup (Windows / Git Bash) - DIAGNOSTIC MODE
# ==============================================================================
# This script configures the AWS CLI to authenticate via IAM Identity Center (SSO)
# and retrieves a CodeArtifact token to securely configure pip for the user.
# It is designed to be idempotent and safe for repeated execution.
# ==============================================================================

set -Eeuo pipefail
IFS=$'\n\t'

# ------------------------------------------------------------------------------
# Configuration Variables
# ------------------------------------------------------------------------------
AWS_REGION="${AWS_REGION:-us-east-1}"
SSO_REGION="${SSO_REGION:-us-east-1}"
AWS_ACCOUNT_ID="${AWS_ACCOUNT_ID:-123456789012}"
SSO_ROLE_NAME="${SSO_ROLE_NAME:-DeveloperAccess}"
PROFILE_NAME="${PROFILE_NAME:-Developer}"
CA_DOMAIN="${CA_DOMAIN:-my-company-domain}"
CA_DOMAIN_OWNER="${CA_DOMAIN_OWNER:-123456789012}"
CA_REPO="${CA_REPO:-python-packages}"
SSO_START_URL="${SSO_START_URL:-https://my-company.awsapps.com/start}"

DIAGNOSTIC_MODE="${DIAGNOSTIC_MODE:-true}"

# ------------------------------------------------------------------------------
# ANSI Colors for Output
# ------------------------------------------------------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Global State Variables
TMPFILE=$(mktemp)
CURRENT_STEP_NAME=""
CURRENT_COMMAND=""
declare -A STEP_RESULTS
declare -a STEP_LIST
SCRIPT_START_TIME=$(date +%s)

# ------------------------------------------------------------------------------
# Cleanup Function
# ------------------------------------------------------------------------------
cleanup() {
    rm -f "$TMPFILE"
}

trap 'log_error "Error on line $LINENO"; cleanup' ERR
trap cleanup EXIT

# ------------------------------------------------------------------------------
# Logging Helpers
# ------------------------------------------------------------------------------
get_timestamp() {
    date "+%Y-%m-%d %H:%M:%S"
}

log_message() {
    local level=$1
    local color=$2
    local message=$3
    echo -e "[$(get_timestamp)]\n\n${color}${level}${NC}\n\n${message}\n"
}

log_info()    { log_message "INFO" "$CYAN" "$1"; }
log_success() { log_message "SUCCESS" "$GREEN" "$1"; }
log_warn()    { log_message "WARNING" "$YELLOW" "$1"; }
log_error()   { log_message "ERROR" "$RED" "$1"; }
log_debug()   { log_message "DEBUG" "$BLUE" "$1"; }

# ------------------------------------------------------------------------------
# Diagnostic Helpers
# ------------------------------------------------------------------------------
pause_if_diagnostic() {
    if [[ "${DIAGNOSTIC_MODE}" == "true" ]]; then
        echo -e "${YELLOW}Press ENTER to continue...${NC}"
        read -r
    fi
}

start_step() {
    local step_num=$1
    local step_desc=$2
    CURRENT_STEP_NAME="$step_desc"
    STEP_LIST+=("$step_desc")
    
    echo -e "\n========================================================="
    echo -e "STEP $step_num"
    echo -e "$step_desc"
    echo -e "========================================================="
    pause_if_diagnostic
}

end_step() {
    local status=$1
    if [[ "$status" == "0" ]]; then
        echo -e "\n${GREEN}✓ Completed${NC}\n"
        STEP_RESULTS["$CURRENT_STEP_NAME"]="${GREEN}✓${NC} $CURRENT_STEP_NAME"
    else
        echo -e "\n${RED}✗ Failed${NC}\n"
        STEP_RESULTS["$CURRENT_STEP_NAME"]="${RED}✗${NC} $CURRENT_STEP_NAME"
    fi
}

# Wrapper to execute commands, capture exit code, timer and print
run_command() {
    local cmd_string="$1"
    local allow_failure="${2:-false}"
    local resolution="${3:-Check the command output and configuration.}"
    
    CURRENT_COMMAND="$cmd_string"

    echo -e "[COMMAND]\n$cmd_string\n"
    
    echo -e "-----------------------------------------"
    echo -e "Executing:\n\n$cmd_string\n"
    
    local start_cmd
    start_cmd=$(date +%s.%N)
    
    local stderr_tmp
    stderr_tmp=$(mktemp)
    local stdout_tmp
    stdout_tmp=$(mktemp)
    
    set +e
    eval "$cmd_string" > >(tee "$stdout_tmp") 2> >(tee "$stderr_tmp" >&2)
    local exit_code=$?
    set -e
    
    local end_cmd
    end_cmd=$(date +%s.%N)
    
    local duration
    duration=$(echo "$end_cmd - $start_cmd" | awk '{printf "%.2f", $1}')
    
    echo -e "\nExit Code : $exit_code"
    echo -e "Duration  : $duration sec"
    
    # Store output for calling functions to evaluate if needed
    CMD_OUTPUT=$(cat "$stdout_tmp")
    CMD_ERROR=$(cat "$stderr_tmp")
    rm -f "$stdout_tmp" "$stderr_tmp"
    
    if [[ $exit_code -eq 0 ]]; then
        echo -e "Status    : ${GREEN}SUCCESS${NC}"
        echo -e "-----------------------------------------\n"
        return 0
    else
        echo -e "Status    : ${RED}FAILED${NC}"
        echo -e "-----------------------------------------\n"
        
        if [[ "$allow_failure" != "true" ]]; then
            echo -e "${RED}Current Step${NC}: $CURRENT_STEP_NAME"
            echo -e "${RED}Command${NC}: $cmd_string"
            echo -e "${RED}Exit Code${NC}: $exit_code"
            echo -e "${RED}stderr${NC}:"
            echo -e "$CMD_ERROR"
            echo -e "\n${RED}Suggested Resolution${NC}: $resolution\n"
            
            end_step 1
            print_summary
            exit 1
        fi
        return $exit_code
    fi
}

# ------------------------------------------------------------------------------
# 1. Banner
# ------------------------------------------------------------------------------
print_banner() {
    echo -e "${BLUE}==========================================${NC}"
    echo -e "${GREEN} AWS CodeArtifact Developer Setup (Diagnostic)${NC}"
    echo -e "${BLUE}==========================================${NC}"
    echo ""
}

# ------------------------------------------------------------------------------
# 2. Check Operating System
# ------------------------------------------------------------------------------
check_os() {
    start_step "1" "Operating System"
    log_info "Checking operating system..."
    
    local os_name
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [[ -f /etc/os-release ]]; then
            # shellcheck disable=SC1091
            . /etc/os-release
            os_name=${ID:-unknown}
            if [[ ! "$os_name" =~ ^(ubuntu|debian|amzn|rhel|rocky|centos)$ ]]; then
                log_error "Unsupported Linux distribution: $os_name. Supported: ubuntu, debian, amzn, rhel, rocky."
                end_step 1
                print_summary
                exit 1
            fi
        else
            log_error "Cannot determine Linux distribution. /etc/os-release not found."
            end_step 1
            print_summary
            exit 1
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        os_name="macos"
    elif [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "cygwin"* ]]; then
        os_name="windows"
    else
        log_error "Unsupported OS type: $OSTYPE. Supported: Linux, macOS, msys/cygwin (Windows)."
        end_step 1
        print_summary
        exit 1
    fi
    
    log_success "OS validated ($os_name)."
    end_step 0
}

# ------------------------------------------------------------------------------
# 3. Check Internet Connectivity
# ------------------------------------------------------------------------------
check_internet() {
    start_step "2" "Internet"
    log_info "Checking internet connectivity..."
    
    run_command "curl -sI --connect-timeout 5 \"https://sts.${AWS_REGION}.amazonaws.com\" > /dev/null" "false" "Please check your internet connection."
    
    log_success "Internet connectivity verified."
    end_step 0
}

# ------------------------------------------------------------------------------
# 4. Check Required Tools
# ------------------------------------------------------------------------------
check_tools() {
    start_step "3" "Tools"
    log_info "Verifying required tools..."
    
    local required_cmds=(aws grep sed jq curl)
    local missing_cmds=()

    for cmd in "${required_cmds[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_cmds+=("$cmd")
        fi
    done

    if ! command -v python >/dev/null 2>&1 && ! command -v python3 >/dev/null 2>&1; then
        missing_cmds+=("python")
    fi
    if ! command -v pip >/dev/null 2>&1 && ! command -v pip3 >/dev/null 2>&1; then
        missing_cmds+=("pip")
    fi

    if [[ ${#missing_cmds[@]} -gt 0 ]]; then
        log_error "The following required tools are missing: ${missing_cmds[*]}"
        echo -e "\nPlease install them."
        echo -e "Examples on Windows (using PowerShell/Winget):"
        echo -e "  winget install Amazon.AWSCLI"
        echo -e "  winget install Python.Python.3"
        echo -e "  winget install jqlang.jq"
        end_step 1
        print_summary
        exit 1
    fi
    
    log_success "All required tools are installed."
    end_step 0
}

# ------------------------------------------------------------------------------
# 5. Verify AWS CLI Version
# ------------------------------------------------------------------------------
verify_aws_cli_version() {
    start_step "4" "AWS CLI"
    log_info "Verifying AWS CLI version..."
    
    set +e
    run_command "aws --version 2>&1" "true"
    set -e
    
    local aws_version
    aws_version=$(echo "$CMD_OUTPUT" | cut -d ' ' -f 1 || true)
    
    if [[ "$aws_version" != aws-cli/2.* ]]; then
        log_error "AWS CLI version 2 is required. Found: $aws_version. Please upgrade."
        end_step 1
        print_summary
        exit 1
    fi
    log_success "AWS CLI v2 detected."
    end_step 0
}

# ------------------------------------------------------------------------------
# 6. Configure AWS SSO Profile
# ------------------------------------------------------------------------------
configure_sso_profile() {
    start_step "5" "AWS Profile"
    log_info "Configuring AWS SSO profile '${PROFILE_NAME}'..."
    
    # Backup existing AWS profile check
    set +e
    run_command "aws configure list-profiles 2>/dev/null" "true"
    if echo "$CMD_OUTPUT" | grep -q "^${PROFILE_NAME}$"; then
        log_warn "Existing ${PROFILE_NAME} profile found. Updating profile..."
    fi
    set -e
    
    run_command "aws configure set sso_start_url \"$SSO_START_URL\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    run_command "aws configure set sso_region \"$SSO_REGION\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    run_command "aws configure set sso_account_id \"$AWS_ACCOUNT_ID\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    run_command "aws configure set sso_role_name \"$SSO_ROLE_NAME\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    run_command "aws configure set region \"$AWS_REGION\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    run_command "aws configure set output \"json\" --profile \"$PROFILE_NAME\"" "false" "Ensure AWS CLI is configured."
    echo -e "✓ Profile updated\n"
    
    log_success "AWS profile '${PROFILE_NAME}' configured successfully."
    end_step 0
}

# ------------------------------------------------------------------------------
# 7 & 8. Detect Existing Session and Perform AWS SSO Login
# ------------------------------------------------------------------------------
perform_sso_login() {
    start_step "6" "AWS Login"
    log_info "Checking for existing SSO session..."
    
    local session_valid=false
    
    if [[ -d ~/.aws/sso/cache ]]; then
        set +e
        run_command "aws sts get-caller-identity --profile \"$PROFILE_NAME\"" "true"
        if [[ $? -eq 0 ]]; then
            session_valid=true
        fi
        set -e
    fi
    
    if [[ "$session_valid" == true ]]; then
        log_success "Existing SSO session found and is valid. Skipping login."
    else
        log_warn "No active session found or session expired."
        log_info "Initiating AWS SSO login via device code..."
        
        echo -e "\n${YELLOW}=== Action Required ===${NC}"
        echo -e "1. Open the URL shown below in your browser."
        echo -e "   ${CYAN}(Tip: Use an Incognito/Private window if you need to switch users or force a password/MFA prompt)${NC}"
        echo -e "2. Enter the device code provided."
        echo -e "3. Complete any MFA required."
        echo -e "4. Return to this terminal once authorized.\n"
        
        echo -e "========================================================="
        echo -e "ABOUT TO START AWS IAM IDENTITY CENTER LOGIN"
        echo -e ""
        echo -e "If Windows asks for Administrator credentials,"
        echo -e "note exactly when it appears."
        echo -e "=========================================================\n"

        local max_retries=3
        local attempt=1
        local login_success=false
        
        while [[ $attempt -le $max_retries ]]; do
            set +e
            run_command "aws sso login --profile \"$PROFILE_NAME\" --use-device-code" "true"
            if [[ $? -eq 0 ]]; then
                log_success "SSO login completed."
                login_success=true
                break
            else
                log_warn "Login attempt $attempt failed."
                attempt=$((attempt + 1))
                if [[ $attempt -le $max_retries ]]; then
                    log_info "Retrying login... (Attempt $attempt of $max_retries)"
                fi
            fi
            set -e
        done
        
        if [[ "$login_success" == false ]]; then
            log_error "Failed to authenticate via AWS SSO after $max_retries attempts."
            end_step 1
            print_summary
            exit 1
        fi
    fi
    end_step 0
}

# ------------------------------------------------------------------------------
# 9 & 10. Verify Identity and Print Expiration
# ------------------------------------------------------------------------------
verify_identity() {
    start_step "7" "Identity"
    log_info "Verifying AWS identity context..."
    
    set +e
    run_command "aws sts get-caller-identity --profile \"$PROFILE_NAME\"" "true"
    local identity_exit_code=$?
    local caller_identity="$CMD_OUTPUT"
    set -e
    
    if [[ $identity_exit_code -ne 0 ]]; then
        echo -e "\n${RED}=== Role Access Denied ===${NC}"
        echo -e "Your SSO login succeeded, but you do not have permission to assume the configured role."
        echo -e "  Target Account ID: ${AWS_ACCOUNT_ID}"
        echo -e "  Target Role Name:  ${SSO_ROLE_NAME}"
        
        # Automatically discover available roles
        local cache_dir=~/.aws/sso/cache
        if [[ -d "$cache_dir" ]]; then
            local token_file
            token_file=$(grep -l accessToken "$cache_dir"/*.json 2>/dev/null | xargs ls -t 2>/dev/null | head -n 1 || true)
            if [[ -n "$token_file" ]]; then
                local access_token
                access_token=$(jq -r .accessToken "$token_file")
                if [[ -n "$access_token" && "$access_token" != "null" ]]; then
                    echo -e "\n${CYAN}Discovering your available roles in account ${AWS_ACCOUNT_ID}...${NC}"
                    
                    set +e
                    run_command "aws sso list-account-roles --account-id \"$AWS_ACCOUNT_ID\" --access-token \"$access_token\"" "true"
                    local roles_json="$CMD_OUTPUT"
                    set -e
                    
                    if [[ -n "$roles_json" ]]; then
                        local available_roles
                        available_roles=$(echo "$roles_json" | jq -r '.roleList[].roleName')
                        if [[ -n "$available_roles" ]]; then
                            echo -e "${GREEN}Good news! We found the following roles assigned to you:${NC}"
                            while IFS= read -r role; do
                                echo -e "  - ${YELLOW}$role${NC}"
                            done <<< "$available_roles"
                            echo -e "\nTo use one of these roles, simply run:"
                            echo -e "  ${YELLOW}export SSO_ROLE_NAME=\"<RoleName>\"${NC}"
                            echo -e "  ${YELLOW}./gitbashsup.sh${NC}"
                            echo -e "\n${CYAN}------------------------------------------${NC}"
                        fi
                    fi
                fi
            fi
        fi

        echo -e "\n${YELLOW}Other Troubleshooting Steps:${NC}"
        echo -e "1. If you accidentally logged in as the wrong user, log out locally first:"
        echo -e "     aws sso logout"
        echo -e "   Then run ./gitbashsup.sh again and open the login link in an ${CYAN}Incognito/Private${NC} browser window."
        echo -e "2. Verify with your DevOps team that you have been granted access to this account.\n"
        
        log_error "Could not retrieve caller identity."
        end_step 1
        print_summary
        exit 1
    fi
    
    local account_id
    local arn
    local user_id
    
    account_id=$(echo "$caller_identity" | jq -r .Account)
    arn=$(echo "$caller_identity" | jq -r .Arn)
    user_id=$(echo "$caller_identity" | jq -r .UserId)
    
    # Try to find the session expiration from SSO cache
    local expires_at=""
    local cache_dir=~/.aws/sso/cache
    if [[ -d "$cache_dir" ]]; then
        local latest_cache
        # Get the most recently modified json file in the cache directory
        latest_cache=$(ls -t "$cache_dir"/*.json 2>/dev/null | head -n 1 || true)
        if [[ -n "$latest_cache" && -f "$latest_cache" ]]; then
            expires_at=$(jq -r '.expiresAt // empty' "$latest_cache")
        fi
    fi
    
    echo -e "\n========================================================="
    echo -e "IDENTITY"
    echo -e "========================================================="
    echo -e "Account            : ${account_id}"
    echo -e "ARN                : ${arn}"
    echo -e "UserID             : ${user_id}"
    if [[ -n "$expires_at" ]]; then
        echo -e "Session Expiration : ${expires_at} UTC"
    else
        echo -e "Session Expiration : Unknown"
    fi
    echo -e "=========================================================\n"
    
    end_step 0
}

# ------------------------------------------------------------------------------
# 11. Configure pip (CodeArtifact)
# ------------------------------------------------------------------------------
configure_pip() {
    start_step "8" "CodeArtifact"
    log_info "Configuring pip for CodeArtifact..."
    
    echo -e "Starting aws codeartifact login at: $(get_timestamp)\n"
    
    run_command "aws codeartifact login --tool pip --repository \"$CA_REPO\" --domain \"$CA_DOMAIN\" --domain-owner \"$CA_DOMAIN_OWNER\" --profile \"$PROFILE_NAME\"" "false" "Ensure you have the right permissions to access CodeArtifact."
    
    echo -e "\nFinished aws codeartifact login at: $(get_timestamp)\n"
    
    log_success "pip configured successfully via 'aws codeartifact login'."
    end_step 0
}

# ------------------------------------------------------------------------------
# 12. Verify Repository Access
# ------------------------------------------------------------------------------
verify_repo_access() {
    start_step "9" "Repository"
    log_info "Verifying repository access via AWS CLI..."
    
    run_command "aws codeartifact list-packages --domain \"$CA_DOMAIN\" --domain-owner \"$CA_DOMAIN_OWNER\" --repository \"$CA_REPO\" --max-results 1 --profile \"$PROFILE_NAME\"" "false" "Ensure you have permissions to list packages."
    
    log_success "Repository Connected Successfully."
    end_step 0
}

# ------------------------------------------------------------------------------
# 13. Test Package Installation (Dry Run)
# ------------------------------------------------------------------------------
test_package_install() {
    start_step "10" "Package Validation"
    log_info "Testing package access (dry verification)..."
    
    local test_pkg="company-web"
    local pip_cmd="pip"
    
    if command -v python >/dev/null 2>&1; then
        pip_cmd="python -m pip"
    elif command -v python3 >/dev/null 2>&1; then
        pip_cmd="python3 -m pip"
    fi

    set +e
    run_command "$pip_cmd index versions \"$test_pkg\"" "true"
    local test_pkg_exit=$?
    set -e
    
    if [[ $test_pkg_exit -eq 0 ]]; then
         log_success "✔ Repository connection verified (found ${test_pkg})."
    else
         log_warn "Could not verify '${test_pkg}'. Ensure the package exists in the repository."
    fi
    end_step 0
}

# ------------------------------------------------------------------------------
# Final Summary
# ------------------------------------------------------------------------------
print_summary() {
    local end_time=$(date +%s)
    local total_time=$((end_time - SCRIPT_START_TIME))
    
    echo -e "\n========================================================="
    echo -e "Execution Summary"
    echo -e "========================================================="
    echo ""
    
    for step in "${STEP_LIST[@]}"; do
        if [[ -n "${STEP_RESULTS[$step]+set}" ]]; then
            echo -e "${STEP_RESULTS[$step]}"
        else
            echo -e "  Not executed: $step"
        fi
    done
    
    echo -e "\nElapsed Time : $total_time seconds"
    echo -e "\n=========================================================\n"
}

# ------------------------------------------------------------------------------
# Main Execution Flow
# ------------------------------------------------------------------------------
main() {
    print_banner
    check_os
    check_internet
    check_tools
    verify_aws_cli_version
    configure_sso_profile
    perform_sso_login
    verify_identity
    configure_pip
    verify_repo_access
    test_package_install
    
    print_summary
    
    echo -e "${BLUE}==========================================${NC}"
    echo -e "${GREEN} Setup Complete${NC}"
    echo -e "${BLUE}==========================================${NC}"
    echo -e "\n${CYAN}Tip: If you are on Windows, create a virtual environment first:${NC}"
    echo -e "  python -m venv .venv && source .venv/Scripts/activate\n"
    
    echo -e "You can now install packages using:"
    echo -e "  ${YELLOW}pip install <package-name>${NC}"
    echo -e "\nExample:"
    echo -e "  ${YELLOW}pip install company-web==1.0.5${NC}\n"
    echo -e "${CYAN}Note: CodeArtifact tokens expire every 12 hours. If you receive authentication errors, simply rerun this setup script to refresh the token while reusing your active SSO session.${NC}\n"
}

# Execute main function
main "$@"
