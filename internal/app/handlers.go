package app

import (
	"fmt"
	"strings"
	"vibe-check/internal/git"
	"vibe-check/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyPress processes keyboard input based on current state
func (a App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.CurrentState {
	case models.StateMenu:
		return a.handleMenuKeys(msg)
	case models.StateCheckpointCreation:
		return a.handleCheckpointCreationKeys(msg)
	case models.StateCheckpointNoteInput:
		return a.handleNoteInputKeys(msg)
	case models.StateCheckpointSelection:
		return a.handleCheckpointSelectionKeys(msg)
	case models.StateFinalizeOptions:
		return a.handleFinalizeOptionsKeys(msg)
	case models.StateFinalizeMessageInput:
		return a.handleFinalizeMessageInputKeys(msg)
	case models.StateResult:
		return a.handleResultKeys(msg)
	}
	return a, nil
}

// handleMenuKeys processes keys in the main menu
func (a App) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "up", "k":
		a.moveCursorUp()
	case "down", "j":
		a.moveCursorDown()
	case "enter", " ":
		return a.executeMenuAction()
	}
	return a, nil
}

// moveCursorUp moves cursor up, skipping disabled items
func (a *App) moveCursorUp() {
	for {
		if a.MenuCursor > 0 {
			a.MenuCursor--
			// If current item is not disabled, we're done
			if !a.DisabledMenuItems[a.MenuCursor] {
				break
			}
		} else {
			// Reached top, wrap to bottom but find first enabled item from bottom
			a.MenuCursor = len(a.MenuChoices) - 1
			for a.MenuCursor >= 0 && a.DisabledMenuItems[a.MenuCursor] {
				a.MenuCursor--
			}
			break
		}
	}
}

// moveCursorDown moves cursor down, skipping disabled items
func (a *App) moveCursorDown() {
	for {
		if a.MenuCursor < len(a.MenuChoices)-1 {
			a.MenuCursor++
			// If current item is not disabled, we're done
			if !a.DisabledMenuItems[a.MenuCursor] {
				break
			}
		} else {
			// Reached bottom, wrap to top but find first enabled item from top
			a.MenuCursor = 0
			for a.MenuCursor < len(a.MenuChoices) && a.DisabledMenuItems[a.MenuCursor] {
				a.MenuCursor++
			}
			break
		}
	}
}

// executeMenuAction executes the selected menu action
func (a App) executeMenuAction() (tea.Model, tea.Cmd) {
	// Check if current item is disabled
	if a.DisabledMenuItems[a.MenuCursor] {
		return a, nil // Do nothing if disabled
	}
	
	selected := a.MenuChoices[a.MenuCursor]

	switch {
	case strings.HasPrefix(selected, "Create"):
		a.CurrentState = models.StateCheckpointCreation
		a.CheckpointOptionsCursor = 0
		return a, nil
	case strings.HasPrefix(selected, "Change"):
		a.CurrentState = models.StateExecuting
		a.Loading = true
		a.LoadingText = "Loading checkpoints..."
		return a.loadCheckpoints()
	case strings.HasPrefix(selected, "Finalize"):
		a.CurrentState = models.StateFinalizeOptions
		a.FinalizeOptionsCursor = 0
		return a, nil
	case strings.HasPrefix(selected, "Exit"):
		return a, tea.Quit
	}
	return a, nil
}

// handleCheckpointCreationKeys processes keys in checkpoint creation menu
func (a App) handleCheckpointCreationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		a.CurrentState = models.StateMenu
		a.updateDisabledItems()
		a.moveToFirstEnabledItem()
		return a, nil
	case "up", "k":
		if a.CheckpointOptionsCursor > 0 {
			a.CheckpointOptionsCursor--
		}
	case "down", "j":
		if a.CheckpointOptionsCursor < len(a.CheckpointOptions)-1 {
			a.CheckpointOptionsCursor++
		}
	case "enter", " ":
		return a.executeCheckpointAction()
	}
	return a, nil
}

// executeCheckpointAction executes the selected checkpoint action
func (a App) executeCheckpointAction() (tea.Model, tea.Cmd) {
	selected := a.CheckpointOptions[a.CheckpointOptionsCursor]

	switch {
	case strings.HasPrefix(selected, "Create Checkpoint with"):
		a.CurrentState = models.StateCheckpointNoteInput
		a.CustomNote = ""
		return a, nil
	case strings.HasPrefix(selected, "Create Checkpoint"):
		return a.createCheckpoint("")
	case strings.HasPrefix(selected, "Back"):
		a.CurrentState = models.StateMenu
		a.updateDisabledItems()
		a.moveToFirstEnabledItem()
		return a, nil
	}
	return a, nil
}

// handleNoteInputKeys processes keys in note input state
func (a App) handleNoteInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		a.CurrentState = models.StateCheckpointCreation
		a.CustomNote = ""
		return a, nil
	case "enter":
		return a.createCheckpoint(a.CustomNote)
	case "backspace":
		if len(a.CustomNote) > 0 {
			a.CustomNote = a.CustomNote[:len(a.CustomNote)-1]
		}
	default:
		// Add character to note
		if len(msg.String()) == 1 && len(a.CustomNote) < 50 {
			a.CustomNote += msg.String()
		}
	}
	return a, nil
}

// handleCheckpointSelectionKeys processes keys in checkpoint selection
func (a App) handleCheckpointSelectionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		a.CurrentState = models.StateMenu
		a.updateDisabledItems()
		a.moveToFirstEnabledItem()
		return a, nil
	case "up", "k":
		if a.CheckpointCursor > 0 {
			a.CheckpointCursor--
		}
	case "down", "j":
		if a.CheckpointCursor < len(a.Checkpoints)-1 {
			a.CheckpointCursor++
		}
	case "enter", " ":
		if len(a.Checkpoints) > 0 {
			selected := a.Checkpoints[a.CheckpointCursor]
			return a.switchToCheckpoint(selected.Hash)
		}
	}
	return a, nil
}

// handleResultKeys processes keys in result display
func (a App) handleResultKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	a.CurrentState = models.StateMenu
	a.Result = ""
	a.updateDisabledItems() // Update disabled items when returning to menu
	a.moveToFirstEnabledItem()
	return a, nil
}

// handleResult processes result messages
func (a App) handleResult(msg resultMsg) (tea.Model, tea.Cmd) {
	a.Loading = false
	a.CurrentState = models.StateResult
	a.Result = msg.Content
	a.IsError = msg.IsError
	return a, nil
}

// showPlaceholder shows a placeholder result message
func (a App) showPlaceholder(message string) (tea.Model, tea.Cmd) {
	a.CurrentState = models.StateResult
	a.Result = message
	a.IsError = false
	return a, nil
}

// createCheckpoint creates a git checkpoint
func (a App) createCheckpoint(customNote string) (tea.Model, tea.Cmd) {
	return a, func() tea.Msg {
		err := git.CreateCheckpoint(customNote)
		if err != nil {
			return resultMsg{
				Content: "Error creating checkpoint: " + err.Error(),
				IsError: true,
			}
		}
		
		message := "Checkpoint created successfully"
		if customNote != "" {
			message += " with note: \"" + customNote + "\""
		}
		
		return resultMsg{
			Content: message,
			IsError: false,
		}
	}
}

// loadCheckpoints loads checkpoints for selection
func (a App) loadCheckpoints() (tea.Model, tea.Cmd) {
	return a, func() tea.Msg {
		checkpoints, err := git.GetCheckpointsFromReflog()
		if err != nil {
			return resultMsg{
				Content: "Error loading checkpoints: " + err.Error(),
				IsError: true,
			}
		}
		
		return checkpointsLoadedMsg{
			Checkpoints: checkpoints,
		}
	}
}

// switchToCheckpoint switches to a specific checkpoint
func (a App) switchToCheckpoint(hash string) (tea.Model, tea.Cmd) {
	return a, func() tea.Msg {
		err := git.SwitchToCheckpoint(hash)
		if err != nil {
			return resultMsg{
				Content: "Error switching to checkpoint: " + err.Error(),
				IsError: true,
			}
		}
		
		return resultMsg{
			Content: "Switched to checkpoint: " + hash,
			IsError: false,
		}
	}
}

// checkpointsLoadedMsg represents loaded checkpoints
type checkpointsLoadedMsg struct {
	Checkpoints []models.Checkpoint
}

// handleCheckpointsLoaded handles loaded checkpoints message
func (a App) handleCheckpointsLoaded(msg checkpointsLoadedMsg) (tea.Model, tea.Cmd) {
	a.Loading = false
	a.CurrentState = models.StateCheckpointSelection
	a.Checkpoints = msg.Checkpoints
	a.CheckpointCursor = 0
	return a, nil
}

// finalizeAndPush finalizes checkpoints and pushes to remote
func (a App) finalizeAndPush() (tea.Model, tea.Cmd) {
	return a.finalizeAndPushWithMessage("")
}

// finalizeAndPushWithMessage finalizes checkpoints and pushes to remote with custom message
func (a App) finalizeAndPushWithMessage(customMessage string) (tea.Model, tea.Cmd) {
	a.CurrentState = models.StateExecuting
	a.Loading = true
	a.LoadingText = "Finalizing and pushing..."
	
	return a, func() tea.Msg {
		// Proceed with finalize and push
		err := git.FinalizeAndPushWithMessage(customMessage)
		if err != nil {
			return resultMsg{
				Content: "Error during finalize and push: " + err.Error(),
				IsError: true,
			}
		}

		var successMessage string
		if customMessage != "" {
			successMessage = fmt.Sprintf("Successfully finalized and pushed with message: \"%s\"", customMessage)
		} else {
			successMessage = "Successfully finalized and pushed checkpoints to remote!"
		}

		return resultMsg{
			Content: successMessage,
			IsError: false,
		}
	}
}


// handleFinalizeOptionsKeys processes keys in finalize options menu
func (a App) handleFinalizeOptionsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		a.CurrentState = models.StateMenu
		a.updateDisabledItems()
		a.moveToFirstEnabledItem()
		return a, nil
	case "up", "k":
		if a.FinalizeOptionsCursor > 0 {
			a.FinalizeOptionsCursor--
		}
	case "down", "j":
		if a.FinalizeOptionsCursor < len(a.FinalizeOptions)-1 {
			a.FinalizeOptionsCursor++
		}
	case "enter", " ":
		return a.executeFinalizeAction()
	}
	return a, nil
}

// executeFinalizeAction executes the selected finalize action
func (a App) executeFinalizeAction() (tea.Model, tea.Cmd) {
	selected := a.FinalizeOptions[a.FinalizeOptionsCursor]

	switch {
	case strings.HasPrefix(selected, "Finalize and Push with Custom"):
		a.CurrentState = models.StateFinalizeMessageInput
		a.CustomCommitMessage = ""
		return a, nil
	case strings.HasPrefix(selected, "Finalize and Push (Auto"):
		return a.finalizeAndPushWithMessage("")
	case strings.HasPrefix(selected, "Back"):
		a.CurrentState = models.StateMenu
		a.updateDisabledItems()
		a.moveToFirstEnabledItem()
		return a, nil
	}
	return a, nil
}

// handleFinalizeMessageInputKeys processes keys in finalize message input state
func (a App) handleFinalizeMessageInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		a.CurrentState = models.StateFinalizeOptions
		return a, nil
	case "enter":
		if len(a.CustomCommitMessage) > 0 {
			return a.finalizeAndPushWithMessage(a.CustomCommitMessage)
		}
		return a, nil
	case "backspace":
		if len(a.CustomCommitMessage) > 0 {
			a.CustomCommitMessage = a.CustomCommitMessage[:len(a.CustomCommitMessage)-1]
		}
	default:
		// Add character to message if it's printable and under limit
		if len(msg.String()) == 1 && len(a.CustomCommitMessage) < 100 {
			a.CustomCommitMessage += msg.String()
		}
	}
	return a, nil
}