# GPG Key Management Scripts

Scripts for managing GPG keys used to sign the apt-bundle Debian repository.

## Automated Key Generation (Recommended)

Use the GitHub Actions workflow `.github/workflows/generate-gpg-key.yml` to automatically generate and manage GPG keys.

### Prerequisites

This workflow uses a GitHub App for secure, scoped access to repository secrets.

**1. Create a GitHub App:**

1. Go to https://github.com/organizations/apt-bundle/settings/apps/new
2. Configure the app:
   - **GitHub App name:** `apt-bundle-gpg-manager`
   - **Homepage URL:** `https://github.com/apt-bundle/apt-bundle`
   - **Webhook:**
     - ☐ Uncheck "Active" (we don't need webhooks)
   - **Repository permissions:**
     - **Secrets:** Read and write (required to update GPG secrets)
     - **Contents:** Read and write (required to commit public key)
   - **Where can this GitHub App be installed?**
     - ⚪ Only on this account
3. Click "Create GitHub App"

**2. Generate Private Key for the App:**

1. On the app settings page, scroll to "Private keys"
2. Click "Generate a private key"
3. Save the downloaded `.pem` file securely

**3. Install the App:**

1. On the app settings page, click "Install App" in the left sidebar
2. Click "Install" next to the `apt-bundle` organization
3. Choose "Only select repositories"
4. Select **only** `apt-bundle/apt-bundle`
5. Click "Install"

**4. Add App Credentials as Secrets:**

1. Note the **App ID** (found at top of app settings page)
2. Go to repository secrets: `https://github.com/apt-bundle/apt-bundle/settings/secrets/actions`
3. Add two secrets:

   **Secret 1:**
   - Name: `GPG_APP_ID`
   - Value: Your App ID (numeric, e.g., `123456`)

   **Secret 2:**
   - Name: `GPG_APP_PRIVATE_KEY`
   - Value: Contents of the `.pem` file you downloaded
     ```
     -----BEGIN RSA PRIVATE KEY-----
     ... (entire private key content) ...
     -----END RSA PRIVATE KEY-----
     ```

### Running the Workflow

1. Go to Actions → "Generate GPG Key for Repository Signing"
2. Click "Run workflow"
3. Leave "Regenerate key even if it exists?" unchecked (unless regenerating)
4. Click "Run workflow"

The workflow will:
- ✅ Generate a GPG key pair (RSA 4096-bit)
- ✅ Store the private key in `GPG_PRIVATE_KEY` secret
- ✅ Store the key ID in `GPG_KEY_ID` secret
- ✅ Commit the public key to `repo/public.key`

### Security Considerations

- ✅ **Scoped Access:** The GitHub App only has access to the `apt-bundle/apt-bundle` repository
- ✅ **Limited Permissions:** Only secrets (read/write) and contents (read/write) permissions
- ✅ **Short-lived Tokens:** The workflow generates short-lived tokens (1 hour) automatically
- ✅ **Auditable:** All app actions are logged in the organization audit log
- ⚠️ **App Private Key:** The `.pem` file is sensitive; store it securely as a GitHub Secret
- ℹ️ **GPG Private Key:** Only visible in GitHub Actions logs (private) and stored as a secret
- ✅ **Public Key:** Safe to commit and distribute

### Why GitHub App vs PAT?

| Feature | GitHub App | Classic PAT | Fine-grained PAT |
|---------|-----------|-------------|------------------|
| Scope | Single repo ✅ | All repos ❌ | Single repo ✅ |
| Token lifetime | 1 hour ✅ | Manual expiry ⚠️ | Manual expiry ⚠️ |
| Audit trail | Full ✅ | Limited ⚠️ | Good ✅ |
| Permissions | Minimal ✅ | Broad ❌ | Minimal ✅ |
| Setup complexity | Medium | Low | Low |
| Org approval | Automatic ✅ | Not needed | Required ⚠️ |

---

## Manual Key Generation (Alternative)

If you prefer to generate keys locally:

### Generate Key

```bash
gpg --batch --gen-key .github/workflows/scripts/gpg-keygen.conf
```

### Export Keys

```bash
.github/workflows/scripts/gpg-export-keys.sh
```

This creates:
- `public.key` - Commit this to the repository at `repo/public.key`
- `private.key` - Add to GitHub Secret `GPG_PRIVATE_KEY`, then delete

### Add to GitHub Secrets Manually

```bash
# First, generate a GitHub App token (you'll need GPG_APP_ID and GPG_APP_PRIVATE_KEY secrets set)
export GITHUB_TOKEN=$(gh api /app/installations/$(gh api /repos/apt-bundle/apt-bundle/installation -q .id)/access_tokens -X POST -q .token)

# Set private key
cat private.key | gh secret set GPG_PRIVATE_KEY

# Set key ID (replace with your actual key ID)
echo "YOUR_KEY_ID" | gh secret set GPG_KEY_ID

# Delete private key file
shred -u private.key
```

Or use the automated workflow instead of manual steps.

---

## Files

- `gpg-keygen.conf` - GPG key generation configuration
- `gpg-export-keys.sh` - Script to export keys after manual generation
- `README.md` - This file

## Key Details

- **Key Type:** RSA 4096-bit
- **Subkey Type:** RSA 4096-bit
- **Expiration:** None (repository signing keys typically don't expire)
- **Passphrase:** None (required for CI/CD automation)
- **Email:** maintainers@apt-bundle.org

