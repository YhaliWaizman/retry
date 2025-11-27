package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var apiKey string

var rootCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry is a CLI tool for simple elegant and fast retrying of commands",
	Long:  `A fast and simple CLI tool to convert between currencies using exchangerate.host.`,
}

func Execute(key string) {

	apiKey = key
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
