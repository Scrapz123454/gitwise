package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/aymenhmaidiwastaken/gitwise/internal/config"
	"github.com/spf13/cobra"
)

// loadConfig loads config and applies CLI flag overrides.
func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if p, _ := cmd.Flags().GetString("provider"); p != "" {
		cfg.Provider = p
	}
	if m, _ := cmd.Flags().GetString("model"); m != "" {
		cfg.Model = m
	}

	return cfg, nil
}

func contextWithTimeout(seconds int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

func cleanMsg(msg string) string {
	msg = strings.TrimSpace(msg)
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	lower := strings.ToLower(msg)
	for _, prefix := range []string{"commit message:", "commit:", "message:"} {
		if strings.HasPrefix(lower, prefix) {
			msg = strings.TrimSpace(msg[len(prefix):])
			lower = strings.ToLower(msg)
		}
	}
	return msg
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func splitString(s, sep string) []string {
	return strings.Split(s, sep)
}
