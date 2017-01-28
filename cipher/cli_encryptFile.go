package main

import (
	"fmt"
	"github.com/freignat91/cipher/rsa"
	"github.com/spf13/cobra"
	"time"
)

var EncryptFileCmd = &cobra.Command{
	Use:   "encryptFile [sourcefilePath] [targetFilePath] [publicKeyFilePath]",
	Short: "encrypt file",
	Long:  `encrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.encryptFile(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(EncryptFileCmd)
}

func (m *cipherCLI) encryptFile(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage cipher encryptFile [sourcefilePath] [targetFilePath] [publicKeyFilePath]")
	}
	t0 := time.Now()
	if err := rsa.EncryptFile(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Printf("done time=%ds\n", time.Now().Sub(t0).Nanoseconds()/1000000000)
	return nil
}
