package command

import (
    "context"
    "fmt"
    "strconv"
    "time"

    "github.com/yhaliwaizman/retry/internal/config"
    "github.com/yhaliwaizman/retry/internal/executor"
    "github.com/yhaliwaizman/retry/internal/logger"
)

// Runner orchestrates the retry logic
type Runner struct {
    config   *config.Config
    logger   *logger.Logger
    executor *executor.Executor
}

// NewRunner creates a new runner
func NewRunner(cfg *config.Config, log *logger.Logger, exec *executor.Executor) *Runner {
    return &Runner{
        config:   cfg,
        logger:   log,
        executor: exec,
    }
}

// Run executes the retry logic
func (r *Runner) Run(args []string) error {
    if len(args) < 2 {
        return fmt.Errorf("usage: retry <times> <command> [args...]")
    }

    times, err := r.parseAttempts(args[0])
    if err != nil {
        return err
    }

    return r.executeWithRetry(times, args[1], args[2:])
}

// parseAttempts validates and parses the number of retry attempts
func (r *Runner) parseAttempts(attemptsStr string) (int, error) {
    times, err := strconv.Atoi(attemptsStr)
    if err != nil || times < 1 {
        return 0, fmt.Errorf("invalid number of times: %s (must be a positive integer)", attemptsStr)
    }
    return times, nil
}

// executeWithRetry contains the main retry logic
func (r *Runner) executeWithRetry(times int, command string, args []string) error {
    var lastErr error

    for attempt := 1; attempt <= times; attempt++ {
        ctx, cancel := r.createContext(r.config.CommandTimeout)
        
        r.logger.LogAttempt(attempt, times, command, args)
		
        result := r.executor.Execute(ctx, command, args)
        
        cancel()

        if result.Err == nil {
            r.logger.LogSuccess(attempt)
            return nil
        }

        lastErr = result.Err
        r.logger.LogFailure(attempt, result.Err)

        if result.TimedOut {
            r.logger.LogTimeout(r.config.CommandTimeout.String())
        }

        if attempt < times {
            r.delayBeforeRetry()
        }
    }

    return fmt.Errorf("all %d attempts failed: %w", times, lastErr)
}

// createContext creates a context with timeout if configured
func (r *Runner) createContext(timeout time.Duration) (context.Context, context.CancelFunc) {
    if timeout > 0 {
        return context.WithTimeout(context.Background(), timeout)
    }
    return context.WithCancel(context.Background())
}

// delayBeforeRetry waits before the next retry attempt
func (r *Runner) delayBeforeRetry() {
    r.logger.LogRetryDelay(r.config.Delay.String())
    time.Sleep(r.config.Delay)
}