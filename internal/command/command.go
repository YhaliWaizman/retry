package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	verbose        bool
	delay          time.Duration
	quiet          bool
	timeout        time.Duration
	commandTimeout time.Duration
)

func runner(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: retry <times> <command> [args...]")
	}
	times, err := strconv.Atoi(args[0])
	if err != nil || times < 1 {
		return fmt.Errorf("invalid number of times: %s", args[0])
	}

	runCmd := args[1]
	runArgs := args[2:]

	for i := 1; i <= times; i++ {
		var ctx context.Context
		var cancel context.CancelFunc
		if commandTimeout > 0 {
			ctx, cancel = context.WithTimeout(context.Background(), commandTimeout)
		} else {
			ctx, cancel = context.WithCancel(context.Background())
		}
		defer cancel()
		if verbose {
			fmt.Printf("Attempt %d/%d: %s %v\n", i, times, runCmd, runArgs)
		}

		c := exec.CommandContext(ctx, runCmd, runArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		err := c.Run()
		if err == nil {
			if verbose {
				fmt.Printf("[Success] Command succeeded on attempt %d\n", i)
			}
			return nil
		}

		if verbose {
			fmt.Printf("[Failed] Attempt %d failed: %v\n", i, err)
		}
		if i < times {
			if !quiet {
				fmt.Printf("Retrying in %s...\n", delay)
			}
			time.Sleep(delay)
		}
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Command timed out after %s\n", commandTimeout)
		}		
	}

	return fmt.Errorf("all attempts failed")
}

var rootCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry is a CLI tool for simple elegant and fast retrying of commands",
	Long:  `A fast and simple CLI tool to repeatedly execute a command until it succeeds or a specified number of attempts is reached. Supports configurable delays, per-command and overall timeouts, and verbose or quiet output modes.`,
	RunE:  runner,
}

func init() {
	// Flags usable BEFORE the positional args
	// rootCmd.Flags().BoolVarP(&backoff, "backoff", "b", false, "Enable backoff")
	// rootCmd.Flags().BoolVarP(&fatal, "fatal-on-change", "f", false, "Exit on change of error")
	// rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	// rootCmd.Flags().IntVarP(&on, "on", "o", 0, "Which error code to retry on (0 for any error)")
	rootCmd.Flags().DurationVarP(&delay, "delay", "d", time.Second, "Delay between attempts")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "t", 0, "Overall timeout for all attempts")
	rootCmd.Flags().DurationVarP(&commandTimeout, "command-timeout", "c", 0, "Timeout for each individual command execution")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
