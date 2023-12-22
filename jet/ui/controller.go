package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet"
	"github.com/masonkmeyer/jet/jet/ui/view"
	"github.com/masonkmeyer/jet/jet/ui/viewmodel"
)

// constants for view names
const (
	BRANCHES    = "branches"
	LOGS        = "logs"
	RECENT_LOGS = "recent_logs"
)

// Controller is the main controller for the UI
type Controller struct {
	g                   *gocui.Gui
	git                 *jet.Git
	recentCommitMessage string
	gitGraph            string
	exitChannel         chan string
}

// NewController creates a new UI controller
func NewController(g *gocui.Gui, exitChannel chan string) *Controller {
	c := &Controller{
		g:                   g,
		git:                 &jet.Git{},
		recentCommitMessage: "",
		gitGraph:            "",
		exitChannel:         exitChannel,
	}

	return c
}

// Quit is the handler for the quit keybinding
func (c *Controller) Quit(g *gocui.Gui, v *gocui.View) error {
	go func() { c.exitChannel <- "" }()
	return gocui.ErrQuit
}

// Layout is the handler for the layout of the UI
func (c *Controller) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	items := []*viewmodel.MenuItem{}

	branches := c.git.ListBranches("--sort=-committerdate", "--format=%(align:13,left)%(authordate:relative)%(end) %(refname:short)")

	if strings.Contains(branches[0], "fatal") {
		go func() { c.exitChannel <- "Error: Not a git repository" }()
		return gocui.ErrQuit
	}

	for _, branch := range branches {
		parts := strings.Split(branch, " ")
		value := parts[len(parts)-1]

		items = append(items, &viewmodel.MenuItem{
			Title: c.pad(branch, " ", maxX/2-len(branch)),
			Value: value,
		})
	}

	menu := viewmodel.Menu{
		Items:      items,
		OnSelected: c.onSelected,
		OnChange:   c.onChange,
	}

	if v, err := g.SetView(BRANCHES, 0, 0, maxX/2-1, maxY/2-1); err != nil {
		menuView, _ := view.NewMenu(g, menu, BRANCHES)
		menuView.Render(v,
			view.WithHighlight(true),
			view.WithSelBgColor(gocui.ColorCyan),
			view.WithSelFgColor(gocui.ColorBlack),
			view.WithTitle("Recent Branches"))
	}

	if v, err := g.SetView(LOGS, maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.recentCommitMessage}, LOGS)
		textView.Render(v, view.WithWrap(true), view.WithTitle("Recent Commit Message"), view.WithFgColor(gocui.ColorYellow))
	}

	if v, err := g.SetView(RECENT_LOGS, 0, maxY/2, maxX-1, maxY-1); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.gitGraph}, LOGS)
		textView.Render(v, view.WithWrap(true), view.WithTitle("Recent Commits"))
	}

	g.SetCurrentView(BRANCHES)

	return nil
}

// onSelected is the handler for when a branch is selected
func (c *Controller) onSelected(item *viewmodel.MenuItem) error {
	output, cmd := c.git.Checkout(item.Value)

	go func() { c.exitChannel <- fmt.Sprintf(">> %s \n %s", cmd, output) }()

	return gocui.ErrQuit
}

// onChange is the handler for when the selected branch changes
func (c *Controller) onChange(item *viewmodel.MenuItem) error {

	results := c.git.Logs(item.Value, "-1")

	c.recentCommitMessage = results
	c.g.DeleteView(LOGS)

	graphLog := c.git.Logs(item.Value, "-n 30", "--oneline", "--decorate", "--color", "--abbrev-commit", "--date=relative", "--format=format:%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)")
	c.gitGraph = graphLog
	c.g.DeleteView(RECENT_LOGS)
	return nil
}

// Pad the string with chars of length
func (c *Controller) pad(s string, padStr string, pLen int) string {
	if pLen <= 0 {
		return s
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(padStr, pLen))
}
