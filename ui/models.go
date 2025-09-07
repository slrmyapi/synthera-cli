// Package ui for user interface
package ui

import (
	"fmt"
	"synthera/api"
	"synthera/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppState int

const (
	StateTokenInput AppState = iota
	StateTraceNameInput
	StateTraceNameResults
	StateTraceIDInput
	StateTraceDetails
	StateError
	StateLoading
)

func InitialModel(initialToken string, logo string) MainModel {
	var state AppState

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.Spinner = spinner.Monkey

	docStyle := lipgloss.NewStyle().Margin(1, 2)

	tokenInput := textinput.New()
	tokenInput.Placeholder = "Enter your API token"
	tokenInput.CharLimit = 256
	tokenInput.Width = 50
	tokenInput.Cursor.Blink = true

	nameInput := textinput.New()
	nameInput.Placeholder = "Search a name"
	nameInput.CharLimit = 256
	nameInput.Width = 50
	nameInput.Cursor.Blink = true

	if initialToken == "" {
		state = StateTokenInput
		tokenInput.Focus()
	} else {
		state = StateTraceNameInput
		nameInput.Focus()
	}

	return MainModel{
		State: state,
		Spinner: s,
		Doc: docStyle,
		TokenInput: tokenInput,
		NameInput: nameInput,
		APIClient: api.NewClient(initialToken),
		APIToken: initialToken,
		Logo: logo,
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return m, tea.Quit
		}

		switch m.State {
		case StateTokenInput:
			switch msg.Type {
			case tea.KeyEnter:
				m.APIToken = m.TokenInput.Value()
				m.APIClient = api.NewClient(m.APIToken)
				m.State = StateTraceNameInput
				m.NameInput.Focus()
				utils.SaveToken(m.APIToken)
			}

			m.TokenInput, cmd = m.TokenInput.Update(msg)
		case StateTraceNameInput:
			switch msg.Type {
			case tea.KeyEnter:
				name := m.NameInput.Value()
				m.State = StateLoading
				return m, m.FetchName(name)
			}

			m.NameInput, cmd = m.NameInput.Update(msg)
		case StateTraceNameResults:
			switch msg.Type {
			case tea.KeyEnter:
				if item, ok := m.List.SelectedItem().(nameItem); ok {
					m.State = StateLoading
					return m, m.FetchID(item.item.ID)
				}
			}
			m.List, cmd = m.List.Update(msg)
		case StateTraceDetails:
			m.Relations = nil
			switch msg.String() {
			case "n", "N":
				m.State = StateLoading
				m.Offset += 1
				return m, m.FetchRelations(m.UserID)
			case "p", "P":
				m.State = StateLoading
				m.Offset -= 1
				return m, m.FetchRelations(m.UserID)
			default:
				m.State = StateTraceNameInput
			}
		case StateError:
			m.State = StateTraceNameInput
		}
	case tea.WindowSizeMsg:
		h, v := m.Doc.GetFrameSize()
		m.Width = msg.Width-h
		m.Height = msg.Height-v
		switch m.State {
		case StateTraceNameResults:
			m.List.SetSize(m.Width, m.Height)
		}
	case api.TraceNameMsg:
		if msg.Err != nil {
			m.State = StateError
			m.ErrorMessage = fmt.Sprintf("Error fetching names: %s", msg.Err.Error())
			return m, nil
		}

		if len(msg.Items) == 0 {
			m.State = StateError
			m.ErrorMessage = "No results found for that name"
		} else {
			var items []list.Item
			for _, item := range msg.Items {
				items = append(items, nameItem{
					item: item,
				})
			}

			listKeys := newNameListKeyMap()
			m.List = list.New(items, list.NewDefaultDelegate(), m.Width, m.Height)
			m.List.Title = "Trace Results"
			m.List.AdditionalFullHelpKeys = func() []key.Binding {
				return []key.Binding{
					listKeys.toggleNextPage,
					listKeys.togglePreviousPage,
				}
			}
			m.State = StateTraceNameResults
		}
	case api.TraceDetailsMsg:
		if msg.Err != nil {
			m.State = StateError
			m.ErrorMessage = fmt.Sprintf("Error fetching details: %s", msg.Err.Error())
			return m, nil
		}

		if len(msg.Details) == 0 {
			m.State = StateError
			m.ErrorMessage = "No results found"
		} else {
			m.UserDetails = &msg.Details[0]
			m.UserID = m.UserDetails.ID
			m.Offset = 0
			m.State = StateTraceDetails
		}
	case api.TraceRelationsMsg:
		if msg.Err != nil {
			m.State = StateError
			m.ErrorMessage = fmt.Sprintf("Can't find relationships: %s", msg.Err.Error())
			return m, nil
		}

		if len(msg.Details) == 0 {
			m.State = StateError
			m.ErrorMessage = "No results found"
		} else {
			m.UserDetails = &msg.Details[0]
			m.Relations = &msg.Relations
			m.State = StateTraceDetails
		}
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
	}
	return m, cmd
}

func (m MainModel) FetchName(name string) tea.Cmd {
	return func() tea.Msg {
		items, err := m.APIClient.TraceName(name)
		return api.TraceNameMsg{
			Items: items,
			Err: err,
		}
	}
}

func (m MainModel) FetchID(id int) tea.Cmd {
	return func() tea.Msg {
		details, err := m.APIClient.TraceDetail(id)
		return api.TraceDetailsMsg{
			Details: details,
			Err: err,
		}
	}
}

func (m MainModel) FetchRelations(id int) tea.Cmd {
	return func() tea.Msg {
		details, relations, err := m.APIClient.TraceRelations(id, m.Offset)
		return api.TraceRelationsMsg{
			Details: details,
			Relations: relations,
			Err: err,
		}
	}
}
