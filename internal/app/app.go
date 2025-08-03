package app

import (
	"os"
	"time"
	"vibe-check/internal/git"
	"vibe-check/internal/models"
	"vibe-check/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// refreshMsg indicates it's time to refresh the disabled state
type refreshMsg struct{}

// doRefresh returns a command that sends refresh message after delay
func doRefresh() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

// Menu options
var MenuOptions = []string{
	"Create Checkpoint",
	"Change Checkpoint",
	"Finalize and Push",
	"Exit",
}

var CheckpointCreationOptions = []string{
	"Create Checkpoint",
	"Create Checkpoint with Custom Note",
	"Back to Main Menu",
}

var FinalizeOptions = []string{
	"Finalize and Push (Auto Message)",
	"Finalize and Push with Custom Message",
	"Back to Main Menu",
}

// App wraps the models.AppModel and implements tea.Model
type App struct {
	models.AppModel
}

// InitialModel creates the initial application model
func InitialModel() App {
	app := App{
		AppModel: models.AppModel{
			CurrentState:      models.StateMenu,
			MenuChoices:       MenuOptions,
			CheckpointOptions: CheckpointCreationOptions,
			FinalizeOptions:   FinalizeOptions,
			DisabledMenuItems: make(map[int]bool),
			DisabledReasons:   make(map[int]string),
		},
	}
	// Update disabled items based on current state
	app.updateDisabledItems()
	// Move cursor to first enabled item
	app.moveToFirstEnabledItem()
	return app
}

// updateDisabledItems updates which menu items should be disabled
func (a *App) updateDisabledItems() {
	hasCheckpoints := git.HasCheckpoints()
	hasChanges := git.HasUncommittedChanges()
	
	for i, choice := range a.MenuChoices {
		switch choice {
		case "Create Checkpoint":
			if !hasChanges {
				a.DisabledMenuItems[i] = true
				a.DisabledReasons[i] = "(no changes)"
			} else {
				a.DisabledMenuItems[i] = false
				a.DisabledReasons[i] = ""
			}
		case "Change Checkpoint":
			a.DisabledMenuItems[i] = !hasCheckpoints
			a.DisabledReasons[i] = ""
		case "Finalize and Push":
			a.DisabledMenuItems[i] = !hasCheckpoints
			a.DisabledReasons[i] = ""
		default:
			a.DisabledMenuItems[i] = false
			a.DisabledReasons[i] = ""
		}
	}
}

// moveToFirstEnabledItem moves cursor to first enabled menu item
func (a *App) moveToFirstEnabledItem() {
	for i := 0; i < len(a.MenuChoices); i++ {
		if !a.DisabledMenuItems[i] {
			a.MenuCursor = i
			break
		}
	}
}

// Init initializes the Bubble Tea program
func (a App) Init() tea.Cmd {
	// Start periodic refresh for main menu
	return doRefresh()
}

// Update handles messages and updates the model
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	case resultMsg:
		return a.handleResult(msg)
	case checkpointsLoadedMsg:
		return a.handleCheckpointsLoaded(msg)
	case refreshMsg:
		return a.handleRefresh(msg)
	}
	return a, nil
}

// handleRefresh updates disabled state and schedules next refresh
func (a App) handleRefresh(msg refreshMsg) (tea.Model, tea.Cmd) {
	// Only refresh when on main menu to avoid unnecessary work
	if a.CurrentState == models.StateMenu {
		a.updateDisabledItems()
		// Schedule next refresh
		return a, doRefresh()
	}
	// If not on main menu, still schedule next refresh but don't update
	return a, doRefresh()
}

// View renders the current view
func (a App) View() string {
	switch a.CurrentState {
	case models.StateMenu:
		return ui.RenderMenu(a.AppModel)
	case models.StateCheckpointCreation:
		return ui.RenderCheckpointCreation(a.AppModel)
	case models.StateCheckpointNoteInput:
		return ui.RenderNoteInput(a.AppModel)
	case models.StateCheckpointSelection:
		return ui.RenderCheckpointSelection(a.AppModel)
	case models.StateFinalizeOptions:
		return ui.RenderFinalizeOptions(a.AppModel)
	case models.StateFinalizeMessageInput:
		return ui.RenderFinalizeMessageInput(a.AppModel)
	case models.StateExecuting:
		return ui.RenderLoading(a.AppModel)
	case models.StateResult:
		return ui.RenderResult(a.AppModel)
	}
	return "Unknown state"
}

// resultMsg represents a command result
type resultMsg struct {
	Content string
	IsError bool
}

// RunApp starts the Bubble Tea application
func RunApp() error {
	p := tea.NewProgram(
		InitialModel(),
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stderr),
	)
	_, err := p.Run()
	return err
}