package main

import (
	"fmt"

	"github.com/MaxHalford/xgp"
	"github.com/MaxHalford/xgp/tree"
	"github.com/urfave/cli"
)

var toEquationCmd = cli.Command{
	Name:  "toeq",
	Usage: "Produces an equation representation of a program",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "program, p",
			Value: "program.json",
			Usage: "Path to the program",
		},
	},
	Action: func(c *cli.Context) error {
		// Load the program
		program, err := xgp.LoadProgramFromJSON(c.String("program"))
		if err != nil {
			return exitCLI(err)
		}
		// Make the equation representation
		var str = tree.EquationDisplay{}.Apply(program.Root)
		// Output in the shell
		fmt.Println(str)
		return nil
	},
}
