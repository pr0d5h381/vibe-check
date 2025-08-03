package models

import "time"

// AppState represents the current state of the application
type AppState int

const (
	StateMenu AppState = iota
	StateCheckpointCreation
	StateCheckpointNoteInput
	StateCheckpointSelection
	StateFinalizeOptions
	StateFinalizeMessageInput
	StateExecuting
	StateResult
)

// Checkpoint represents a git checkpoint
type Checkpoint struct {
	Hash    string
	Message string
	Time    time.Time
}

// AppModel represents the main application model for Bubble Tea
type AppModel struct {
	// Current state
	CurrentState AppState

	// Menu navigation
	MenuCursor   int
	MenuChoices  []string
	DisabledMenuItems map[int]bool // tracks which menu items are disabled
	DisabledReasons   map[int]string // reasons why items are disabled (for display)

	// Checkpoint creation
	CheckpointOptions       []string
	CheckpointOptionsCursor int
	CustomNote              string

	// Checkpoint selection
	Checkpoints       []Checkpoint
	CheckpointCursor  int

	// Finalize options
	FinalizeOptions       []string
	FinalizeOptionsCursor int
	CustomCommitMessage   string

	// Execution state
	Loading     bool
	LoadingText string

	// Result display
	Result  string
	IsError bool
}

// Result represents an operation result
type Result struct {
	Content string
	IsError bool
}