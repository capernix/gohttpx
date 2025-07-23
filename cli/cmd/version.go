package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the version number of the GoHTTPx CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoHTTPx CLI v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
