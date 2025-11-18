# Troubleshooting Guide

A systematic approach to diagnosing and fixing issues, inspired by lessons learned from debugging CI failures.

## Core Principles

### 1. Don't Rush to Fix Based on First Impressions
- **The Problem:** It's tempting to fix the first error you see, but this often leads to treating symptoms rather than root causes.
- **The Solution:** Take time to understand the complete picture before implementing fixes.

### 2. Deep Analysis Over Quick Fixes
- Read error messages completely, not just the first line
- Look at the full context, not just the failing step
- Trace the execution flow to understand what's happening
- Ask "why" multiple times to get to the root cause

### 3. Verify Your Assumptions
- Don't assume you know what a parameter does - verify it
- Check documentation for the actual behavior
- Test your understanding before implementing changes

## Troubleshooting Process

### Step 1: Gather Complete Information
```bash
# Get detailed logs from failed runs
gh run view <run-id> --log | grep -A 20 -B 5 "Error"

# Look at the full context, not just error lines
gh run view <run-id> --log | less

# Check recent changes that might have caused the issue
git log --oneline -10
git diff HEAD~1
```

### Step 2: Identify the Real Error
- Filter out noise (warnings, cache errors, etc.)
- Find the actual failure point
- Read the complete error message, not just the summary

**Example:** In our CI failure, there were multiple errors:
- ❌ "Unexpected input 'ruby-version-file'" - This was a distraction
- ✅ "File to import not found: ./color_schemes/auto" - This was the real issue

### Step 3: Understand the Root Cause
Ask yourself:
- Why did this error occur?
- What changed recently that could cause this?
- Is this error a symptom of a deeper issue?

**Example:** The Jekyll error wasn't about Jekyll itself - it was about a configuration option (`color_scheme: auto`) that was valid for one theme but not another.

### Step 4: Trace Execution Flow
For workflow failures:
- Read the workflow file step by step
- Note the order of operations
- Identify dependencies between steps
- Look for state changes that might affect later steps

**Example:** In the release workflow, we discovered:
1. Download packages → 2. Checkout code → 3. Use packages
   
The checkout step was wiping out the downloaded packages!

### Step 5: Implement and Verify
- Make targeted fixes based on root cause analysis
- Test the fix (locally if possible)
- Verify all related functionality still works
- Document what was learned

## Common Pitfalls to Avoid

### 1. Fixing Symptoms Instead of Causes
❌ **Bad:** "The build fails, let me add a retry"
✅ **Good:** "The build fails, let me understand why"

### 2. Making Multiple Changes at Once
❌ **Bad:** "Let me fix all these issues in one commit"
✅ **Good:** "Let me fix and test each issue separately"

### 3. Ignoring Warning Signs
❌ **Bad:** "The test passed, but there's a warning - I'll ignore it"
✅ **Good:** "Let me understand what this warning means"

### 4. Not Checking Dependencies
❌ **Bad:** "This step looks independent"
✅ **Good:** "What state does this step depend on? What does it modify?"

## Real-World Example: CI Failure Debug Session

### Initial State
- CI failing with "Unexpected input 'ruby-version-file'"
- Pages build failing with mysterious Jekyll error

### First Impression (Wrong)
"The ruby/setup-ruby action needs the ruby-version-file parameter changed"

### Deep Analysis Revealed
1. The ruby-version parameter was indeed wrong
2. **BUT** the real failure was in Jekyll build
3. Jekyll was trying to load `color_scheme: auto`
4. This color scheme existed in TeXt theme (previously used)
5. But doesn't exist in just-the-docs theme (currently used)
6. The config was never updated when themes were switched

### The Fix
Remove the unsupported `color_scheme: auto` configuration

### Bonus Issue Found
Release workflow had incorrect step ordering:
- Packages were downloaded
- Then checkout wiped them out
- Then release tried to use them (now gone!)

## Key Takeaways

1. **Read the full error, not just the summary**
   - Error summaries can be misleading
   - The real cause is often buried in details

2. **Understand the system's behavior**
   - How do the pieces interact?
   - What's the order of operations?
   - What state is passed between steps?

3. **Verify before implementing**
   - Check documentation
   - Understand what each change does
   - Test your assumptions

4. **Think critically about "obvious" solutions**
   - Is this fixing the root cause or a symptom?
   - Could there be a deeper issue?

5. **Learn from the process**
   - Document what you learned
   - Share insights with the team
   - Build better debugging instincts

## Tools and Commands

### GitHub CLI for CI Debugging
```bash
# List recent runs
gh run list --limit 10

# View specific run
gh run view <run-id>

# Get detailed logs
gh run view <run-id> --log

# Filter logs for errors (excluding noise)
gh run view <run-id> --log | grep -E "(Error|Failed)" | \
  grep -v "Failed to save: Unable to reserve cache" | \
  grep -v "Failed to restore"
```

### Git for Context
```bash
# What changed recently?
git log --oneline -10
git log --all --follow -p <file>

# Find when something was introduced
git log --all --oneline --grep="<search-term>"

# See what a specific commit changed
git show <commit-hash>
```

### Local Testing
```bash
# Test locally before pushing
make test
make build

# For workflows, use act (GitHub Actions locally)
act -l                    # List workflows
act push                  # Run push workflows locally
```

## Remember

> "The first version of the fix is rarely the right one. Take time to understand the problem deeply, and the solution will become obvious."

When in doubt:
1. Pause
2. Read the full context
3. Trace the flow
4. Understand the cause
5. Then fix

This systematic approach saves time in the long run and builds better understanding of the system.

