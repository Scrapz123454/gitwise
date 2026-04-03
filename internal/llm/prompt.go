package llm

import (
	"fmt"
	"strings"
)

// CommitPrompt generates the prompt for commit message generation.
func CommitPrompt(diff string, convention string, scopeFromPath bool, scope string, language string, commitHistory string, ticket string) string {
	scopeInstruction := ""
	if scopeFromPath && scope != "" {
		scopeInstruction = fmt.Sprintf("\n- Suggested scope based on file paths: %s", scope)
	}

	conventionRules := ""
	switch convention {
	case "conventional", "angular":
		conventionRules = `- Use format: type(scope): description
- Types: feat, fix, refactor, docs, test, chore, style, perf, ci, build
- Scope: infer from the most relevant directory/module changed` + scopeInstruction
	default:
		conventionRules = "- Write a clear, descriptive commit message in imperative mood"
	}

	langInstruction := ""
	if language != "" && language != "english" {
		langInstruction = fmt.Sprintf("\n- Write the commit message in %s", language)
	}

	styleInstruction := ""
	if commitHistory != "" {
		styleInstruction = fmt.Sprintf(`

Match the style and tone of these recent commit messages from the repository:
%s`, commitHistory)
	}

	ticketInstruction := ""
	if ticket != "" {
		ticketInstruction = fmt.Sprintf("\n- Reference ticket %s in the commit message footer", ticket)
	}

	return fmt.Sprintf(`Analyze the following git diff and generate a commit message.

Rules:
%s
- Description: imperative mood, lowercase, no period, max 72 chars
- If changes are complex, add a body separated by a blank line with bullet points
- If there are breaking changes, add a "BREAKING CHANGE:" footer%s%s%s
- Output ONLY the commit message, nothing else — no explanations, no markdown fences
%s
Diff:
%s`, conventionRules, langInstruction, ticketInstruction, styleInstruction, "", diff)
}

// MultiCommitPrompt generates a prompt requesting multiple commit message options.
func MultiCommitPrompt(diff string, convention string, scopeFromPath bool, scope string, count int, language string, commitHistory string, ticket string) string {
	scopeInstruction := ""
	if scopeFromPath && scope != "" {
		scopeInstruction = fmt.Sprintf("\n- Suggested scope based on file paths: %s", scope)
	}

	conventionRules := ""
	switch convention {
	case "conventional", "angular":
		conventionRules = `- Use format: type(scope): description
- Types: feat, fix, refactor, docs, test, chore, style, perf, ci, build
- Scope: infer from the most relevant directory/module changed` + scopeInstruction
	default:
		conventionRules = "- Write a clear, descriptive commit message in imperative mood"
	}

	langInstruction := ""
	if language != "" && language != "english" {
		langInstruction = fmt.Sprintf("\n- Write all commit messages in %s", language)
	}

	styleInstruction := ""
	if commitHistory != "" {
		styleInstruction = fmt.Sprintf(`

Match the style and tone of these recent commit messages from the repository:
%s`, commitHistory)
	}

	ticketInstruction := ""
	if ticket != "" {
		ticketInstruction = fmt.Sprintf("\n- Reference ticket %s in each commit message footer", ticket)
	}

	return fmt.Sprintf(`Analyze the following git diff and generate %d different commit message options.

Rules:
%s
- Description: imperative mood, lowercase, no period, max 72 chars
- If changes are complex, add a body separated by a blank line with bullet points
- If there are breaking changes, add a "BREAKING CHANGE:" footer%s%s%s
- Output ONLY the commit messages, each separated by the line "---"
- No explanations, no markdown fences, no numbering
- Each option should have a different perspective or emphasis

Diff:
%s`, count, conventionRules, langInstruction, ticketInstruction, styleInstruction, diff)
}

// PRPrompt generates the prompt for PR description generation.
func PRPrompt(branch, base, commits, diff, template string) string {
	templateInstruction := ""
	if template != "" {
		templateInstruction = fmt.Sprintf(`

Use the following PR template as the output format. Fill in each section based on the changes:

%s`, template)
	}

	defaultFormat := ""
	if template == "" {
		defaultFormat = `

Output the PR description in this exact format (no other text):

## Summary
Brief description of what this PR does and why.

## Changes
- Bullet list of key changes with file references

## Testing
- How to test these changes

## Breaking Changes
- List any breaking changes (or "None")`
	}

	return fmt.Sprintf(`Analyze the following commits and diffs for a pull request from branch '%s' to '%s'.
Generate a structured PR description.

Commits:
%s

Combined diff:
%s%s%s`, branch, base, commits, diff, templateInstruction, defaultFormat)
}

// DiffSummaryPrompt generates a prompt for summarizing a diff in plain English.
func DiffSummaryPrompt(diff string) string {
	return fmt.Sprintf(`Summarize the following git diff in plain English. Be concise and focus on what changed and why it matters. Use bullet points for multiple changes. Do not include code — just describe the changes.

Diff:
%s`, diff)
}

// ReviewPrompt generates a prompt for AI code review.
func ReviewPrompt(diff string) string {
	return fmt.Sprintf(`You are a senior software engineer performing a code review. Analyze the following git diff and provide feedback.

Focus on:
1. Bugs or logic errors
2. Security vulnerabilities (injection, XSS, auth issues, secrets in code)
3. Performance issues
4. Code style and readability concerns
5. Missing error handling
6. Potential edge cases

For each issue found, output:
- **File**: the file path
- **Line**: approximate line number or context
- **Severity**: critical / warning / suggestion
- **Issue**: description of the problem
- **Fix**: suggested fix

If the code looks good, say so briefly. Do not repeat the diff back.

Diff:
%s`, diff)
}

// ChangelogPrompt generates a prompt for changelog generation.
func ChangelogPrompt(commits string, fromTag, toTag string) string {
	tagRange := ""
	if fromTag != "" && toTag != "" {
		tagRange = fmt.Sprintf(" from %s to %s", fromTag, toTag)
	} else if toTag != "" {
		tagRange = fmt.Sprintf(" for %s", toTag)
	}

	return fmt.Sprintf(`Generate a CHANGELOG entry%s from the following conventional commits.

Group changes by type:
### Added (feat)
### Fixed (fix)
### Changed (refactor, perf)
### Documentation (docs)
### Other (chore, ci, build, test, style)

Rules:
- Write in past tense
- Each entry should be a single clear sentence
- Include scope in parentheses if present
- Skip merge commits
- Skip entries with no user-visible impact

Commits:
%s`, tagRange, commits)
}

// LabelPrompt generates a prompt for suggesting PR labels.
func LabelPrompt(files []string, diff string) string {
	fileList := strings.Join(files, "\n")
	return fmt.Sprintf(`Based on the following changed files and diff, suggest GitHub PR labels.

Choose from these common labels (pick all that apply):
bug, feature, enhancement, documentation, refactor, performance, security, breaking-change, dependencies, frontend, backend, api, database, infrastructure, ci, tests, config

Output ONLY the labels, one per line, nothing else.

Changed files:
%s

Diff summary (first 200 lines):
%s`, fileList, truncate(diff, 200))
}

func truncate(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return strings.Join(lines, "\n")
}

// LintPrompt generates a prompt for linting commit messages.
func LintPrompt(commits string) string {
	return fmt.Sprintf(`Analyze the following commit messages for conventional commit format compliance.

Rules to check:
1. Format: type(scope): description
2. Valid types: feat, fix, refactor, docs, test, chore, style, perf, ci, build, revert
3. Description: imperative mood, lowercase start, no trailing period, max 72 chars
4. Body (if present): separated by blank line, wrapped at 72 chars
5. Footer (if present): proper format (e.g., BREAKING CHANGE:, Refs:)

For each commit, output:
- The commit hash and message
- PASS or FAIL with specific violations
- Suggested fix if FAIL

Commits:
%s`, commits)
}
