package jet

import (
	"fmt"
	"os/exec"
	"strings"
)

type (
	// Git is a wrapper around the git command
	Git struct {
	}
)

// Checkout executes a git checkout command with the branch name
// and returns the output and the command that was executed
func (g Git) Checkout(branchName string) (string, string) {
	result, cmd := g.exec("checkout", branchName)
	return result, cmd
}

// ListBranches executes a git branch command with --list and any other provided args
func (g Git) ListBranches(gitBranchArgs ...string) []string {
	output, _ := g.exec(append([]string{"branch", "--list"}, gitBranchArgs...)...)
	return strings.Split(output, "\n")
}

// Log executes a git log command with any provided args and --no-pager to prevent interactive paging
func (g Git) Logs(args ...string) string {
	output, _ := g.exec(append([]string{"--no-pager", "log"}, args...)...)
	return output
}

// exec executes a git command with the provided args
// it returns the output of the command and the command that was executed
func (g Git) exec(args ...string) (string, string) {
	out, _ := exec.Command("git", args...).CombinedOutput()
	return string(out), fmt.Sprintf("git %s", strings.Join(args, " "))
}
