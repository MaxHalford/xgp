package main

import (
	"os"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/urfave/cli"
)

var toDOTCmd = cli.Command{
	Name:  "todot",
	Usage: "Creates a .dot file from a program",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "program, p",
			Value: "program.json",
			Usage: "Path to the program",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "program.dot",
			Usage: "Path for the output file",
		},
	},
	Action: func(c *cli.Context) error {
		// Load the program
		program, err := xgp.LoadProgramFromJSON(c.String("program"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		// Create the output file
		file, err := os.Create(c.String("output"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer file.Close()
		// Make the Graphviz representation
		var graphviz = tree.GraphvizDisplay{}
		file.WriteString(graphviz.Apply(program.Root))
		file.Sync()
		return nil
	},
}
