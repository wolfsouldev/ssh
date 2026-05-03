# setup-hooks.ps1
# Sets up Git hooks for the SSHH project (similar to Husky for Node.js)

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "  Setting up Git hooks for SSHH..." -ForegroundColor Cyan
Write-Host ""

# Ensure we're in a git repo
if (-not (Test-Path ".git")) {
    Write-Host "  Error: Not a Git repository. Run 'git init' first." -ForegroundColor Red
    exit 1
}

# Create hooks directory if needed
$hooksDir = ".git\hooks"
if (-not (Test-Path $hooksDir)) {
    New-Item -ItemType Directory -Path $hooksDir -Force | Out-Null
}

# Pre-commit hook
$preCommit = @'
#!/bin/sh
# SSHH pre-commit hook
# Runs formatting, vetting, and tests before each commit

echo ""
echo "  Running pre-commit checks..."
echo ""

# Check Go formatting
echo "  [1/4] Checking formatting..."
UNFORMATTED=$(gofmt -l . 2>&1)
if [ -n "$UNFORMATTED" ]; then
    echo "  ERROR: The following files are not formatted:"
    echo "$UNFORMATTED"
    echo ""
    echo "  Run: gofmt -s -w ."
    exit 1
fi
echo "    OK"

# Run go vet
echo "  [2/4] Running go vet..."
if ! go vet ./... 2>&1; then
    echo "  ERROR: go vet failed"
    exit 1
fi
echo "    OK"

# Run tests
echo "  [3/4] Running tests..."
if ! go test -short ./... 2>&1; then
    echo "  ERROR: Tests failed"
    exit 1
fi
echo "    OK"

# Check for secrets (basic check)
echo "  [4/4] Scanning for secrets..."
if git diff --cached --name-only | xargs grep -l "PRIVATE KEY" 2>/dev/null; then
    echo "  WARNING: Possible private key detected in staged files!"
    echo "  Please review before committing."
    exit 1
fi
echo "    OK"

echo ""
echo "  All pre-commit checks passed!"
echo ""
'@

$preCommitPath = Join-Path $hooksDir "pre-commit"
$preCommit | Out-File -FilePath $preCommitPath -Encoding utf8 -NoNewline

# Commit-msg hook (conventional commits)
$commitMsg = @'
#!/bin/sh
# SSHH commit-msg hook
# Validates commit messages follow Conventional Commits format

commit_msg=$(cat "$1")

# Pattern: type(scope): description OR type: description
pattern="^(feat|fix|docs|style|refactor|perf|test|chore|ci|build|revert)(\(.+\))?: .{1,}"

if ! echo "$commit_msg" | grep -qE "$pattern"; then
    echo ""
    echo "  ERROR: Invalid commit message format!"
    echo ""
    echo "  Expected: <type>(<scope>): <description>"
    echo ""
    echo "  Types: feat, fix, docs, style, refactor, perf, test, chore, ci, build, revert"
    echo ""
    echo "  Examples:"
    echo "    feat(vault): add credential tags"
    echo "    fix: handle empty password input"
    echo "    docs(readme): update installation guide"
    echo ""
    exit 1
fi
'@

$commitMsgPath = Join-Path $hooksDir "commit-msg"
$commitMsg | Out-File -FilePath $commitMsgPath -Encoding utf8 -NoNewline

Write-Host "  Installed hooks:" -ForegroundColor Green
Write-Host "    pre-commit  - Format, vet, test, secret scan" -ForegroundColor White
Write-Host "    commit-msg  - Conventional Commits validation" -ForegroundColor White
Write-Host ""
Write-Host "  Git hooks are ready!" -ForegroundColor Green
Write-Host ""
