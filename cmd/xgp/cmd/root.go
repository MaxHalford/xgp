package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "xgp",
	Short: "Machine learning tool based on genetic programming",
	Long:  "Machine learning tool based on genetic programming",
}

func init() {
	RootCmd.AddCommand(newFitCmd().Command)
	RootCmd.AddCommand(newPredictCmd().Command)
	RootCmd.AddCommand(newScoreCmd().Command)
	RootCmd.AddCommand(newToDOTCmd().Command)
}

// Execute RootCmd and catch error.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
