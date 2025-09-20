# Branch Protection Configuration

This document outlines the recommended branch protection settings for the `choochoo` repository to maintain code quality and security.

## Recommended Branch Protection Rules for `main` branch

### Required Status Checks
- ✅ **Require status checks to pass before merging**
- ✅ **Require branches to be up to date before merging**
- Required status checks:
  - `lint-and-validate` (from CI workflow)
  - `security-check` (from CI workflow)

### Pull Request Requirements
- ✅ **Require a pull request before merging**
- ✅ **Require approvals**: At least 1 reviewer
- ✅ **Dismiss stale pull request approvals when new commits are pushed**
- ✅ **Require review from code owners** (if CODEOWNERS file exists)

### Additional Restrictions
- ✅ **Restrict pushes that create files larger than 100 MB**
- ✅ **Require conversation resolution before merging**
- ✅ **Include administrators** (applies rules to repository admins)

### Administrative Settings
- ✅ **Allow force pushes**: Disabled
- ✅ **Allow deletions**: Disabled

## How to Configure

### Via GitHub Web Interface

1. Navigate to repository **Settings** → **Branches**
2. Click **Add rule** or edit existing rule for `main`
3. Configure the settings as outlined above

### Via GitHub CLI

```bash
# Install GitHub CLI if not already installed
# gh auth login

# Create branch protection rule
gh api repos/deedubs/choochoo/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["lint-and-validate","security-check"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field restrictions=null \
  --field allow_force_pushes=false \
  --field allow_deletions=false
```

### Via GitHub REST API

```bash
curl -X PUT \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  https://api.github.com/repos/deedubs/choochoo/branches/main/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": ["lint-and-validate", "security-check"]
    },
    "enforce_admins": true,
    "required_pull_request_reviews": {
      "required_approving_review_count": 1,
      "dismiss_stale_reviews": true
    },
    "restrictions": null,
    "allow_force_pushes": false,
    "allow_deletions": false
  }'
```

## Benefits

These branch protection settings provide:

1. **Code Quality**: All changes must pass CI checks
2. **Peer Review**: Human review of all changes
3. **Security**: Prevention of force pushes and accidental deletions
4. **Consistency**: Up-to-date branches ensure clean merge history
5. **Accountability**: Clear audit trail of who approved what changes

## CI Workflow

The repository includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that provides:

- **Validation**: Ensures README.md exists and has proper structure
- **Security**: Basic checks for sensitive files and secrets
- **Consistency**: Standardized checks for all pull requests

These checks serve as the required status checks for branch protection.

## Customization

Adjust the settings based on your team's needs:

- **Increase review count** for more critical repositories
- **Add additional status checks** as your CI/CD pipeline grows
- **Configure code owners** by adding a `.github/CODEOWNERS` file
- **Add custom checks** to the CI workflow as needed