package git

import (
	"fmt"
	"strings"
	"time"
	"vibe-check/internal/models"
)

// FinalizeAndPush squashes consecutive checkpoints and pushes to remote
func FinalizeAndPush() error {
	return FinalizeAndPushWithMessage("")
}

// FinalizeAndPushWithMessage squashes consecutive checkpoints and pushes to remote with custom message
func FinalizeAndPushWithMessage(customMessage string) error {
	if !IsRepo() {
		return fmt.Errorf("not in a Git repository")
	}

	// Check if we're in detached HEAD state and fix it
	detectedBranch, branchErr := GetCurrentBranch()
	if branchErr == nil && detectedBranch == "HEAD" {
		// We're in detached HEAD - switch back to main/master branch
		_, err := RunCommand("checkout", "main")
		if err != nil {
			// Try master if main doesn't exist
			_, err = RunCommand("checkout", "master")
			if err != nil {
				return fmt.Errorf("in detached HEAD state and cannot switch to main/master branch. Please checkout a branch first: %v", err)
			}
		}
	}

	// Get current commit hash
	currentCommit, err := GetCurrentCommit()
	if err != nil {
		return fmt.Errorf("error getting current commit: %v", err)
	}

	// Get all checkpoints from reflog to include navigation history
	checkpoints, err := GetCheckpointsFromReflog()
	if err != nil {
		return fmt.Errorf("error getting checkpoints: %v", err)
	}

	if len(checkpoints) == 0 {
		return fmt.Errorf("no checkpoints found")
	}

	// Find current position in checkpoint list
	var currentIndex = -1
	for i, cp := range checkpoints {
		if cp.Hash == currentCommit {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return fmt.Errorf("current commit is not a checkpoint. Push feature only works from checkpoints")
	}

	// Find consecutive checkpoints going backwards from current position
	consecutiveCheckpoints := []models.Checkpoint{checkpoints[currentIndex]}
	
	// Go backwards and collect consecutive checkpoints
	for i := currentIndex + 1; i < len(checkpoints); i++ {
		// Check if this is still a checkpoint (has CHECKPOINT: prefix)
		if !strings.Contains(checkpoints[i].Message, "CHECKPOINT:") {
			break
		}
		consecutiveCheckpoints = append(consecutiveCheckpoints, checkpoints[i])
	}

	// Find the commit before the oldest checkpoint we're squashing
	var baseCommit string
	if currentIndex + len(consecutiveCheckpoints) < len(checkpoints) {
		// There's a commit before our checkpoints
		baseCommit = checkpoints[currentIndex + len(consecutiveCheckpoints)].Hash
	} else {
		// We're at the very beginning, find first non-checkpoint commit
		output, err := RunCommand("log", "--oneline", "--format=%H %s")
		if err != nil {
			return fmt.Errorf("error getting commit history: %v", err)
		}

		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if !strings.Contains(line, "CHECKPOINT:") {
				parts := strings.SplitN(line, " ", 2)
				if len(parts) >= 1 {
					baseCommit = parts[0]
					break
				}
			}
		}
	}

	if baseCommit == "" {
		return fmt.Errorf("cannot find base commit for squashing. All commits appear to be checkpoints")
	}
	

	// Remove any checkpoints newer than current position (if any)
	var removedNewer int
	for i := 0; i < currentIndex; i++ {
		removedNewer++
	}

	// Create backup branch
	backupBranch := fmt.Sprintf("vibe-check-backup-%d", time.Now().Unix())
	_, err = RunCommand("branch", backupBranch, "HEAD")
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// Soft reset to base commit to preserve changes but remove checkpoint commits
	_, err = RunCommand("reset", "--soft", baseCommit)
	if err != nil {
		// Restore backup on failure
		RunCommand("reset", "--hard", backupBranch)
		RunCommand("branch", "-D", backupBranch)
		return fmt.Errorf("failed to reset to base: %v", err)
	}

	// Generate commit message (custom or automatic)
	var commitMessage string
	if customMessage != "" {
		commitMessage = customMessage
	} else {
		// Use timestamp-based auto message
		timestamp := GetTimestamp()
		commitMessage = fmt.Sprintf("Update: %s", timestamp)
	}

	// Check if there are changes to commit after soft reset
	// Use --cached to check staged changes specifically
	stagedStatus, statusErr := RunCommand("diff", "--cached", "--name-only")
	if statusErr == nil && strings.TrimSpace(stagedStatus) == "" {
		// No staged changes after soft reset - check if working directory has changes
		workingStatus, _ := RunCommand("status", "--porcelain")
		if strings.TrimSpace(workingStatus) == "" {
			// No changes at all - this means checkpoints were identical to base
			RunCommand("branch", "-D", backupBranch)
			return fmt.Errorf("no changes to commit after squashing checkpoints. This usually means:\n" +
				"1. All checkpoints had identical content to the base commit\n" +
				"2. The soft reset resulted in no differences\n" +
				"Solution: Your checkpoints have been consolidated - no new commit was needed")
		} else {
			// Working directory has changes but nothing staged - stage them
			RunCommand("add", ".")
		}
	}

	// Create the final commit
	output, err := RunCommand("commit", "-m", commitMessage)
	if err != nil {
		// Restore backup on failure
		RunCommand("reset", "--hard", backupBranch)
		RunCommand("branch", "-D", backupBranch)
		
		diagnosis := diagnoseCommitError(output, err)
		return fmt.Errorf("failed to create final commit:\n%s\n\nDiagnosis: %s", err, diagnosis)
	}

	// Get current branch name for push
	currentBranch, branchErr := GetCurrentBranch()
	if branchErr != nil {
		currentBranch = "main" // fallback to main if can't detect
	}
	
	// Push to remote (force with lease for safety when rewriting history)
	pushOutput, err := RunCommand("push", "--force-with-lease", "origin", currentBranch)
	if err != nil {
		// Don't restore backup here - commit was successful, just push failed
		RunCommand("branch", "-D", backupBranch)
		
		// Provide detailed error diagnosis
		diagnosis := diagnosePushError(pushOutput, err)
		return fmt.Errorf("commit created successfully but push failed:\nError: %s\nOutput: %s\n\nDiagnosis: %s\n\nNote: You can manually push with:\ngit push --force-with-lease origin %s", err, pushOutput, diagnosis, currentBranch)
	}

	// Clean up backup branch
	RunCommand("branch", "-D", backupBranch)

	// Clean up reflog ONLY after successful push
	// This removes ALL checkpoints from the UI lists (including current)
	RunCommand("reflog", "expire", "--expire=now", "--all")
	RunCommand("gc", "--prune=now")

	return nil
}

// GetFinalizeInfo returns information about what would be finalized
func GetFinalizeInfo() (string, error) {
	if !IsRepo() {
		return "", fmt.Errorf("not in a Git repository")
	}

	currentCommit, err := GetCurrentCommit()
	if err != nil {
		return "", fmt.Errorf("error getting current commit: %v", err)
	}

	checkpoints, err := GetCheckpointsFromReflog()
	if err != nil {
		return "", fmt.Errorf("error getting checkpoints: %v", err)
	}

	if len(checkpoints) == 0 {
		return "No checkpoints found", nil
	}

	// Find current position
	var currentIndex = -1
	for i, cp := range checkpoints {
		if cp.Hash == currentCommit {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return "Current commit is not a checkpoint", nil
	}

	// Count consecutive checkpoints
	consecutiveCount := 1
	for i := currentIndex + 1; i < len(checkpoints) && strings.Contains(checkpoints[i].Message, "CHECKPOINT:"); i++ {
		consecutiveCount++
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("Will squash %d consecutive checkpoints\n", consecutiveCount))
	
	if currentIndex > 0 {
		info.WriteString(fmt.Sprintf("Will remove %d newer checkpoints\n", currentIndex))
	}
	
	info.WriteString("\nCheckpoints to be squashed:\n")
	for i := 0; i < consecutiveCount && currentIndex+i < len(checkpoints); i++ {
		cp := checkpoints[currentIndex+i]
		marker := "  "
		if i == 0 {
			marker = "> " // current position
		}
		info.WriteString(fmt.Sprintf("%s[%s] %s\n", marker, cp.Hash, cp.Message))
	}

	return info.String(), nil
}


// diagnosePushError analyzes push failure and provides helpful suggestions
func diagnosePushError(output string, err error) string {
	errorMsg := strings.ToLower(output + " " + err.Error())
	
	if strings.Contains(errorMsg, "not a full refname") || strings.Contains(errorMsg, "refs/heads") {
		return "Git refname issue (HEAD push problem). Solutions:\n" +
			"1. This should now be fixed - using branch name instead of HEAD\n" +
			"2. Check current branch: git branch\n" +
			"3. Make sure you're on a proper branch, not detached HEAD\n" +
			"4. Try: git checkout -b main (if no branch exists)"
	}
	
	if strings.Contains(errorMsg, "permission denied") || strings.Contains(errorMsg, "authentication") {
		return "Authentication failed. Solutions:\n" +
			"1. Check if you have push access to this repository\n" +
			"2. Verify your Git credentials: git config --list | grep user\n" +
			"3. For GitHub, check if you need a personal access token\n" +
			"Manual push: git push --force-with-lease origin main"
	}
	
	if strings.Contains(errorMsg, "no such remote") || strings.Contains(errorMsg, "does not exist") {
		return "Remote repository not found. Solutions:\n" +
			"1. Check remote URL: git remote -v\n" +
			"2. Add remote if missing: git remote add origin <your-repo-url>\n" +
			"3. Update remote URL: git remote set-url origin <correct-url>\n" +
			"Manual push: git push --force-with-lease origin main"
	}
	
	if strings.Contains(errorMsg, "rejected") || strings.Contains(errorMsg, "non-fast-forward") {
		return "Push rejected (remote has newer commits). Solutions:\n" +
			"1. Someone else pushed to the repository\n" +
			"2. Pull latest changes: git pull origin main\n" +
			"3. Then retry finalize and push\n" +
			"Manual push: git push --force-with-lease origin main"
	}
	
	if strings.Contains(errorMsg, "network") || strings.Contains(errorMsg, "timeout") {
		return "Network connection failed. Solutions:\n" +
			"1. Check your internet connection\n" +
			"2. Try again in a moment\n" +
			"3. Check if GitHub/GitLab is accessible\n" +
			"Manual push: git push --force-with-lease origin main"
	}
	
	return "Unknown push error. Solutions:\n" +
		"1. Check repository access and credentials\n" +
		"2. Verify remote repository exists\n" +
		"3. Try manual push: git push --force-with-lease origin main\n" +
		"4. Check git status and git remote -v"
}

// diagnoseCommitError analyzes commit failure and provides helpful suggestions
func diagnoseCommitError(output string, err error) string {
	errorMsg := strings.ToLower(output + " " + err.Error())
	
	if strings.Contains(errorMsg, "nothing to commit") || strings.Contains(errorMsg, "no changes") {
		return "No changes to commit. This means:\n" +
			"1. All checkpoints had identical content\n" +
			"2. The working directory is clean after squashing\n" +
			"Solution: This is normal - your checkpoints have been consolidated"
	}
	
	if strings.Contains(errorMsg, "pre-commit") || strings.Contains(errorMsg, "hook") {
		return "Pre-commit hook failed. Solutions:\n" +
			"1. Fix the issues reported by the pre-commit hook\n" +
			"2. Or bypass hooks temporarily: git commit --no-verify -m \"your message\"\n" +
			"3. Check what hooks are configured in .git/hooks/"
	}
	
	if strings.Contains(errorMsg, "index.lock") || strings.Contains(errorMsg, "unable to create") {
		return "Git index locked or permission issue. Solutions:\n" +
			"1. Remove lock file: rm .git/index.lock\n" +
			"2. Check file permissions in .git directory\n" +
			"3. Try the operation again"
	}
	
	if strings.Contains(errorMsg, "pathspec") || strings.Contains(errorMsg, "did not match") {
		return "File path issue. Solutions:\n" +
			"1. Check if all files exist\n" +
			"2. Verify working directory is correct\n" +
			"3. Run git status to see current state"
	}
	
	return "Unknown commit error. Solutions:\n" +
		"1. Check git status\n" +
		"2. Verify repository is in good state\n" +
		"3. Try manual commit: git commit -m \"Manual squash\"\n" +
		"4. Check .git/hooks for problematic hooks"
}