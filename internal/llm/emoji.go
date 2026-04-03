package llm

import "strings"

// Gitmoji maps conventional commit types to emojis.
var Gitmoji = map[string]string{
	"feat":     "✨",
	"fix":      "🐛",
	"docs":     "📝",
	"style":    "💄",
	"refactor": "♻️",
	"perf":     "⚡",
	"test":     "✅",
	"build":    "📦",
	"ci":       "👷",
	"chore":    "🔧",
	"revert":   "⏪",
}

// AddEmoji prepends a gitmoji to a conventional commit message.
func AddEmoji(msg string) string {
	for commitType, emoji := range Gitmoji {
		if strings.HasPrefix(msg, commitType+"(") || strings.HasPrefix(msg, commitType+":") {
			return emoji + " " + msg
		}
	}
	return msg
}
