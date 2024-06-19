package cmd

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/paulscherrerinstitute/scicat/datasetUtils"
	"github.com/spf13/cobra"
)

var datasetPublishDataRetrieveCmd = &cobra.Command{
	Use:   "datasetPublishDataRetrieve [options]",
	Short: "Create a job to retrieve all datasets of a given PublishedData item",
	Long:  `Create a job to retrieve all datasets of a given PublishedData item.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// consts & vars
		const PROD_API_SERVER string = "https://dacat.psi.ch/api/v3"
		const TEST_API_SERVER string = "https://dacat-qa.psi.ch/api/v3"
		const DEV_API_SERVER string = "https://dacat-development.psi.ch/api/v3"

		var APIServer string = PROD_API_SERVER

		var client = &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}},
			Timeout:   10 * time.Second}

		// retrieve params
		retrieveFlag, _ := cmd.Flags().GetBool("retrieve")
		publishedDataId, _ := cmd.Flags().GetString("publisheddata")
		userpass, _ := cmd.Flags().GetString("user")
		token, _ := cmd.Flags().GetString("token")
		testenvFlag, _ := cmd.Flags().GetBool("testenv")
		devenvFlag, _ := cmd.Flags().GetBool("devenv")
		showVersion, _ := cmd.Flags().GetBool("version")

		// execute command
		if showVersion {
			fmt.Printf("%s\n", VERSION)
			return
		}

		var env string
		if testenvFlag {
			APIServer = TEST_API_SERVER
			env = "test"
		} else if devenvFlag {
			APIServer = DEV_API_SERVER
			env = "dev"
		} else {
			APIServer = PROD_API_SERVER
			env = "production"
		}

		color.Set(color.FgGreen)
		log.Printf("You are about to trigger a retrieve job for publish dataset(s) from the === %s === retrieve server...", env)
		color.Unset()

		if !retrieveFlag {
			color.Set(color.FgRed)
			log.Printf("Note: you run in 'dry' mode to simply check which data would be retrieved.\n")
			log.Printf("Use the -retrieve flag to actually retrieve the datasets.\n")
			color.Unset()
		}

		if publishedDataId == "" { /* && *datasetId == "" && *ownerGroup == "" */
			fmt.Println("\n\nTool to retrieve datasets to the intermediate cache server of the tape archive")
			fmt.Printf("Run script without arguments, but specify options:\n\n")
			fmt.Printf("datasetPublishDataRetrieve [options] \n\n")
			fmt.Printf("Use -publisheddata option to define the datasets which should be published.\n\n")
			fmt.Printf("For example:\n")
			fmt.Printf("./datasetPublishDataRetrieve -user archiveManager:password -publisheddata 10.16907/05a50450-767f-421d-9832-342b57c201\n\n")
			fmt.Printf("The script should be run as archiveManager\n\n")
			flag.PrintDefaults()
			return
		}

		auth := &datasetUtils.RealAuthenticator{}
		user, _ := datasetUtils.Authenticate(auth, client, APIServer, &token, &userpass)

		datasetList, _, _ := datasetUtils.GetDatasetsOfPublication(client, APIServer, publishedDataId)

		// get sourceFolder and other dataset related info for all Datasets and print them
		datasetUtils.GetDatasetDetailsPublished(client, APIServer, datasetList)

		if !retrieveFlag {
			color.Set(color.FgRed)
			log.Printf("\n\nNote: you run in 'dry' mode to simply check what would happen.")
			log.Printf("Use the -retrieve flag to actually retrieve data from tape.\n")
			color.Unset()
		} else {
			// create retrieve Job
			jobId, err := datasetUtils.CreateRetrieveJob(client, APIServer, user, datasetList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(jobId)
		}
	},
}

func init() {
	rootCmd.AddCommand(datasetPublishDataRetrieveCmd)

	datasetPublishDataRetrieveCmd.Flags().Bool("retrieve", false, "Defines if this command is meant to actually retrieve data (default: retrieve actions are only displayed)")
	datasetPublishDataRetrieveCmd.Flags().String("publisheddata", "", "Defines to publish data from a given publishedData document ID")
	datasetPublishDataRetrieveCmd.Flags().String("user", "", "Defines optional username:password string")
	datasetPublishDataRetrieveCmd.Flags().String("token", "", "Defines optional API token instead of username:password")
	datasetPublishDataRetrieveCmd.Flags().Bool("testenv", false, "Use test environment (qa) (default is to use production system)")
	datasetPublishDataRetrieveCmd.Flags().Bool("devenv", false, "Use development environment (default is to use production system)")
	datasetPublishDataRetrieveCmd.Flags().Bool("version", false, "Show version number and exit")

	datasetPublishDataRetrieveCmd.MarkFlagsMutuallyExclusive("testenv", "devenv")
}
