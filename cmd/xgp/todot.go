package koza

import (
	"fmt"
	"os"

	"github.com/MaxHalford/koza/tree"
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

	toDOTCmd.Flags().StringVarP(&toDOTOutputName, "output", "o", "program.dot", "path for the output file")
	toDOTCmd.Flags().StringVarP(&toDOTProgramName, "program", "p", "program.json", "path to the program")
	toDOTCmd.Flags().BoolVarP(&toDOTSave, "save", "c", false, "save to a file or not")
	toDOTCmd.Flags().BoolVarP(&toDOTShell, "shell", "t", true, "output in the terminal or not")
}

var toDOTCmd = &cobra.Command{
	Use:   "todot",
	Short: "Produces a DOT language representation of a program",
	Long:  "Produces a DOT language representation of a program",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the program
		program, err := koza.LoadProgramFromJSON(toDOTProgramName)
		if err != nil {
			return err
		}
		// Make the Graphviz representation
		var str = tree.GraphvizDisplay{}.Apply(program.Tree)
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
