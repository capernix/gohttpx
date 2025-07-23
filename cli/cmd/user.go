package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  "Create, list, show, and delete users via the GoHTTPx API",
}

var userCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new user",
	Long:    "Create a new user with the specified name",
	Example: `  gohttpx-cli user create --name "John Doe"`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if name == "" {
			PrintError("name is required")
			return
		}

		body := map[string]string{"name": name}
		resp, err := MakeRequest("POST", "/users", body)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var user User
		if err := ParseAPIResponse(resp, &user); err != nil {
			PrintError(err.Error())
			return
		}

		if jsonOutput {
			PrintUser(user, true)
		} else {
			PrintSuccess(fmt.Sprintf("User created (ID: %d)", user.ID))
		}
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  "List all users from the GoHTTPx API",
	Example: `  gohttpx-cli user list
  gohttpx-cli user list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		resp, err := MakeRequest("GET", "/users", nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var users []User
		if err := ParseAPIResponse(resp, &users); err != nil {
			PrintError(err.Error())
			return
		}

		PrintUsers(users, jsonOutput)
	},
}

var userShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific user",
	Long:  "Show details of a specific user by ID",
	Example: `  gohttpx-cli user show --id 1
  gohttpx-cli user show --id 1 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		idStr, _ := cmd.Flags().GetString("id")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		id, err := ValidateID(idStr)
		if err != nil {
			PrintError(err.Error())
			return
		}

		resp, err := MakeRequest("GET", "/users/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var user User
		if err := ParseAPIResponse(resp, &user); err != nil {
			PrintError(err.Error())
			return
		}

		PrintUser(user, jsonOutput)
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Long:  "Delete a user by ID",
	Example: `  gohttpx-cli user delete --id 1
  gohttpx-cli user delete --id 1 --confirm`,
	Run: func(cmd *cobra.Command, args []string) {
		idStr, _ := cmd.Flags().GetString("id")
		skipConfirm, _ := cmd.Flags().GetBool("confirm")

		id, err := ValidateID(idStr)
		if err != nil {
			PrintError(err.Error())
			return
		}

		resp, err := MakeRequest("GET", "/users/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var user User
		if err := ParseAPIResponse(resp, &user); err != nil {
			PrintError(err.Error())
			return
		}

		if !skipConfirm {
			prompt := fmt.Sprintf("Are you sure you want to delete user '%s' (ID: %d)?", user.Name, user.ID)
			if !ConfirmAction(prompt) {
				fmt.Println("Deletion cancelled")
				return
			}
		}

		resp, err = MakeRequest("DELETE", "/users/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		if resp.StatusCode == 200 || resp.StatusCode == 204 {
			PrintSuccess(fmt.Sprintf("User deleted (ID: %d)", id))
		} else {
			PrintError("Failed to delete user")
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userShowCmd)
	userCmd.AddCommand(userDeleteCmd)

	userCreateCmd.Flags().StringP("name", "n", "", "Name of the user (required)")
	userCreateCmd.MarkFlagRequired("name")

	userShowCmd.Flags().StringP("id", "i", "", "ID of the user (required)")
	userShowCmd.MarkFlagRequired("id")

	userDeleteCmd.Flags().StringP("id", "i", "", "ID of the user (required)")
	userDeleteCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
	userDeleteCmd.MarkFlagRequired("id")
}
