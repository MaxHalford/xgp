package xgp

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(todotCmd)

	todotCmd.Flags().StringVarP(&programName, "program", "p", "program.json", "Path to the program")
	todotCmd.Flags().BoolVarP(&shell, "shell", "t", true, "Output in the terminal or not")
	todotCmd.Flags().BoolVarP(&save, "save", "c", false, "Save to a file or not")
	todotCmd.Flags().StringVarP(&outputName, "output", "o", "program.dot", "Path for the output file")
}

var todotCmd = &cobra.Command{
	Use:   "todot",
	Short: "Produces a DOT language representation of a program",
	Long:  "Produces a DOT language representation of a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the program
		program, err := xgp.LoadProgramFromJSON(programName)
		if err != nil {
			return err
		}
		// Make the Graphviz representation
		var str = tree.GraphvizDisplay{}.Apply(program.Root)
		// Output in the shell
		if shell {
			fmt.Println(str)
		}
		// Create the output file
		if save {
			return nil
		}
		file, err := os.Create(outputName)
		if err != nil {
			return err
		}
		defer file.Close()
		file.WriteString(str)
		file.Sync()
		return nil
	},
}
