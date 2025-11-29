package executor

import (
    "context"
    "os"
    "os/exec"
)

// Result represents the outcome of a command execution
type Result struct {
    ExitCode int
    Err      error
    TimedOut bool
}

// Executor handles command execution
type Executor struct{}

// NewExecutor creates a new executor
func NewExecutor() *Executor {
    return &Executor{}
}

// Execute runs a single command with context
func (e *Executor) Execute(ctx context.Context, command string, args []string) *Result {
    cmd := exec.CommandContext(ctx, command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    
    result := &Result{
        Err:      err,
        TimedOut: ctx.Err() == context.DeadlineExceeded,
    }

    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            result.ExitCode = exitErr.ExitCode()
        } else {
            result.ExitCode = -1
        }
    }

    return result
}