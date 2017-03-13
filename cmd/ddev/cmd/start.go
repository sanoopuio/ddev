package cmd

import (
	"log"
	"strings"

	"github.com/drud/ddev/pkg/plugins/platform"
	"github.com/spf13/cobra"
)

const netName = "drud_default"

var (
	serviceType string
	webImage    string
	dbImage     string
	webImageTag string
	dbImageTag  string
	skipYAML    bool
	appClient   string
)

// StartCmd represents the add command
var StartCmd = &cobra.Command{
	Use:     "start [app_name] [environment_name]",
	Aliases: []string{"add"},
	Short:   "Start the local development environment for a site.",
	Long:    `Start initializes and configures the web server and database containers to provide a working environment for development.`,
	PreRun: func(cmd *cobra.Command, args []string) {

		client, err := platform.GetDockerClient()
		if err != nil {
			log.Fatal(err)
		}

		err = EnsureNetwork(client, netName)
		if err != nil {
			log.Fatal(err)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {

		app := platform.PluginMap[strings.ToLower(plugin)]

		opts := platform.AppOptions{
			Name:        activeApp,
			Environment: activeDeploy,
			Client:      appClient,
			WebImage:    webImage,
			WebImageTag: webImageTag,
			DbImage:     dbImage,
			DbImageTag:  dbImageTag,
			SkipYAML:    skipYAML,
			CFG:         cfg,
		}

		app.Init(opts)

		err := app.Start()
		if err != nil {
			Failed("Failed to start %s: %s", app.GetName(), err)
		}

		Success("Successfully added %s-%s", activeApp, activeDeploy)
		Success("Your application can be reached at: %s", app.URL())

	},
}

func init() {
	StartCmd.Flags().StringVarP(&webImage, "web-image", "", "", "Change the image used for the app's web server")
	StartCmd.Flags().StringVarP(&dbImage, "db-image", "", "", "Change the image used for the app's database server")
	StartCmd.Flags().StringVarP(&webImageTag, "web-image-tag", "", "", "Override the default web image tag")
	StartCmd.Flags().StringVarP(&dbImageTag, "db-image-tag", "", "", "Override the default web image tag")
	StartCmd.Flags().StringVarP(&plugin, "plugin", "p", "legacy", "Choose which plugin to use")
	StartCmd.Flags().BoolVarP(&skipYAML, "skip-yaml", "", false, "Skip creating the docker-compose.yaml.")
	StartCmd.Flags().StringVarP(&appClient, "client", "c", "", "Client name")

	RootCmd.AddCommand(StartCmd)
}