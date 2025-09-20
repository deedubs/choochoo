# Branch Protection Setup

This document provides step-by-step instructions for configuring branch protection rules to require passing CI tests before merging pull requests.

## Required Configuration

To ensure code quality and prevent broken code from reaching the main branch, configure the following branch protection settings:

### Steps to Configure

1. **Navigate to Repository Settings**
   - Go to your GitHub repository
   - Click on **Settings** tab
   - Select **Branches** from the left sidebar

2. **Add or Edit Branch Protection Rule**
   - Click **Add rule** (or edit existing rule for `main`)
   - Set **Branch name pattern** to: `main`

3. **Configure Required Settings**
   Check the following options:

   #### Pull Request Requirements
   - ☑️ **Require a pull request before merging**
     - ☑️ **Require approvals** (recommended: 1 approval)
     - ☑️ **Dismiss stale PR approvals when new commits are pushed**

   #### Status Check Requirements
   - ☑️ **Require status checks to pass before merging**
   - ☑️ **Require branches to be up to date before merging**
   - In the status checks search box, add: **CI** (this matches our workflow name)

   #### Additional Protections
   - ☑️ **Restrict pushes that create files larger than 100MB**
   - ☑️ **Include administrators** (recommended for consistency)
   - ☑️ **Allow force pushes** (unchecked - keep disabled)
   - ☑️ **Allow deletions** (unchecked - keep disabled)

4. **Save the Rule**
   - Click **Create** or **Save changes**

## What This Accomplishes

Once configured, this branch protection rule will:

- **Prevent direct pushes** to the main branch
- **Require pull requests** for all changes
- **Mandate passing CI tests** before any merge
- **Ensure code review** before integration
- **Keep main branch stable** and deployable

## Testing the Setup

1. Create a test branch and make a small change
2. Open a pull request to `main`
3. Verify that the CI workflow runs automatically
4. Confirm that the merge button is disabled until CI passes
5. Check that the PR shows required status checks

## Troubleshooting

- **Status check not appearing**: Make sure your CI workflow has run at least once and shows up in the repository's Actions tab
- **Wrong status check name**: The status check name should match the job name in your workflow file (`test` or `CI`)
- **Administrators bypassing rules**: If "Include administrators" is unchecked, repository admins can bypass these restrictions

For more information, see the [GitHub documentation on branch protection rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches).