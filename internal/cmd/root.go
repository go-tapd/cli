package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-tapd/cli/internal/app"
	"github.com/spf13/cobra"
)

func Execute() int {
	rt := app.NewRuntime()
	rootCmd := newRootCmd(rt)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rt.Stderr, "Error:", err)
		return 1
	}

	return 0
}

func newRootCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tapd",
		Short:         "TAPD CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			configPath, err := app.ResolveConfigPath(rt.ConfigPath)
			if err != nil {
				return err
			}
			rt.ConfigPath = configPath
			return nil
		},
	}

	cmd.SetOut(rt.Stdout)
	cmd.SetErr(rt.Stderr)
	cmd.PersistentFlags().StringVar(&rt.ConfigPath, "config", "", "config file path")
	cmd.PersistentFlags().StringVar(&rt.BaseURLOverride, "base-url", "", "override TAPD API base URL")
	cmd.PersistentFlags().StringVarP(&rt.OutputFormat, "format", "F", "table", "output format: table|json")
	cmd.SetContext(signalContext())

	cmd.AddCommand(
		newLoginCmd(rt),
		newAuthCmd(rt),
		newWorkspaceCmd(rt),
		newStoryCmd(rt),
		newBugCmd(rt),
		newIterationCmd(rt),
		newTaskCmd(rt),
		newReportCmd(rt),
	)

	return cmd
}

func signalContext() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
	}()
	return ctx
}
