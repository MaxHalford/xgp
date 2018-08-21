package cmd

import "fmt"

type errUnknownFlavor struct {
	mode string
}

func (e errUnknownFlavor) Error() string {
	return fmt.Sprintf("'%s' is not a recognized flavor, has to be one of ('vanilla', 'boosting')", e.mode)
}
