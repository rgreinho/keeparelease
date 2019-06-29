package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/hub/github"
	"github.com/github/hub/ui"
	"github.com/github/hub/utils"

	"github.com/rgreinho/keeparelease/keeparelease"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keeparelease",
	Short: "Create beautiful GitHub releases.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the flags.
		extract, err := cmd.Flags().GetBool("extract")
		utils.Check(err)

		// Read the last release information from the Changelog.
		title, content, err := keeparelease.ReadChangelog("")
		utils.Check(err)

		// If extract only, simply display the content of the last release.
		if extract {
			fmt.Println(content)
			return nil
		}

		// Inspect the local repo.
		localRepo, err := github.LocalRepo()
		utils.Check(err)
		project, nil := localRepo.MainProject()
		utils.Check(err)

		// Prepare the Hub client.
		gh := github.NewClient(project.Host)

		// Prepare the release.
		tag, err := cmd.Flags().GetString("tag")
		utils.Check(err)
		params := &github.Release{
			TagName: tag,
			Name:    title,
			Body:    content,
		}

		// Create the release.
		var release *github.Release
		release, err = gh.CreateRelease(project, params)
		utils.Check(err)

		// Upload assets.
		assets, err := cmd.Flags().GetStringArray("attach")
		utils.Check(err)
		uploadAssets(gh, release, assets)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keeparelease.yaml)")
	rootCmd.Flags().BoolP("extract", "x", false, "Only extract the last release information")
	rootCmd.Flags().StringP("tag", "t", "", "Use a specific tag")
	rootCmd.Flags().StringArrayP("attach", "a", []string{}, "Specify the assets to include into the release")
}

func uploadAssets(gh *github.Client, release *github.Release, assets []string) {
	for _, asset := range assets {
		var label string
		parts := strings.SplitN(asset, "#", 2)
		asset = parts[0]
		if len(parts) > 1 {
			label = parts[1]
		}

		for _, existingAsset := range release.Assets {
			if existingAsset.Name == filepath.Base(asset) {
				err := gh.DeleteReleaseAsset(&existingAsset)
				utils.Check(err)
				break
			}
		}
		ui.Errorf("Attaching release asset `%s'...\n", asset)
		_, err := gh.UploadReleaseAsset(release, asset, label)
		utils.Check(err)
	}
}
