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

func newNameListKeyMap() *nameListKeyMap {
	return &nameListKeyMap{
		toggleNextPage: key.NewBinding(
			key.WithHelp("n", "next page"),
			key.WithKeys("n"),
		),
		togglePreviousPage: key.NewBinding(
			key.WithHelp("p", "previous page"),
			key.WithKeys("p"),
		),
	}
}
