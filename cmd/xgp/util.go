package xgp

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaxHalford/xgp/tree"
	"github.com/urfave/cli"
)

func fileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("No file named '%s'", file), 1)
	}
	return nil
}

func parseStringFuncs(str string) ([]tree.Operator, error) {
	var (
		strs  = strings.Split(str, ",")
		funcs = make([]tree.Operator, len(strs))
	)
	for i, s := range strs {
		var f, err = tree.GetFunction(s)
		if err != nil {
			return nil, err
		}
		funcs[i] = f
	}
	return funcs, nil
}
