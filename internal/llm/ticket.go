package llm

import "regexp"

// Common ticket patterns: JIRA (PROJ-123), Linear (PROJ-123), GitHub (#123)
var ticketPatterns = []*regexp.Regexp{
	regexp.MustCompile(`([A-Z][A-Z0-9]+-\d+)`),  // JIRA/Linear style: PROJ-123
	regexp.MustCompile(`#(\d+)`),                  // GitHub issue: #123
}

// ExtractTicket tries to extract a ticket reference from a branch name.
func ExtractTicket(branch string) string {
	for _, pat := range ticketPatterns {
		if match := pat.FindString(branch); match != "" {
			return match
		}
	}
	return ""
}
