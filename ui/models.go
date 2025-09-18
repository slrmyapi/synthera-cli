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
	StateMainMenu
	StateTraceNRICInput
	StateHistory
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
	tokenInput.Focus()

	nameInput := textinput.New()
	nameInput.Placeholder = "Search a name"
	nameInput.CharLimit = 256
	nameInput.Width = 50
	nameInput.Cursor.Blink = true
	nameInput.Focus()

	nricInput := textinput.New()
	nricInput.Placeholder = "Search IC"
	nricInput.CharLimit = 256
	nricInput.Width = 50
	nricInput.Cursor.Blink = true
	nricInput.Focus()

	items := []list.Item{
		menuItem{
			Name:  "Trace Name",
			Desc:  "Find people using their name",
			State: StateTraceNameInput,
		},
		menuItem{
			Name:  "Trace Mykad",
			Desc:  "Find people using their mykad",
			State: StateTraceNRICInput,
		},
		menuItem{
			Name:  "History",
			Desc:  "Your search history",
			State: StateHistory,
		},
		menuItem{
			Name:  "Replace Token",
			Desc:  "Replace your old token with a new one",
			State: StateTokenInput,
		},
	}
	menuList := list.New(items, list.NewDefaultDelegate(), 0, 0)

	if initialToken == "" {
		state = StateTokenInput
		tokenInput.Focus()
	} else {
		state = StateMainMenu
	}

	return MainModel{
		State:      state,
		Spinner:    s,
		Doc:        docStyle,
		TokenInput: tokenInput,
		NameInput:  nameInput,
		NRICInput:  nricInput,
		APIClient:  api.NewClient(initialToken),
		APIToken:   initialToken,
		Logo:       logo,
		Menu:       menuList,
		Page:       1,
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
				m.State = StateMainMenu
				m.NameInput.Focus()
				utils.SaveToken(m.APIToken)
			}

			m.TokenInput, cmd = m.TokenInput.Update(msg)
		case StateTraceNameInput:
			switch msg.Type {
			case tea.KeyEnter:
				name := m.NameInput.Value()
				m.State = StateLoading
				m.Page = 1
				return m, m.FetchName(name)
			}

			m.NameInput, cmd = m.NameInput.Update(msg)
		case StateTraceNameResults:
			switch msg.String() {
			case "enter":
				if item, ok := m.List.SelectedItem().(nameItem); ok {
					m.State = StateLoading
					return m, m.FetchID(item.item.ID)
				}
			case "n", "N":
				m.Page += 1
				m.State = StateLoading
				name := m.NameInput.Value()
				return m, m.FetchName(name)
			case "p", "P":
				m.Page -= 1
				m.State = StateLoading
				name := m.NameInput.Value()
				return m, m.FetchName(name)
			case "m":
				m.State = StateMainMenu
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
				m.State = StateMainMenu
			}
		case StateError:
			m.State = StateMainMenu
		case StateMainMenu:
			switch msg.Type {
			case tea.KeyEnter:
				if item, ok := m.Menu.SelectedItem().(menuItem); ok {
					switch item.State {
					case StateHistory:
						m.State = StateLoading
						return m, m.FetchHistory()
					default:
						m.State = item.State
					}
				}
			}
			m.Menu, cmd = m.Menu.Update(msg)
		case StateTraceNRICInput:
			switch msg.Type {
			case tea.KeyEnter:
				nric := m.NRICInput.Value()
				m.State = StateLoading
				return m, m.FetchNRIC(nric)
			}
			m.NRICInput, cmd = m.NRICInput.Update(msg)
		case StateHistory:
			switch msg.String() {
			case "n", "N":
				m.State = StateLoading
				m.Page += 1
				return m, m.FetchHistory()
			case "p", "P":
				m.State = StateLoading
				m.Page -= 1
				return m, m.FetchHistory()
			case "m", "M":
				m.State = StateMainMenu
			}
			m.List, cmd = m.List.Update(msg)
		}
	case tea.WindowSizeMsg:
		h, v := m.Doc.GetFrameSize()
		m.Width = msg.Width - h
		m.Height = msg.Height - v
		switch m.State {
		case StateTraceNameResults:
			m.List.SetSize(m.Width, m.Height)
		case StateMainMenu:
			m.Menu.SetSize(m.Width, m.Height)
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

			listKeys := newListKeyMap()
			m.List = list.New(items, list.NewDefaultDelegate(), m.Width, m.Height)
			m.List.Title = "Trace Results"
			m.List.AdditionalShortHelpKeys = func() []key.Binding {
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
	case api.HistoryMsg:
		m.State = StateHistory
		var items []list.Item

		for _, item := range msg.Data {
			items = append(items, historyItem{
				item: item,
			})
		}

		listKeys := newListKeyMap()
		m.List = list.New(items, list.NewDefaultDelegate(), m.Width, m.Height)
		m.List.Title = "Search History"
		m.List.FullHelp()
		m.List.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{
				listKeys.toggleNextPage,
				listKeys.togglePreviousPage,
			}
		}
	}
	return m, cmd
}

func (m MainModel) FetchName(name string) tea.Cmd {
	return func() tea.Msg {
		items, err := m.APIClient.TraceName(name, m.Page)
		return api.TraceNameMsg{
			Items: items,
			Err:   err,
		}
	}
}

func (m MainModel) FetchID(id int) tea.Cmd {
	return func() tea.Msg {
		details, err := m.APIClient.TraceDetail(id)
		return api.TraceDetailsMsg{
			Details: details,
			Err:     err,
		}
	}
}

func (m MainModel) FetchRelations(id int) tea.Cmd {
	return func() tea.Msg {
		details, relations, err := m.APIClient.TraceRelations(id, m.Offset)
		return api.TraceRelationsMsg{
			Details:   details,
			Relations: relations,
			Err:       err,
		}
	}
}

func (m MainModel) FetchNRIC(nric string) tea.Cmd {
	return func() tea.Msg {
		details, err := m.APIClient.TraceNRIC(nric)
		return api.TraceDetailsMsg{
			Details: details,
			Err:     err,
		}
	}
}

func (m MainModel) FetchHistory() tea.Cmd {
	return func() tea.Msg {
		data, err := m.APIClient.History(m.Page)
		return api.HistoryMsg{
			Data: data,
			Err:  err,
		}
	}
}
