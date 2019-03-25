package cmd

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/genny/command"
	"github.com/gobuffalo/genny/genny/new"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/release/genny/initgen"
	"github.com/spf13/cobra"
)

var commandOptions = struct {
	Cmd    *command.Options
	New    *new.Options
	dryRun bool
}{
	Cmd: &command.Options{},
	New: &new.Options{},
}

// commandCmd represents the command command
var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "generates a command genny stub",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if commandOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		app := meta.New(".")

		copts := commandOptions.Cmd
		copts.App = app

		var name string
		if len(args) > 0 {
			name = args[0]
		}
		copts.Name = name

		r.Root = filepath.Join(app.Root, name)

		if err := r.WithNew(command.New(copts)); err != nil {
			return err
		}

		nopts := commandOptions.New
		nopts.Name = name
		nopts.Prefix = copts.Prefix

		if err := r.WithNew(new.New(nopts)); err != nil {
			return err
		}

		iopts := &initgen.Options{
			MainFile: "main.go",
			Root:     r.Root,
		}
		g, err := initgen.New(iopts)
		if err != nil {
			return err
		}
		r.WithGroup(g)
		return r.Run()
	},
}

func init() {
	commandCmd.Flags().BoolVarP(&commandOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	commandCmd.Flags().StringVarP(&commandOptions.Cmd.Prefix, "prefix", "p", "", "path prefix for the generator")
	rootCmd.AddCommand(commandCmd)
}
