package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	Long:  "Create, list, show, and delete notes via the GoHTTPx API",
}

var noteCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new note",
	Long:    "Create a new note with the specified title and content",
	Example: `  gohttpx-cli note create --title "Meeting" --content "Discuss project roadmap"`,
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		if title == "" {
			PrintError("title is required")
			return
		}

		if content == "" {
			PrintError("content is required")
			return
		}

		body := map[string]string{
			"title":   title,
			"content": content,
		}

		resp, err := MakeRequest("POST", "/notes", body)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var note Note
		if err := ParseAPIResponse(resp, &note); err != nil {
			PrintError(err.Error())
			return
		}

		if jsonOutput {
			PrintNote(note, true)
		} else {
			PrintSuccess(fmt.Sprintf("Note created (ID: %d)", note.ID))
		}
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  "List all notes from the GoHTTPx API",
	Example: `  gohttpx-cli note list
  gohttpx-cli note list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")

		resp, err := MakeRequest("GET", "/notes", nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var notes []Note
		if err := ParseAPIResponse(resp, &notes); err != nil {
			PrintError(err.Error())
			return
		}

		PrintNotes(notes, jsonOutput)
	},
}

var noteShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific note",
	Long:  "Show details of a specific note by ID",
	Example: `  gohttpx-cli note show --id 1
  gohttpx-cli note show --id 1 --json`,
	Run: func(cmd *cobra.Command, args []string) {
		idStr, _ := cmd.Flags().GetString("id")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		id, err := ValidateID(idStr)
		if err != nil {
			PrintError(err.Error())
			return
		}

		resp, err := MakeRequest("GET", "/notes/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var note Note
		if err := ParseAPIResponse(resp, &note); err != nil {
			PrintError(err.Error())
			return
		}

		PrintNote(note, jsonOutput)
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a note",
	Long:  "Delete a note by ID",
	Example: `  gohttpx-cli note delete --id 1
  gohttpx-cli note delete --id 1 --confirm`,
	Run: func(cmd *cobra.Command, args []string) {
		idStr, _ := cmd.Flags().GetString("id")
		skipConfirm, _ := cmd.Flags().GetBool("confirm")

		id, err := ValidateID(idStr)
		if err != nil {
			PrintError(err.Error())
			return
		}

		resp, err := MakeRequest("GET", "/notes/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		var note Note
		if err := ParseAPIResponse(resp, &note); err != nil {
			PrintError(err.Error())
			return
		}

		if !skipConfirm {
			prompt := fmt.Sprintf("Are you sure you want to delete note '%s' (ID: %d)?", note.Title, note.ID)
			if !ConfirmAction(prompt) {
				fmt.Println("Deletion cancelled")
				return
			}
		}

		resp, err = MakeRequest("DELETE", "/notes/"+strconv.Itoa(id), nil)
		if err != nil {
			PrintError(err.Error())
			return
		}

		if resp.StatusCode == 200 || resp.StatusCode == 204 {
			PrintSuccess(fmt.Sprintf("Note deleted (ID: %d)", id))
		} else {
			PrintError("Failed to delete note")
		}
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)

	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteShowCmd)
	noteCmd.AddCommand(noteDeleteCmd)

	noteCreateCmd.Flags().StringP("title", "t", "", "Title of the note (required)")
	noteCreateCmd.Flags().StringP("content", "c", "", "Content of the note (required)")
	noteCreateCmd.MarkFlagRequired("title")
	noteCreateCmd.MarkFlagRequired("content")

	noteShowCmd.Flags().StringP("id", "i", "", "ID of the note (required)")
	noteShowCmd.MarkFlagRequired("id")

	noteDeleteCmd.Flags().StringP("id", "i", "", "ID of the note (required)")
	noteDeleteCmd.Flags().BoolP("confirm", "y", false, "Skip confirmation prompt")
	noteDeleteCmd.MarkFlagRequired("id")
}
