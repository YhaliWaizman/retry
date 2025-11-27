package command

import (
	"fmt"
    "os"
    "os/exec"
    "strconv"
    "time"

	"github.com/spf13/cobra"
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
		fmt.Printf("Attempt %d/%d: %s %v\n", i, times, runCmd, runArgs)

		c := exec.Command(runCmd, runArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		err := c.Run()
		if err == nil {
            fmt.Printf("[Success] Command succeeded on attempt %d\n", i)
            return nil
        }
		fmt.Printf("[Failed] Attempt %d failed: %v\n", i, err)
        if i < times {
            fmt.Printf("Retrying in %v...\n", 2*time.Second)
            time.Sleep(2*time.Second)
        }
	}

	return fmt.Errorf("all attempts failed")
}




var rootCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry is a CLI tool for simple elegant and fast retrying of commands",
	Long:  `A fast and simple CLI tool to retry a given command multiple times with delays between attempts. Useful for handling flaky commands or transient errors by automatically retrying until success or until the maximum number of attempts is reached.`,
	RunE: runner,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}