package main

import (
	"fmt"
	"os"
	"vibe-check/internal/app"
	"vibe-check/internal/git"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vibe-check",
	Short: "Simple Git workflow made easy - no more complex commands!",
	Long: `Vibe Check - Simple Git workflow made easy

Tired of complicated Git commands? Vibe Check gives you a clean, intuitive interface 
for creating checkpoints, switching between versions, and cleaning up commit history.
Perfect for AI-assisted development and experimental coding.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.RunApp(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var createCmd = &cobra.Command{
	Use:   "create [note]",
	Short: "Create a new checkpoint",
	Long:  "Create a new checkpoint with optional custom note",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var note string
		if len(args) > 0 {
			note = args[0]
		}
		
		err := git.CreateCheckpoint(note)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		if note != "" {
			fmt.Printf("✅ Checkpoint created with note: %s\n", note)
		} else {
			fmt.Printf("✅ Checkpoint created\n")
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all checkpoints",
	Long:  "Display all available checkpoints with their hashes and messages",
	Run: func(cmd *cobra.Command, args []string) {
		checkpoints, err := git.GetCheckpointsFromReflog()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		if len(checkpoints) == 0 {
			fmt.Println("No checkpoints found")
			return
		}
		
		// Get current commit for highlighting
		currentCommit, _ := git.GetCurrentCommit()
		
		fmt.Println("📋 Checkpoints:")
		for _, cp := range checkpoints {
			marker := "  "
			if cp.Hash == currentCommit {
				marker = "* " // Current checkpoint
			}
			fmt.Printf("%s[%s] %s\n", marker, cp.Hash, cp.Message)
		}
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <hash>",
	Short: "Switch to a specific checkpoint",
	Long:  "Switch to a checkpoint using its commit hash",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]
		
		err := git.SwitchToCheckpoint(hash)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Switched to checkpoint %s\n", hash)
	},
}

var finalizeCmd = &cobra.Command{
	Use:   "finalize [message]",
	Short: "Finalize and push checkpoints",
	Long:  "Squash consecutive checkpoints into clean commits and push to remote",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var message string
		if len(args) > 0 {
			message = args[0]
		}
		
		err := git.FinalizeAndPushWithMessage(message)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("✅ Successfully finalized and pushed!\n")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd) 
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(finalizeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}