package main

import (
	"fmt"
	"github.com/freignat91/cipher/rsa"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var CreateKeysCmd = &cobra.Command{
	Use:   "createKeys [keysPath/name]",
	Short: "Create public and private RSA keys",
	Long:  `Create public and private RSA keys`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.createRSAKeys(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(CreateKeysCmd)
	CreateKeysCmd.Flags().String("size", "8192", `RSA Keys size (bit) should be a multiple of 64`)
}

func (m *cipherCLI) createRSAKeys(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("need key file path as argument. usage createKeys [keyFilePath]")
	}
	keyBitSize, err := strconv.Atoi(cmd.Flag("size").Value.String())
	if err != nil {
		return fmt.Errorf("option --size is not a number")
	}
	path := args[0]
	t0 := time.Now()
	publicKey, privateKey, err := rsa.CreateRSAKey(keyBitSize, m.verbose, m.debug)
	if err != nil {
		return err
	}
	if m.verbose {
		fmt.Printf("Compute time=%ds\n", time.Now().Sub(t0).Nanoseconds()/1000000000)
		fmt.Printf("Public key: %s\n", publicKey.ToHexa())
		fmt.Printf("Private key: %s\n", privateKey.ToHexa())
	}
	if err := rsa.SaveKeys(path, publicKey, privateKey); err != nil {
		return err
	}
	return nil
}
