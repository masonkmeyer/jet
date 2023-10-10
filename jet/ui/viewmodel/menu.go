package viewmodel

type (
	Menu struct {
		Items      []*MenuItem
		OnSelected func(item *MenuItem) error
	}

	MenuItem struct {
		Title string
		Value string
	}
)
