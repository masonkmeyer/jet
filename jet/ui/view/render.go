package view

import "github.com/jroimartin/gocui"

type (
	RenderOption func(*gocui.View)
	Renderer     interface {
		Render(*gocui.View, ...RenderOption) error
	}
)

// WithWrap sets the wrap option for a view
func WithWrap(val bool) RenderOption {
	return func(v *gocui.View) {
		v.Wrap = val
	}
}

// WithEditable sets the editable option for a view
func WithEditable(val bool) RenderOption {
	return func(v *gocui.View) {
		v.Editable = val
	}
}

// WithAutoscroll sets the autoscroll option for a view
func WithAutoscroll(val bool) RenderOption {
	return func(v *gocui.View) {
		v.Autoscroll = val
	}
}

// WithHighlight sets the highlight option for a view
func WithHighlight(val bool) RenderOption {
	return func(v *gocui.View) {
		v.Highlight = val
	}
}

// WithSelBgColor sets the selected background color for a view
func WithSelBgColor(color gocui.Attribute) RenderOption {
	return func(v *gocui.View) {
		v.SelBgColor = color
	}
}

// WithSelFgColor sets the selected foreground color for a view
func WithSelFgColor(color gocui.Attribute) RenderOption {
	return func(v *gocui.View) {
		v.SelFgColor = color
	}
}

// WithTitle sets the title for a view
func WithTitle(title string) RenderOption {
	return func(v *gocui.View) {
		v.Title = title
	}
}

// WithFgColor sets the foreground color for a view
func WithFgColor(color gocui.Attribute) RenderOption {
	return func(v *gocui.View) {
		v.FgColor = color
	}
}
