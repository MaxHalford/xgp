package cmd

/*
var (
	toDOTRound      int
	toDOTSave       bool
	toDOTShell      bool
	toDOTOutputName string
)

func init() {
	RootCmd.AddCommand(toDOTCmd)

	toDOTCmd.Flags().IntVarP(&toDOTRound, "round", "", 0, "position of the program in the ensemble")
	toDOTCmd.Flags().StringVarP(&toDOTOutputName, "output", "", "program.dot", "path to the DOT file output")
	toDOTCmd.Flags().BoolVarP(&toDOTSave, "save", "", false, "save to a DOT file or not")
	toDOTCmd.Flags().BoolVarP(&toDOTShell, "shell", "", true, "output in the terminal or not")
}

var toDOTCmd = &cobra.Command{
	Use:   "todot",
	Short: "Produces a DOT language representation of a program",
	Long:  "Produces a DOT language representation of a program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Load the ensemble
		var (
			ensemble   meta.Ensemble
			bytes, err = ioutil.ReadFile(args[0])
		)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &ensemble)
		if err != nil {
			return err
		}

		// Make the Graphviz representation
		var str = op.GraphvizDisplay{}.Apply(ensemble.Programs[toDOTRound].Op)

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
}*/
