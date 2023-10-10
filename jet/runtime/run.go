package runtime

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/masonkmeyer/jet/jet/ui"
)

// Run starts the UI
func Run(exitChannel chan string) {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	controller := ui.NewController(g, exitChannel)

	g.SetManager(controller)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, controller.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
