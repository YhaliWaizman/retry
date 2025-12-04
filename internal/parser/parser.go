package parser

import (
	"strconv"
	"strings"
	
	"github.com/yhaliwaizman/retry/internal/config"
)

// ParseOnFlag parses the "on" flag into a slice of strings
func ParseOnFlag(onFlag string) []string {
	if onFlag == "" {
		return nil
	}
	return strings.Split(onFlag, ",")
}

func OnLogic(cfg *config.Config, exitCode int, output string) bool {
	onConditions := ParseOnFlag(cfg.On)
	if len(onConditions) == 0 {
		return false
	}
	exitCodeStr := strconv.Itoa(exitCode)
	for _, condition := range onConditions {
		condition = strings.TrimSpace(condition)
		if condition == exitCodeStr || strings.Contains(output, condition) {
			return true
		}
	}
	return false
}