package git

import (
	"os/exec"
	"strings"
)

// IsRepo checks if the current directory is a Git repository
func IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// RunCommand executes a git command and returns the output
func RunCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// GetCurrentCommit returns the current commit hash (short format)
func GetCurrentCommit() (string, error) {
	return RunCommand("rev-parse", "--short", "HEAD")
}

// GetTimestamp returns current timestamp in desired format
func GetTimestamp() string {
	cmd := exec.Command("date", "+%d/%m/%Y %H:%M")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() (string, error) {
	return RunCommand("rev-parse", "--abbrev-ref", "HEAD")
}

// HasUncommittedChanges returns true if there are uncommitted changes
func HasUncommittedChanges() bool {
	status, err := RunCommand("status", "--porcelain")
	if err != nil {
		return false
	}
	return strings.TrimSpace(status) != ""
}