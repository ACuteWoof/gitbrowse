// Gitbrowse: a simple web server for git.
// Copyright (C) 2026 Vithushan
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package template

import (
	"html"
	"os/exec"
	"regexp"
	"strings"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func GitShow(repo string, commit string) string {
	if !validateCommitOrTag(repo, commit) {
		return "<div class=\"error\">Invalid hash</div>"
	}
	cmd := exec.Command("git", "show", commit, "--date=format:%a %b %d %H:%M:%S %Y +0000")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	escaped := html.EscapeString(string(out))
	output := colorizeGitShowAdvanced(escaped)
	return "<div class=\"git-show\">" + output + "</div>"
}

func colorizeGitShowAdvanced(text string) string {
	patterns := []struct {
		name string
		re   *regexp.Regexp
		css  string
	}{
		// Commit header
		{"commit-header", regexp.MustCompile(`^commit ([a-f0-9]{40})(?: \((.*)\))?$`),
			`<span class="git-commit">commit</span> <span class="git-hash">$1</span> <span class="git-refs">$2</span>`},

		// Merge info
		{"merge", regexp.MustCompile(`^Merge: (([a-f0-9]{7} )+)$`),
			`<span class="git-merge">Merge:</span> <span class="git-parent-hashes">$1</span>`},

		// Author/Committer/Tagger (with proper HTML escaping)
		{"author", regexp.MustCompile(`^(Author|Commit(?:ter)?|tagger):\s+(.+?)\s+<(.+?)>\s+(.+)$`),
			`<span class="git-metadata-label">$1:</span> <span class="git-author-name">$2</span> &lt;<span class="git-author-email">$3</span>&gt; <span class="git-date">$4</span>`},

		// Date
		{"date", regexp.MustCompile(`^Date:\s+(.+)$`),
			`<span class="git-date-label">Date:</span> <span class="git-date">$1</span>`},

		// Stats line
		{"stats", regexp.MustCompile(`^ (\d+) file(?:s)? changed(?:, (\d+) insertion(?:s)?\(\+\))?(?:, (\d+) deletion(?:s)?\(-\))?$`),
			`<span class="git-stats"> $1 file(s) changed</span><span class="git-insertions">, $2 insertion(s)(+)</span><span class="git-deletions">, $3 deletion(s)(-)</span>`},

		// Diff stats with file names
		{"diff-stat", regexp.MustCompile(`^ (.*?)\s+\| (\d+) ([+-]*)$`),
			`<span class="git-stat-file">$1</span> <span class="git-stat-bar">|</span> <span class="git-stat-count">$2</span> <span class="git-stat-graph">$3</span>`},

		// File names in diff
		{"diff-files", regexp.MustCompile(`^diff --git a/(.+) b/(.+)$`),
			`<span class="git-diff-header">diff --git</span> <span class="git-file-a">a/$1</span> <span class="git-file-b">b/$2</span>`},

		// Index line with hash ranges (remove duplicate)
		{"index", regexp.MustCompile(`^index ([a-f0-9]+)\.\.([a-f0-9]+)(?: (\d+))?$`),
			`<span class="git-index">index</span> <span class="git-index-old">$1</span>..<span class="git-index-new">$2</span> <span class="git-index-mode">$3</span>`},

		// File mode changes
		{"file-mode", regexp.MustCompile(`^(old mode|new mode) (\d+)$`),
			`<span class="git-mode-label">$1</span> <span class="git-mode-value">$2</span>`},

		// Added/Deleted file
		{"file-status", regexp.MustCompile(`^(new|deleted) file mode (\d+)$`),
			`<span class="git-file-status-$1">$1 file mode</span> <span class="git-mode-value">$2</span>`},

		// Binary files
		{"binary", regexp.MustCompile(`^Binary files (a/.+) and (b/.+) differ$`),
			`<span class="git-binary">Binary files</span> <span class="git-binary-file-a">$1</span> and <span class="git-binary-file-b">$2</span> differ`},

		// Hunk header with optional function context
		{"hunk-header", regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@(?: (.+))?$`),
			`<span class="git-hunk">@@</span> <span class="git-hunk-old">-$1$2</span> <span class="git-hunk-new">+$3$4</span> <span class="git-hunk">@@</span><span class="git-function-context"> $5</span>`},

		// File paths (--- / +++ lines)
		{"file-paths", regexp.MustCompile(`^(---|\+\+\+) (a/|b/|/dev/null)?(.*)$`),
			`<span class="git-file-path-$1">$1</span> <span class="git-file-location">$2</span><span class="git-file-name">$3</span>`},

		// Tag object header
		{"tag-header", regexp.MustCompile(`^tag ([^\s]+)$`),
			`<span class="git-tag">tag</span> <span class="git-tag-name">$1</span>`},

		// Tag object
		{"tag-object", regexp.MustCompile(`^object ([a-f0-9]{40})$`),
			`<span class="git-tag-object-label">object</span> <span class="git-tag-object-hash">$1</span>`},

		// Tag type (commit/tag)
		{"tag-type", regexp.MustCompile(`^type (commit|tag)$`),
			`<span class="git-tag-type-label">type</span> <span class="git-tag-type">$1</span>`},

		// Tag message lines (start with 4 spaces) - careful with commit messages too
		{"message-line", regexp.MustCompile(`^    (.*)$`),
			`<span class="git-message-line">    $1</span>`},

		// Added lines (starts with + but not +++ which is file header)
		{"line-added", regexp.MustCompile(`^(\+)([^\+].*)$`),
			`<span class="git-line-added"><span class="git-diff-plus">+</span>$2</span>`},

		// Removed lines (starts with - but not --- which is file header)
		{"line-removed", regexp.MustCompile(`^(-)([^-].*)$`),
			`<span class="git-line-removed"><span class="git-diff-minus">-</span>$2</span>`},

		// Context lines (starts with space)
		{"line-context", regexp.MustCompile(`^( )(.*)$`),
			`<span class="git-line-context"><span class="git-diff-space">&nbsp;</span>$2</span>`},
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		matched := false
		for _, pattern := range patterns {
			if pattern.re.MatchString(line) {
				line = pattern.re.ReplaceAllString(line, pattern.css)
				matched = true
				break
			}
		}

		// If no specific pattern matched, check for diff lines
		if !matched {
			line = colorizeDiffLine(line)
		}

		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

func colorizeDiffLine(line string) string {
	// Check line prefixes for diff content
	switch {
	case strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- "):
		// Already handled by patterns
		return line
	case strings.HasPrefix(line, "+") && line != "+++":
		return `<span class="git-line-added"><span class="git-diff-plus">+</span>` + line[1:] + `</span>`
	case strings.HasPrefix(line, "-") && line != "---":
		return `<span class="git-line-removed"><span class="git-diff-minus">-</span>` + line[1:] + `</span>`
	case strings.HasPrefix(line, " "):
		return `<span class="git-line-context"><span class="git-diff-space">&nbsp;</span>` + line[1:] + `</span>`
	default:
		return line
	}
}

func validateCommitOrTag(repoPath, ref string) bool {
    repo, err := git.PlainOpen(repoPath)
    if err != nil {
        return false
    }
    
    // Try to resolve it as a revision (tag, branch, or commit hash)
    hash, err := repo.ResolveRevision(plumbing.Revision(ref))
    if err != nil {
        return false
    }
    
    // Make sure it points to a commit
    _, err = repo.CommitObject(*hash)
    return err == nil
}
