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

	"github.com/spf13/cobra"
)

// cCmd represents the c command
var cCmd = &cobra.Command{
	Use:   "c",
	Short: "Encrypt file",
	Long: `Use this command to encript your file, a master key will be asket two times, and you never forget it
For example:

crp c readme.me = encrypt readme.me file
crp c mydir     = encrypt content of dir mydir
`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteOriginalFile, err := cmd.Flags().GetBool("delete")
		if err != nil {
			panic(err)
		}
		crypt(args[0], deleteOriginalFile)
	},
}

func init() {

	rootCmd.AddCommand(cCmd)

	cCmd.PersistentFlags().BoolP("delete", "d", false, "Delete original file")
}

func crypt(f string, delete bool) {
	fileInfo, err := os.Stat(f)
	if err != nil {
		fmt.Println("file", f, "does not exist")
		os.Exit(1)
	}

	if fileInfo.Mode().IsRegular() {

		password, err := input.ReadPasswordAndConfirm()
		if err != nil {
			panic(err)
		}
		nf, err := encoder.EncryptFile(f, password, delete)
		fmt.Printf("crypt %s...", f)
		if err != nil {
			panic(err)
		}
		fmt.Printf("done -> %s\n", nf)
	}

	if fileInfo.IsDir() {
		fmt.Println("crypt entire dir", f)
		password, err := input.ReadPasswordAndConfirm()
		if err != nil {
			panic(err)
		}
		_ = filepath.WalkDir(f, func(path string, di fs.DirEntry, err error) error {
			s, _ := os.Stat(path)
			if s.Mode().IsRegular() {
				fmt.Printf("crypt %s...", path)
				nf, err := encoder.EncryptFile(path, password, delete)
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
