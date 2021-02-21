package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RootCMD(ccmd *cobra.Command, args []string) {
	fmt.Println("CC")
}
