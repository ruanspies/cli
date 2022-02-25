package cmd

import (
	"context"
	"embed"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pbOperations "go.protobuf.alis.alis.exchange/alis/os/resources/operations/v1"
	pbProducts "go.protobuf.alis.alis.exchange/alis/os/resources/products/v1"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	organisationID       string
	productID            string
	neuronID             string
	releaseType          string
	debugFlag            bool
	cfgFile              string
	homeDir              string
	asyncFlag            bool
	alisProductsClient   pbProducts.ServiceClient
	alisOperationsClient pbOperations.ServiceClient
	TemplateFs           embed.FS
	ptermTip             pterm.PrefixPrinter
	ptermInput           pterm.PrefixPrinter
)

const VERSION = "3.3.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alis",
	Short: pterm.Green("alis_ Technologies LLC - Command Line Interface"),
	Long: pterm.Green("The alis CLI manages authentication, local configuration, developer workflow, \n" +
		"and interactions with the alis_ os resources"),
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugFlag {
			pterm.EnableDebugMessages()
		}
	},
}

//var testCmd = &cobra.Command{
//	Use:   "test",
//	Short: pterm.Green("a test short command"),
//	Run: func(cmd *cobra.Command, args []string) {
//
//		//generateFileDescriptorSet("organisations/alis/products/in/neurons/resources-instruments-v2")
//		generateFileDescriptorSet("organisations/alis/products/fe/neurons/resources-store-v1")
//
//		return
//	},
//}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Printf("\033[32m%s\033[0m", err)
	}
	// Initialise alis Products client
	var connProducts *grpc.ClientConn
	connProducts, err = NewServerConnection(context.Background(), "resources-products-v1-ntaj7kcaca-ew.a.run.app")
	if err != nil {
		log.Fatalf("alis.NewServerConnection: %s", err)
	}
	alisProductsClient = pbProducts.NewServiceClient(connProducts)

	// Initialise alis Services client
	var connOperations *grpc.ClientConn
	connOperations, err = NewServerConnection(context.Background(), "resources-operations-v1-ntaj7kcaca-ew.a.run.app")
	if err != nil {
		log.Fatalf("alis.NewServerConnection: %s", err)
	}
	alisOperationsClient = pbOperations.NewServiceClient(connOperations)

	cobra.OnInitialize(initConfig)
	rootCmd.Version = VERSION
	//rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, pterm.Green("Run the commands in DEBUG mode."))
	rootCmd.PersistentFlags().BoolVarP(&asyncFlag, "async", "a", false, pterm.Green("Return immediately, without waiting for the operation in progress to complete.\nOnly relevant if the command involves a long-running operation"))
	//rootCmd.PersistentFlags().BoolVarP(&localFlag, "local", "l", false, pterm.Green("If possible run command in local mode.  For example, use with neuron '--genproto-go' tag to generate protos only locally."))
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	// Define own commandline message type to use for tips.
	ptermTip = pterm.PrefixPrinter{
		Prefix: pterm.Prefix{
			Text:  " TIP ",
			Style: pterm.NewStyle(pterm.FgLightYellow, pterm.BgDarkGray),
		},
		Scope:          pterm.Scope{},
		MessageStyle:   pterm.NewStyle(pterm.FgLightYellow),
		Fatal:          false,
		ShowLineNumber: false,
		Debugger:       false,
	}

	ptermInput = pterm.PrefixPrinter{
		Prefix: pterm.Prefix{
			Text:  " INPUT ",
			Style: pterm.NewStyle(pterm.FgLightWhite, pterm.BgLightRed),
		},
		Scope:          pterm.Scope{},
		MessageStyle:   pterm.NewStyle(pterm.FgLightWhite),
		Fatal:          false,
		ShowLineNumber: false,
		Debugger:       false,
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".alis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".alis")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
