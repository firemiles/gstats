package main

import (
	"fmt"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/firemiles/gstats/pkg/model"
	"github.com/firemiles/gstats/pkg/parser"
	"github.com/firemiles/gstats/pkg/statistician"
	externalapi "github.com/firemiles/gstats/pkg/statistician/externalAPI"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

const (
	version = "0.1"

	excludeMatchPattern = "^" + generator.GenfilePrefix + ".*.go$"
)

var (
	inputDirOption string
)

var rootCmd = &cobra.Command{
	Use:   "gstats",
	Short: "gstats is a golang statistics tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		parsedSources, err := parser.New().ParseSourceDir(inputDirOption, "^.*.go$", excludeMatchPattern)
		if err != nil {
			klog.Errorf("Error parsing golang source in %s: %s", inputDirOption, err)
			return err
		}
		return runAllStatistician(inputDirOption, parsedSources)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&inputDirOption, "input-dir", "", "Directory to be statistics")
	rootCmd.MarkPersistentFlagRequired("input-dir")

	rootCmd.SetVersionTemplate(version)
}

func runAllStatistician(inputDir string, parsedSources model.ParsedSources) error {
	for name, s := range map[string]statistician.Statistician{
		"external-api": externalapi.NewStatistician(),
	} {
		infos, err := s.Statistics(inputDir, parsedSources)
		if err != nil {
			klog.Errorf("Error statistics module %s: %s", name, err)
			return err
		}
		fmt.Println(s.PrettyFormat(infos))
	}
	return nil
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
