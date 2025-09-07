package ui

import (
	"synthera/api"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	State AppState
	TokenInput textinput.Model
	NameInput textinput.Model
	UserDetails *api.TraceDetailID
	UserID int
	Relations *api.TraceRelationsItem
	APIToken string
	ErrorMessage string
	APIClient *api.Client
	Offset int

	Spinner spinner.Model
	Doc lipgloss.Style
	List list.Model
	Width int
	Height int
	Logo string
}

type nameItem struct {
	item api.TraceNameItem
}

type nameListKeyMap struct {
	toggleNextPage key.Binding
	togglePreviousPage key.Binding
}
