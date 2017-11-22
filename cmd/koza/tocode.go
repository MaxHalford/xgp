package koza

import (
	"fmt"

	"github.com/MaxHalford/koza"
	"github.com/MaxHalford/koza/tree"
	"github.com/spf13/cobra"
)

var (
	toCodeProgramName string
)

func init() {
	RootCmd.AddCommand(toCode)

	toCode.Flags().StringVarP(&toCodeProgramName, "program", "p", "program.json", "path to the program")
}

var toCode = &cobra.Command{
	Use:   "toeq",
	Short: "Produces a code-like representation of a program",
	Long:  "Produces a code-like representation of a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the program
		program, err := koza.LoadProgramFromJSON(toCodeProgramName)
		if err != nil {
			return err
		}
		// Make the equation representation
		var str = tree.CodeDisplay{}.Apply(program.Tree)
		// Output in the shell
		fmt.Println(str)
		return nil
	},
}
