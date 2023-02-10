/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"mi0772/crp/encoder"
	"mi0772/crp/input"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// dCmd represents the d command
var dCmd = &cobra.Command{
	Use:   "d",
	Short: "Decrypt file",
	Long: `Use this command to decrypt your file

	For example:

	crp d readme.me = decrypt readme.me file
	crp d mydir     = decrypt content of dir mydir`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteFile, err := cmd.Flags().GetBool("delete")
		if err != nil {
			panic(err)
		}
		decrypt(args[0], deleteFile)
	},
}

func decrypt(f string, delete bool) {
	fileInfo, err := os.Stat(f)
	if err != nil {
		fmt.Println("file", f, "does not exist")
		os.Exit(1)
	}

	password, err := input.ReadPassword()
	if err != nil {
		panic(err)
	}
	if fileInfo.Mode().IsRegular() {
		nf, err := encoder.DecryptFile(f, password, delete)
		fmt.Printf("\ndecrypt %s...", f)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("done -> %s\n", nf)
	}

	if fileInfo.IsDir() {
		fmt.Println("decrypt entire dir", f)
		fmt.Println()
		_ = filepath.WalkDir(f, func(path string, di fs.DirEntry, err error) error {
			s, _ := os.Stat(path)
			if s.Mode().IsRegular() && strings.Contains(path, ".encrypted") {
				fmt.Printf("decrypt %s...", path)
				nf, err := encoder.DecryptFile(path, password, delete)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("done -> %s\n", nf)
			}

			return nil
		})
	}
}

func init() {
	rootCmd.AddCommand(dCmd)
	dCmd.PersistentFlags().BoolP("delete", "d", false, "Delete encryted file")
}
