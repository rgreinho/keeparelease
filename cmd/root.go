package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v28/github"
	"github.com/rgreinho/keeparelease/keeparelease"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// Version is used by the build system.
var Version string

// The log flag value.
var l string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "keeparelease",
	Short:   "Create beautiful GitHub releases.",
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		if err := setUpLogs(l); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the flags.
		extract, _ := cmd.Flags().GetBool("extract")
		file, _ := cmd.Flags().GetString("file")

		// Read the last release information from the Changelog.
		title, content, err := keeparelease.ReadChangelog(file)
		if err != nil {
			return fmt.Errorf("cannot read the changelog file %q: %s", file, err)
		}

		// If extract only, simply display the content of the last release.
		if extract {
			fmt.Println(content)
			return nil
		}

		// Get the owner and repository.
		owner, repository := keeparelease.GetInfo()
		if owner == "" || repository == "" {
			return fmt.Errorf("cannot extract project information: owner:%q, repository:%q", owner, repository)
		}

		// Get token.
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			return fmt.Errorf("cannot retrieve the authentication token")
		}

		// Get the tag.
		tag, err := cmd.Flags().GetString("tag")
		if tag == "" {
			tag, err = keeparelease.GetTag()
		}
		if err != nil || tag == "" {
			return fmt.Errorf("cannot retrieve the tag: %s", err)
		}
		log.Debugf("tag %q found", tag)

		// Prepare the client.
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		gh := github.NewClient(tc)

		// Prepare the release.
		releaseOptions := &github.RepositoryRelease{
			TagName: &tag,
			Name:    &title,
			Body:    &content,
		}

		// Create the release.
		release, _, err := gh.Repositories.CreateRelease(ctx, owner, repository, releaseOptions)
		log.Debugf("release options: %s", releaseOptions.String())
		if err != nil {
			return fmt.Errorf("cannot create the release: %s", err)
		}

		// Upload assets.
		assets, _ := cmd.Flags().GetStringArray("attach")
		for _, asset := range assets {
			// Validate the asset.
			info, err := os.Stat(asset)
			if err != nil {
				return fmt.Errorf("cannot validate the asset file %q", asset)
			}
			if info.Size() == 0 {
				return fmt.Errorf("cannot upload an asset of size 0 byte")
			}

			// Prepare options.
			uploadOptions := &github.UploadOptions{
				Name: filepath.Base(asset),
			}

			// Open the asset file.
			assetFile, err := os.Open(asset)
			if err != nil {
				return fmt.Errorf("cannot open file %q", asset)
			}

			// Attach the asset.
			log.Debugf("attaching asset %q", assetFile.Name())
			_, _, err = gh.Repositories.UploadReleaseAsset(ctx, owner, repository, release.GetID(), uploadOptions, assetFile)
			if err != nil {
				return fmt.Errorf("cannot upload asset %q: %s", asset, err)
			}
		}
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
	rootCmd.Flags().StringArrayP("attach", "a", []string{}, "Specify the assets to include into the release")
	rootCmd.Flags().BoolP("extract", "x", false, "Only extract the last release information")
	rootCmd.Flags().StringP("file", "f", "CHANGELOG.md", "Specify a changelog file")
	rootCmd.PersistentFlags().StringVarP(&l, "log", "l", "", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.Flags().StringP("tag", "t", "", "Use a specific tag")
}

// func uploadAssets(gh *github.Client, release *github.Release, assets []string) {
// 	for _, asset := range assets {
// 		var label string
// 		parts := strings.SplitN(asset, "#", 2)
// 		asset = parts[0]
// 		if len(parts) > 1 {
// 			label = parts[1]
// 		}

// 		for _, existingAsset := range release.Assets {
// 			if existingAsset.Name == filepath.Base(asset) {
// 				err := gh.DeleteReleaseAsset(&existingAsset)
// 				utils.Check(err)
// 				break
// 			}
// 		}
// 		ui.Errorf("Attaching release asset `%s'...\n", asset)
// 		_, err := gh.UploadReleaseAsset(release, asset, label)
// 		utils.Check(err)
// 	}
// }

// SetUpLogs sets the log level.
func setUpLogs(level string) error {
	// Read the log level
	//  1. from the CLI first
	//  2. then the ENV vars
	//  3. then use the default value.
	if level == "" {
		level = os.Getenv("ARMAMENT_LOG_LEVEL")
		if level == "" {
			level = logrus.WarnLevel.String()
		}
	}

	// Parse the log level.
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Set the log level.
	logrus.SetLevel(lvl)
	return nil
}
