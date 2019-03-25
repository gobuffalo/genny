package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/genny/command/widget/genny/widget"
	"github.com/gobuffalo/logger"
	"github.com/spf13/cobra"
)

var rootOptions = struct {
	*widget.Options
	dryRun  bool
	verbose bool
}{
	Options: &widget.Options{},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "widget",
	Short: "A brief description of your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		run := genny.WetRunner(ctx)
		if rootOptions.dryRun {
			run = genny.DryRunner(ctx)
		}

		if rootOptions.verbose {
			run.Logger = logger.New(logger.DebugLevel)
		}

		opts := rootOptions.Options
		if err := run.WithNew(widget.New(opts)); err != nil {
			return err
		}
		return run.Run()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&rootOptions.dryRun, "dry-run", "d", false, "performs a dry run")
	rootCmd.Flags().BoolVarP(&rootOptions.verbose, "verbose", "v", false, "turns on debug logging")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
