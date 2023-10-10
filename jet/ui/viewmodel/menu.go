package viewmodel

type (
	// Menu is the viewmodel for a menu
	Menu struct {
		// Items is the list of menu items
		Items []*MenuItem
		// OnSelected is the function to call when a menu item is selected
		OnSelected func(item *MenuItem) error
		// OnChange is the function to call when the selected menu item changes
		OnChange func(item *MenuItem) error
	}

	// MenuItem is a single item in a menu
	MenuItem struct {
		Title string
		Value string
	}
)
