/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gorest-cli/cmd/transports/http/fiber"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := fiber.AddNewHandler("goserve", "ProductVariant"); err != nil {
			cmd.PrintErrln("failed to add new handler:", err)
			return
		}

		if err := fiber.RegisterNewHandler("ProductVariant"); err != nil {
			cmd.PrintErrln("failed to register new handler:", err)
			return
		}

		if err := fiber.AddNewRouter("goserve", "ProductVariant"); err != nil {
			cmd.PrintErrln("failed to add new router:", err)
			return
		}

		if err := fiber.RegisterNewRouter("ProductVariant"); err != nil {
			cmd.PrintErrln("failed to register new router:", err)
			return
		}

		cmd.Println("ProductVariant scaffolding created successfully")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
