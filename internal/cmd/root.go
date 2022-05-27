package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/backoff/v2"
	"github.com/spf13/cobra"

	"github.com/jamillosantos/rdy"
)

var (
	timeoutString  = "30s"
	url            string
	verbose        bool
	verboseVerbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rdy [flags] [url]",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		timeout, err := time.ParseDuration(timeoutString)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "invalid timeout: ", timeoutString)
			os.Exit(1)
		}

		if len(args) == 0 {
			_, _ = fmt.Fprintln(os.Stderr, "error: missing url")
			_, _ = fmt.Fprintln(os.Stderr, "")
			_ = cmd.Help()
			os.Exit(1)
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
		defer cancelFunc()
		err = rdy.Wait(ctx, rdy.WaitRequest{
			URL:     args[0],
			Backoff: backoff.Constant(backoff.WithInterval(time.Second)),
			Reporter: &stderrReporter{
				verbose:        verbose,
				verboseVerbose: verboseVerbose,
			},
		})
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&timeoutString, "timeout", "t", timeoutString, "Timeout for waiting for readiness. Ex: 30s, 1m, 1h30m")
	rootCmd.Flags().BoolVar(&verbose, "v", verbose, "Verbose L1 mode on (request and response statuses)")
	rootCmd.Flags().BoolVar(&verboseVerbose, "vv", verboseVerbose, "Verbose L2 mode on (L1 + response body)")
}

type stderrReporter struct {
	verbose        bool
	verboseVerbose bool
}

func (s stderrReporter) print(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	_, _ = fmt.Fprintln(os.Stderr)
}

func (s stderrReporter) L1(ctx context.Context, format string, args ...interface{}) {
	if !s.verbose && !s.verboseVerbose {
		return
	}
	s.print(format, args...)
}

func (s stderrReporter) L2(ctx context.Context, format string, args ...interface{}) {
	if !s.verboseVerbose {
		return
	}
	s.print(format, args...)
}
