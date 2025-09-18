package ui

import "github.com/charmbracelet/bubbles/key"

func (i nameItem) Title() string {
	return i.item.Name
}

func (i nameItem) Description() string {
	return i.item.Mykad
}

func (i nameItem) FilterValue() string {
	return i.item.Mykad
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleNextPage: key.NewBinding(
			key.WithHelp("n", "next page"),
			key.WithKeys("n"),
		),
		togglePreviousPage: key.NewBinding(
			key.WithHelp("p", "previous page"),
			key.WithKeys("p"),
		),
		toggleMenu: key.NewBinding(
			key.WithHelp("m", "back to menu"),
			key.WithKeys("m"),
		),
	}
}

func (i menuItem) Title() string       { return i.Name }
func (i menuItem) Description() string { return i.Desc }
func (i menuItem) FilterValue() string { return i.Name }

func (i historyItem) Title() string {
	return i.item.Query
}

func (i historyItem) Description() string {
	return i.item.Result
}

func (i historyItem) FilterValue() string {
	return i.item.Query
}
