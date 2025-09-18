package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Width(60).Align(lipgloss.Center).PaddingBottom(5)
	labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#87CEEB"))
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	boxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FF69B4")).Margin(1, 1).Padding(1, 2)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(60).AlignHorizontal(lipgloss.Center)
	inputStyle = lipgloss.NewStyle().Align(lipgloss.Center)
)

func (m MainModel) View() string {
	s := strings.Builder{}

	switch m.State {
	case StateTokenInput:
		s.WriteString(m.Logo)
		s.WriteString(inputStyle.Render(m.TokenInput.View()))
		s.WriteString(helpStyle.Render("\nPress enter to continue"))
	case StateTraceNameInput:
		s.WriteString(m.Logo)
		s.WriteString(inputStyle.Render(m.NameInput.View()))
	case StateTraceNameResults:
		return m.Doc.Render(m.List.View())
	case StateTraceDetails:
		details := m.UserDetails
		s.WriteString(m.Logo)

		fields := []struct {
			Label string
			Value string
		}{
			{"Name", details.Name},
			{"Mykad", details.Mykad},
			{"Address", details.Address},
			{"Gender", details.Gender},
			{"Mobile", details.Mobile},
			{"Phone", details.Phone},
			{"Race", details.Race},
			{"Religion", details.Religion},
			{"Income", details.Income},
		}

		var rows []string
		for _, f := range fields {
			if f.Value == "" {
				continue
			}

			rows = append(rows, fmt.Sprintf("%s: %s", labelStyle.Render(f.Label), valueStyle.Render(f.Value)))
		}

		if m.Relations != nil {
			rows = append(rows, fmt.Sprintf("%s: %s", labelStyle.Render("Relation"), valueStyle.Render(m.Relations.Relation)))
		}

		data := strings.Join(rows, "\n")
		s.WriteString(boxStyle.Render(data))
		s.WriteString(helpStyle.Render("\nPress [n] to view next relationships, [p] for previous, and any other key to return to main menu."))
	case StateLoading:
		s.WriteString(fmt.Sprintf("%s Loading...", m.Spinner.View()))
	case StateError:
		s.WriteString(errorStyle.Render("Error: " + m.ErrorMessage))
		s.WriteString(helpStyle.Render("\nPress any key to return to main menu"))
	case StateMainMenu:
		return m.Doc.Render(m.Menu.View())
	case StateTraceNRICInput:
		s.WriteString(m.Logo)
		s.WriteString(inputStyle.Render(m.NRICInput.View()))
	case StateHistory:
		return m.Doc.Render(m.List.View())
	default:
		s.WriteString(fmt.Sprintf("State: %+v", m.State))
	}
	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center,
		lipgloss.Center,
		s.String(),
	)
}
