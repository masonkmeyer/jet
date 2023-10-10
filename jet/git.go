package jet

import (
	"os/exec"
	"strings"
)

type Git struct {
}

func (g Git) exec(args ...string) string {
	out, _ := exec.Command("git", args...).CombinedOutput()
	return string(out)
}

func (g Git) Checkout(name string) string {
	return g.exec("checkout", name)
}

func (g Git) ListBranches(gitBranchArgs ...string) []string {
	output := g.exec(append([]string{"branch", "--list"}, gitBranchArgs...)...)
	return strings.Split(output, "\n")
}
