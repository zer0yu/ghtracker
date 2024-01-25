package cmd

import (
	"github.com/projectdiscovery/gologger"
	"github.com/zer0yu/ghtracker/pkg/runner"
	"os"

	"github.com/spf13/cobra"
	"github.com/zer0yu/ghtracker/pkg/options"
)

var ghoptions = options.GHTopDepOptions{}

var rootCmd = &cobra.Command{
	Use:   "ghtracker",
	Short: "CLI tool for tracking dependents repositories and sorting result by Stars",
	Long:  ``,
	Run:   GHTracker,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&ghoptions.Version, "version", "v", false, "Show version of ghtracker")
	rootCmd.PersistentFlags().StringVarP(&ghoptions.URL, "url", "u", "", "URL to process")
	rootCmd.PersistentFlags().BoolVarP(&ghoptions.Repositories, "repositories", "r", false,
		"Sort repositories or packages (default repositories)")
	rootCmd.PersistentFlags().BoolVarP(&ghoptions.Table, "table", "t", false, "View mode")
	rootCmd.PersistentFlags().BoolVarP(&ghoptions.Description, "description", "d", false,
		"Show description of packages or repositories (performs additional request per repository)")
	rootCmd.PersistentFlags().IntVarP(&ghoptions.Rows, "rows", "o", 10, "Number of showing repositories (default=10)")
	rootCmd.PersistentFlags().IntVarP(&ghoptions.MinStar, "minstar", "m", 5, "Minimum number of stars (default=5)")
	rootCmd.PersistentFlags().StringVarP(&ghoptions.Search, "search", "s", "",
		"search code at dependents (repositories or packages)")
	rootCmd.PersistentFlags().StringVarP(&ghoptions.Token, "token", "k", os.Getenv("GHTOPDEP_TOKEN"), "GitHub token")
	rootCmd.PersistentFlags().StringVarP(&ghoptions.OutputFile, "output", "f", "", "File to write output to (optional)")
}

func GHTracker(_ *cobra.Command, _ []string) {
	newRunner, err := runner.NewRunner(&ghoptions)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	err = newRunner.RunGHCrawler()
	if err != nil {
		gologger.Fatal().Msgf("Could not run fuzz engine: %s\n", err)
	}
}
