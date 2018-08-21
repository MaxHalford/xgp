package cmd

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/meta"
	"github.com/MaxHalford/xgp/op"
	"github.com/spf13/cobra"
)

type toDOTCmd struct {
	modelPath  string
	round      uint
	shell      bool
	save       bool
	outputPath string // Only applies if save is true

	*cobra.Command
}

func (c *toDOTCmd) run(cmd *cobra.Command, args []string) error {
	// Load the model
	sm, err := readModel(c.modelPath)
	if err != nil {
		return err
	}

	// Build the Graphviz representation
	var (
		disp = op.GraphvizDisplay{}
		str  string
	)
	switch sm.Flavor {
	case "vanilla":
		str = disp.Apply(sm.Model.(xgp.Program).Op)
	case "boosting":
		gb := sm.Model.(meta.GradientBoosting)
		if uint(len(gb.Programs)) < c.round+1 {
			return fmt.Errorf("Ensemble only contains %d programs", len(gb.Programs))
		}
		str = disp.Apply(gb.Programs[c.round].Op)
	default:
		return errUnknownFlavor{sm.Flavor}
	}

	// Output in the shell if instructed
	if c.shell {
		fmt.Println(str)
	}

	// Save the output if instructed
	if !c.save {
		return nil
	}
	file, err := os.Create(c.outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(str)
	file.Sync()
	return nil
}

func newToDOTCmd() *toDOTCmd {
	c := &toDOTCmd{}
	c.Command = &cobra.Command{
		Use:   "todot",
		Short: "Produces a DOT language representation of a program",
		Long:  "Produces a DOT language representation of a program",
		Args:  cobra.ExactArgs(0),
		RunE:  c.run,
	}

	c.Flags().StringVarP(&c.modelPath, "model", "", "model.json", "path to the model used to make predictions")
	c.Flags().UintVarP(&c.round, "round", "", 0, "position of the program in the ensemble")
	c.Flags().BoolVarP(&c.shell, "shell", "", true, "output in the terminal or not")
	c.Flags().BoolVarP(&c.save, "save", "", false, "save to a DOT file or not")
	c.Flags().StringVarP(&c.outputPath, "output", "", "program.dot", "path to the DOT file output")

	return c
}
