package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "xgp"
	app.Usage = "Genetic programming for machine learning tasks"
	app.Compiled = time.Now()

	app.Commands = []cli.Command{
		fitCmd,
		predictCmd,
		toDOTCmd,
	}

	app.Run(os.Args)
}
