package xgp

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/spf13/cobra"
)

var (
	toDOTProgramName string
	toDOTSave        bool
	toDOTShell       bool
	toDOTOutputName  string
)

func init() {
	RootCmd.AddCommand(toDOTCmd)

	toDOTCmd.Flags().StringVarP(&toDOTOutputName, "output", "o", "program.dot", "Path for the output file")
	toDOTCmd.Flags().StringVarP(&toDOTProgramName, "program", "p", "program.json", "Path to the program")
	toDOTCmd.Flags().BoolVarP(&toDOTSave, "save", "c", false, "Save to a file or not")
	toDOTCmd.Flags().BoolVarP(&toDOTShell, "shell", "t", true, "Output in the terminal or not")
}

var toDOTCmd = &cobra.Command{
	Use:   "todot",
	Short: "Produces a DOT language representation of a program",
	Long:  "Produces a DOT language representation of a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the program
		program, err := xgp.LoadProgramFromJSON(toDOTProgramName)
		if err != nil {
			return err
		}
		// Make the Graphviz representation
		var str = tree.GraphvizDisplay{}.Apply(program.Root)
		// Output in the shell
		if toDOTShell {
			fmt.Println(str)
		}
		// Create the output file
		if toDOTSave {
			return nil
		}
		file, err := os.Create(toDOTOutputName)
		if err != nil {
			return err
		}
		defer file.Close()
		file.WriteString(str)
		file.Sync()
		return nil
	},
}
