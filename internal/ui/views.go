package ui

import (
	"fmt"
	"strings"
	"time"
	"vibe-check/internal/git"
	"vibe-check/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// RenderNoteInput renders the note input view
func RenderNoteInput(m models.AppModel) string {
	var s strings.Builder

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoStyle.Render("Custom Note"),
		"  ",
		AppCaption.Render("Enter a note for this checkpoint"),
	)

	// Input field
	noteDisplay := m.CustomNote
	if len(noteDisplay) == 0 {
		noteDisplay = AppCaption.Render("Type your note here...")
	}
	
	// Add cursor
	noteDisplay += MenuPointer.Render("│")
	
	// Character counter
	counter := fmt.Sprintf("(%d/50)", len(m.CustomNote))
	counterStyle := AppCaption
	if len(m.CustomNote) > 40 {
		counterStyle = ErrorStyle
	}
	
	inputSection := fmt.Sprintf("%s\n%s", 
		MenuItem.Render(noteDisplay),
		counterStyle.Render(counter),
	)
	
	footer := HelpStyle.Render("Type to add text • Backspace to delete • Enter to create • Esc to cancel")
	dividerLine := Hairline.Render(strings.Repeat("─", 50))
	
	body := inputSection + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}

func RenderFinalizeMessageInput(m models.AppModel) string {
	var s strings.Builder

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoStyle.Render("Custom Commit Message"),
		"  ",
		AppCaption.Render("Enter commit message for finalize"),
	)

	// Input field
	messageDisplay := m.CustomCommitMessage
	if len(messageDisplay) == 0 {
		messageDisplay = AppCaption.Render("Type your commit message here...")
	}
	
	// Add cursor
	messageDisplay += MenuPointer.Render("│")
	
	// Character counter
	counter := fmt.Sprintf("(%d/100)", len(m.CustomCommitMessage))
	counterStyle := AppCaption
	if len(m.CustomCommitMessage) > 80 {
		counterStyle = ErrorStyle
	}
	
	inputSection := fmt.Sprintf("%s\n%s", 
		MenuItem.Render(messageDisplay),
		counterStyle.Render(counter),
	)
	
	footer := HelpStyle.Render("Type to add text • Backspace to delete • Enter to finalize • Esc to cancel")
	dividerLine := Hairline.Render(strings.Repeat("─", 60))
	
	body := inputSection + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}

// RenderCheckpointSelection renders the checkpoint selection view
func RenderCheckpointSelection(m models.AppModel) string {
	var s strings.Builder

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoStyle.Render("Checkpoints"),
		"  ",
		AppCaption.Render("Select a checkpoint to switch"),
	)

	if len(m.Checkpoints) == 0 {
		body := AppCaption.Render("No checkpoints found\nTip: Create your first checkpoint")
		footer := HelpStyle.Render("Esc back to main menu")
		content := body + "\n" + Hairline.Render(strings.Repeat("─", 30)) + "\n" + footer
		
		s.WriteString(CardAlt.Render(title) + "\n")
		s.WriteString(Card.Render(content))
		return s.String()
	}

	var list strings.Builder
	
	// Get current commit hash for comparison
	currentCommit, _ := git.GetCurrentCommit()
	
	for i, cp := range m.Checkpoints {
		prefix := "  "
		lineStyle := MenuItem
		
		// Check if this is the current commit
		isCurrentCheckpoint := cp.Hash == currentCommit
		
		// Handle cursor selection and styling
		if i == m.CheckpointCursor {
			prefix = MenuPointer.Render("› ")
			lineStyle = MenuItemActive
		}
		
		// If this is the current checkpoint, render with green style regardless of cursor
		if isCurrentCheckpoint {
			line := fmt.Sprintf("%s%s", prefix, CurrentCheckpointStyle.Render(fmt.Sprintf("[%s] — %s", cp.Hash, cp.Message)))
			list.WriteString(line)
		} else {
			line := fmt.Sprintf("%s[%s] — %s", prefix, cp.Hash, cp.Message)
			list.WriteString(lineStyle.Render(line))
		}
		list.WriteString("\n")
	}
	
	footer := HelpStyle.Render("↑/↓ navigate • Enter switch • Esc back")
	dividerLine := Hairline.Render(strings.Repeat("─", 40))
	
	body := strings.TrimRight(list.String(), "\n") + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}

// RenderLoading renders the loading view
func RenderLoading(m models.AppModel) string {
	// Smooth spinner animation
	frames := []string{"⠁", "⠂", "⠄", "⠂"}
	frame := int(time.Now().UnixNano()/140000000) % len(frames)
	
	line := LoadingTextStyle.Render(frames[frame] + " " + m.LoadingText)
	return Card.Render(line)
}

// RenderResult renders the result/message view
func RenderResult(m models.AppModel) string {
	var msg string
	switch {
	case m.IsError:
		msg = ErrorStyle.Render(m.Result)
	case strings.Contains(m.Result, "Checkpoint created") || 
		 strings.Contains(m.Result, "Switched to checkpoint") ||
		 strings.Contains(m.Result, "Successfully"):
		msg = SuccessStyle.Render(m.Result)
	default:
		msg = InfoStyle.Render(m.Result)
	}
	
	footer := HelpStyle.Render("Press any key to continue…")
	dividerLine := Hairline.Render(strings.Repeat("─", 30))
	
	content := msg + "\n" + dividerLine + "\n" + footer
	return Card.Render(content)
}