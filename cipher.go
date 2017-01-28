package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use:   `cipher [OPTIONS] COMMAND [arg...]`,
		Short: "cipher commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}
)

func cli() {
	RootCmd.PersistentFlags().BoolVarP(&agridCli.verbose, "verbose", "v", false, `Verbose output`)
	cobra.OnInitialize(func() {
		if err := agridCli.init(); err != nil {
			fmt.Printf("Init error: %v\n", err)
			os.Exit(1)
		}
	})

	// versionCmd represents the agrid version
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version number of cipher",
		Long:  `Display the version number of cipher`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cipher version: %s, build: %s)\n", Version, Build)
		},
	}
	RootCmd.AddCommand(versionCmd)

	// infoCmd represents the agrid information
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display agrid version and server information",
		Long:  `Display agrid version and server information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("agrid version: %s, build: %s)\n", Version, Build)
			fmt.Printf("Server: %s\n", config.serverAddress)
		},
	}
	RootCmd.AddCommand(infoCmd)

	//Execute commad
	cmd, _, err := RootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error during: %s: %v\n", cmd.Name(), err)
		os.Exit(1)
	}

	os.Exit(0)
}
