package main

import (
	"fmt"
	"os"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/urfave/cli"
)

var toDOTCmd = cli.Command{
	Name:  "todot",
	Usage: "Produces a DOT language representation of a program",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "program, p",
			Value: "program.json",
			Usage: "Path to the program",
		},
		cli.BoolTFlag{
			Name:  "shell, sh",
			Usage: "Output in the shell or not",
		},
		cli.BoolFlag{
			Name:  "save, s",
			Usage: "Save to a file or not",
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
			return exitCLI(err)
		}
		// Make the Graphviz representation
		var str = tree.GraphvizDisplay{}.Apply(program.Root)
		// Output in the shell
		if c.Bool("sh") {
			fmt.Println(str)
		}
		// Create the output file
		if !c.Bool("save") {
			return nil
		}
		file, err := os.Create(c.String("output"))
		if err != nil {
			return exitCLI(err)
		}
		defer file.Close()
		file.WriteString(str)
		file.Sync()
		return nil
	},
}
