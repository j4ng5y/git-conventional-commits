package cli

import (
	"log"

	"github.com/spf13/cobra"
)

// Run is the primary CLI function
func Run() {
	var (
		rootCMD = &cobra.Command{
			Use: "cc",
			Run: RootCMD,
		}

		initCMD = &cobra.Command{
			Use: "init",
			Run: InitCMD,
		}
	)

	rootCMD.AddCommand(initCMD)

	if err := rootCMD.Execute(); err != nil {
		log.Fatal(err)
	}
}
