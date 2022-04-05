package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func main() {
	var cmdPull = &cobra.Command{
		Use:   "pull [OPTIONS] NAME[:TAG|@DIGEST]",
		Short: "Pull an image or a repository fron a registry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("pull: " + strings.Join(args, " "))
		},
	}

	var rootCmd = &cobra.Command{
		Use: "docker",
	}
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print Usage")
	rootCmd.AddCommand(cmdPull)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}
