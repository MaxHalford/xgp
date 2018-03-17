package cmd

import (
	"fmt"
	"os"
	"strings"

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

	toDOTCmd.Flags().StringVarP(&toDOTOutputName, "output", "", "program.dot", "path to the DOT file output")
	toDOTCmd.Flags().BoolVarP(&toDOTSave, "save", "", false, "save to a DOT file or not")
	toDOTCmd.Flags().BoolVarP(&toDOTShell, "shell", "", true, "output in the terminal or not")
}

var toDOTCmd = &cobra.Command{
	Use:   "todot",
	Short: "Produces a DOT language representation of a program",
	Long:  "Produces a DOT language representation of a program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		// Load the program
		var program xgp.Program
		if strings.Contains(args[0], `"`) {
			program, err = xgp.LoadProgramFromJSON(args[0])
			if err != nil {
				return err
			}
		} else {
			tree, err := tree.ParseCode(args[0])
			if err != nil {
				return err
			}
			program = xgp.Program{Tree: tree}
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
