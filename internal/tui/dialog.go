package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dialogModel struct {
	connectionName string
	focusYes       bool
	confirmed      bool
	cancelled      bool
}

func newDialogModel(connectionName string) dialogModel {
	return dialogModel{
		connectionName: connectionName,
		focusYes:       false, // Default to "No"
	}
}

func (m dialogModel) Update(msg tea.Msg) (dialogModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "tab", "shift+tab":
			m.focusYes = !m.focusYes

		case "enter", "y":
			if m.focusYes || msg.String() == "y" {
				m.confirmed = true
			} else {
				m.cancelled = true
			}

		case "n", "esc":
			m.cancelled = true
		}
	}

	return m, nil
}

func (m dialogModel) View() string {
	var s strings.Builder

	// Warning title
	warningTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor).
		Render("⚠️  Delete Connection")
	s.WriteString(warningTitle)
	s.WriteString("\n\n")

	// Question
	question := fmt.Sprintf("Are you sure you want to delete connection '%s'?",
		lipgloss.NewStyle().Bold(true).Foreground(dangerColor).Render(m.connectionName))
	s.WriteString(question)
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("This action cannot be undone."))
	s.WriteString("\n\n")

	// Buttons
	yesButton := " Yes "
	noButton := " No "

	if m.focusYes {
		yesButton = activeButtonStyle.Copy().Background(dangerColor).Render(yesButton)
		noButton = inactiveButtonStyle.Render(noButton)
	} else {
		yesButton = inactiveButtonStyle.Render(yesButton)
		noButton = activeButtonStyle.Render(noButton)
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesButton, noButton)
	s.WriteString(buttons)
	s.WriteString("\n\n")

	// Help
	help := "←/→: Navigate • Enter: Confirm • Esc/n: Cancel"
	s.WriteString(helpStyle.Render(help))

	return dialogBoxStyle.Render(s.String())
}
