package xgp

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/spf13/cobra"
)

var (
	toEqProgramName string
)

func init() {
	RootCmd.AddCommand(toEq)

	toEq.Flags().StringVarP(&toEqProgramName, "program", "p", "program.json", "Path to the program")
}

var toEq = &cobra.Command{
	Use:   "toeq",
	Short: "Produces a equation representation of a program",
	Long:  "Produces a equation representation of a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the program
		program, err := xgp.LoadProgramFromJSON(toEqProgramName)
		if err != nil {
			return err
		}
		// Make the equation representation
		var str = tree.EquationDisplay{}.Apply(program.Root)
		// Output in the shell
		fmt.Println(str)
		return nil
	},
}
