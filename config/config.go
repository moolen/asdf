package config

// DefaultTypeMap contains a mapping of types to groups
// which are used to render the changelog
var DefaultTypeMap = map[string]string{
	"feat":     "Feature",
	"breaking": "Breaking Changes",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverted",
	"docs":     "Documentation",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
}
