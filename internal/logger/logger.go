package logger

import (
    "fmt"
    "io"
    "os"
)

// Logger handles output formatting
type Logger struct {
    verbose bool
    quiet   bool
    out     io.Writer
    err     io.Writer
}

// NewLogger creates a new logger
func NewLogger(verbose, quiet bool) *Logger {
    return &Logger{
        verbose: verbose,
        quiet:   quiet,
        out:     os.Stdout,
        err:     os.Stderr,
    }
}

// LogAttempt logs the current attempt
func (l *Logger) LogAttempt(attempt, total int, command string, args []string) {
    if l.verbose {
        fmt.Fprintf(l.out, "Attempt %d/%d: %s %v\n", attempt, total, command, args)
    }
}

// LogSuccess logs a successful attempt
func (l *Logger) LogSuccess(attempt int) {
    if l.verbose {
        fmt.Fprintf(l.out, "[Success] Command succeeded on attempt %d\n", attempt)
    }
}

// LogFailure logs a failed attempt
func (l *Logger) LogFailure(attempt int, err error) {
    if l.verbose {
        fmt.Fprintf(l.err, "[Failed] Attempt %d failed: %v\n", attempt, err)
    }
}

// LogTimeout logs a timeout event
func (l *Logger) LogTimeout(duration string) {
    if !l.quiet {
        fmt.Fprintf(l.err, "Command timed out after %s\n", duration)
    }
}

// LogRetryDelay logs the retry delay
func (l *Logger) LogRetryDelay(delay string) {
    if !l.quiet {
        fmt.Fprintf(l.out, "Retrying in %s...\n", delay)
    }
}