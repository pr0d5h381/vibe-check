package git

import (
	"fmt"
	"strings"
	"vibe-check/internal/models"
)

// CreateCheckpoint creates a new git checkpoint
func CreateCheckpoint(customNote string) error {
	if !IsRepo() {
		return fmt.Errorf("not in a Git repository")
	}

	// Check if there are changes to commit
	status, err := RunCommand("status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to check git status: %v", err)
	}

	if status == "" {
		return fmt.Errorf("no changes to checkpoint")
	}

	// Add all changes
	_, err = RunCommand("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add changes: %v", err)
	}

	// Create commit message
	timestamp := GetTimestamp()
	var message string
	if customNote != "" {
		message = fmt.Sprintf("CHECKPOINT: %s - %s", timestamp, customNote)
	} else {
		message = fmt.Sprintf("CHECKPOINT: %s", timestamp)
	}

	// Create commit
	_, err = RunCommand("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create checkpoint: %v", err)
	}

	return nil
}

// GetCheckpointsFromHistory returns checkpoints from commit history
func GetCheckpointsFromHistory() ([]models.Checkpoint, error) {
	output, err := RunCommand("log", "--oneline", "--grep=CHECKPOINT", "--format=%h %s")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []models.Checkpoint{}, nil
	}

	lines := strings.Split(output, "\n")
	var checkpoints []models.Checkpoint

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 || !strings.Contains(parts[1], "CHECKPOINT:") {
			continue
		}

		checkpoints = append(checkpoints, models.Checkpoint{
			Hash:    parts[0],
			Message: parts[1],
		})
	}

	return checkpoints, nil
}

// GetCheckpointsFromReflog returns checkpoints from reflog (includes navigation history)
func GetCheckpointsFromReflog() ([]models.Checkpoint, error) {
	// Get all reflog entries to find all commits (including unreachable ones)
	reflogOutput, err := RunCommand("reflog", "--format=%h %s")
	if err != nil {
		return nil, err
	}

	var checkpoints []models.Checkpoint
	seen := make(map[string]bool)

	if strings.TrimSpace(reflogOutput) != "" {
		lines := strings.Split(reflogOutput, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}

			hash := parts[0]
			message := parts[1]

			// Only include checkpoints and avoid duplicates
			if strings.Contains(message, "CHECKPOINT:") && !seen[hash] {
				seen[hash] = true
				checkpoints = append(checkpoints, models.Checkpoint{
					Hash:    hash,
					Message: message,
				})
			}
		}
	}

	// Reverse to show newest first (reflog is oldest to newest)
	for i, j := 0, len(checkpoints)-1; i < j; i, j = i+1, j-1 {
		checkpoints[i], checkpoints[j] = checkpoints[j], checkpoints[i]
	}

	// Add the last non-checkpoint commit at the end if it exists
	lastNonCheckpoint, err := GetLastNonCheckpointCommit()
	if err == nil && lastNonCheckpoint != nil {
		// Check if this commit is not already in the list
		if !seen[lastNonCheckpoint.Hash] {
			// Mark it as non-checkpoint for UI display
			lastNonCheckpoint.Message = "[LAST COMMIT] " + lastNonCheckpoint.Message
			checkpoints = append(checkpoints, *lastNonCheckpoint)
		}
	}

	return checkpoints, nil
}

// SwitchToCheckpoint switches to a specific checkpoint
func SwitchToCheckpoint(hash string) error {
	if !IsRepo() {
		return fmt.Errorf("not in a Git repository")
	}

	// Check if we're already on this checkpoint
	currentCommit, err := GetCurrentCommit()
	if err == nil && currentCommit == hash {
		return fmt.Errorf("you are already on checkpoint %s", hash)
	}

	_, err = RunCommand("checkout", hash)
	if err != nil {
		return fmt.Errorf("failed to switch to checkpoint %s: %v", hash, err)
	}

	return nil
}

// GetLastNonCheckpointCommit finds the last commit that is not a checkpoint
func GetLastNonCheckpointCommit() (*models.Checkpoint, error) {
	output, err := RunCommand("log", "--oneline", "--format=%h %s")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		// If this commit is NOT a checkpoint, return it
		if !strings.Contains(parts[1], "CHECKPOINT:") {
			return &models.Checkpoint{
				Hash:    parts[0],
				Message: parts[1],
			}, nil
		}
	}

	return nil, nil // No non-checkpoint commit found
}

// HasCheckpoints returns true if any checkpoints exist
func HasCheckpoints() bool {
	checkpoints, err := GetCheckpointsFromReflog()
	if err != nil {
		return false
	}
	
	// Count only actual checkpoints (not the last non-checkpoint commit)
	checkpointCount := 0
	for _, cp := range checkpoints {
		if strings.Contains(cp.Message, "CHECKPOINT:") {
			checkpointCount++
		}
	}
	
	return checkpointCount > 0
}

// IsCurrentCommitCheckpoint returns true if the current commit is a checkpoint
func IsCurrentCommitCheckpoint() bool {
	// Get current commit message
	output, err := RunCommand("log", "--oneline", "-1", "--format=%s")
	if err != nil {
		return false
	}
	
	return strings.Contains(output, "CHECKPOINT:")
}