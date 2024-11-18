package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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
	FILTER      = "filter"
)

// Controller is the main controller for the UI
type Controller struct {
	filter              string
	g                   *gocui.Gui
	git                 *jet.Git
	recentCommitMessage string
	gitGraph            string
	exitChannel         chan string
}

// NewController creates a new UI controller
func NewController(g *gocui.Gui, exitChannel chan string) (*Controller, error) {
	repo, err := jet.NewGit()
	if err != nil {
		return nil, err
	}

	c := &Controller{
		filter:              "",
		g:                   g,
		git:                 repo,
		recentCommitMessage: "",
		gitGraph:            "",
		exitChannel:         exitChannel,
	}

	return c, nil
}

func (c *Controller) OnBackspace(g *gocui.Gui, v *gocui.View) error {
	if len(c.filter) == 0 {
		return nil
	}

	c.filter = c.filter[:len(c.filter)-1]
	g.DeleteView(BRANCHES)
	g.DeleteView(FILTER)

	return nil
}

func (c *Controller) OnType(char rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		c.filter += string(char)
		g.DeleteView(BRANCHES)
		g.DeleteView(FILTER)
		return nil
	}
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

	branches := c.git.ListBranches()

	for _, branch := range branches {
		if c.filter != "" && !strings.Contains(branch.Name, c.filter) {
			continue
		}

		elapsed := jet.Ago(branch.When)
		branchName := fmt.Sprintf(" %s (%s)", branch.Name, elapsed)
		items = append(items, &viewmodel.MenuItem{
			Title: c.pad(branchName, " ", maxX/2-len(branchName)),
			Value: branch.Name,
		})
	}

	menu := viewmodel.Menu{
		Items:      items,
		OnSelected: c.onSelected,
		OnChange:   c.onChange,
	}

	if v, err := g.SetView(BRANCHES, 0, 0, maxX/2-1, maxY/2-5); err != nil {
		menuView, _ := view.NewMenu(g, menu, BRANCHES)
		menuView.Render(v,
			view.WithHighlight(true),
			view.WithSelBgColor(gocui.ColorCyan),
			view.WithSelFgColor(gocui.ColorBlack),
			view.WithTitle("Recent Branches"))
	}

	if v, err := g.SetView(LOGS, maxX/2, 0, maxX-1, maxY/2-5); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.recentCommitMessage}, LOGS)
		textView.Render(v, view.WithWrap(true), view.WithTitle("Recent Commit Message"))
	}

	if v, err := g.SetView(RECENT_LOGS, 0, maxY/2-4, maxX-1, maxY-4); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.gitGraph}, LOGS)
		textView.Render(v, view.WithWrap(true), view.WithTitle("Recent Commits"))
	}

	if v, err := g.SetView(FILTER, 0, maxY-3, maxX-1, maxY-1); err != nil {
		textView := view.NewText(g, viewmodel.Text{Value: c.filter}, FILTER)
		textView.Render(v, view.WithTitle("Filter"))
	}

	g.SetCurrentView(BRANCHES)

	return nil
}

// onSelected is the handler for when a branch is selected
func (c *Controller) onSelected(item *viewmodel.MenuItem) error {
	c.git.Checkout(item.Value)

	go func() { c.exitChannel <- fmt.Sprintf(">> checkout %s ", item.Value) }()

	return gocui.ErrQuit
}

// onChange is the handler for when the selected branch changes
func (c *Controller) onChange(item *viewmodel.MenuItem) error {
	results := c.git.Logs(item.Value, 10)

	if len(results) == 0 {
		return nil
	}

	recentCommit := results[0]

	c.recentCommitMessage = format(recentCommit)
	c.g.DeleteView(LOGS)

	remainingCommits := results[1:]
	commitList := []string{}
	for _, commit := range remainingCommits {
		commitList = append(commitList, format(commit))
	}

	c.gitGraph = strings.Join(commitList, "\n\n")

	c.g.DeleteView(RECENT_LOGS)

	return nil
}

// format formats a commit for display
func format(c jet.Commit) string {
	return fmt.Sprintf("\t%s - (%s) %s\n %s", color.BlueString(c.Hash), color.GreenString(jet.Ago(c.When)), c.Author, color.WhiteString(removeLineBreaks(c.Message)))
}

// remove starting and ending line breaks
func removeLineBreaks(s string) string {
	s = strings.TrimPrefix(s, "\n")
	s = strings.TrimSuffix(s, "\n")
	s = strings.ReplaceAll(s, "\n", "\n\t")
	return s
}

func (c *Controller) pad(s string, padStr string, pLen int) string {
	if pLen <= 0 {
		return s
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(padStr, pLen))
}
