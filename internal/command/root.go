package command

import (
    "fmt"
    "os"
    "time"

    "github.com/spf13/cobra"
    "github.com/yhaliwaizman/retry/internal/config"
    "github.com/yhaliwaizman/retry/internal/executor"
    "github.com/yhaliwaizman/retry/internal/logger"
)

var cfg = config.NewConfig()

var rootCmd = &cobra.Command{
    Use:   "retry <times> <command> [args...]",
    Short: "Retry is a CLI tool for simple, elegant, and fast retrying of commands",
    Long: `Retry is a CLI tool that executes a command multiple times until it succeeds
or reaches the maximum number of attempts. It supports configurable delays,
timeouts, and logging options.

Examples:
  retry 3 curl https://example.com
  retry --delay 2s 5 ./flaky-script.sh
  retry --command-timeout 5s 3 npm test`,
    RunE: runCommand,
    Args: cobra.MinimumNArgs(2),
}

func runCommand(cmd *cobra.Command, args []string) error {
    log := logger.NewLogger(cfg.Verbose, cfg.Quiet)
    exec := executor.NewExecutor()
    runner := NewRunner(cfg, log, exec)
    
    return runner.Run(args)
}

func init() {
    rootCmd.Flags().DurationVarP(&cfg.Delay, "delay", "d", time.Second,
        "Delay between retry attempts")
    rootCmd.Flags().BoolVarP(&cfg.Quiet, "quiet", "q", false,
        "Suppress retry messages")
    rootCmd.Flags().DurationVarP(&cfg.Timeout, "timeout", "t", 0,
        "Overall timeout for all attempts (0 = no timeout)")
    rootCmd.Flags().DurationVarP(&cfg.CommandTimeout, "command-timeout", "c", 0,
        "Timeout for each individual command execution (0 = no timeout)")
    rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false,
        "Enable verbose logging")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}