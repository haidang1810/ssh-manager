package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sm/internal/config"
	"sm/internal/models"
)

type listModel struct {
	table        table.Model
	connections  []connectionItem
	cursor       int
	width        int
	height       int
	shouldAdd    bool
	shouldEdit   bool
	shouldDelete bool
}

func newListModel(connections []connectionItem) listModel {
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "User", Width: 15},
		{Title: "Host", Width: 25},
		{Title: "Port", Width: 6},
		{Title: "Auth", Width: 10},
	}

	rows := makeRows(connections)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("62"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("62")).
		Bold(true)

	t.SetStyles(s)

	return listModel{
		table:       t,
		connections: connections,
	}
}

func makeRows(connections []connectionItem) []table.Row {
	rows := make([]table.Row, len(connections))
	for i, conn := range connections {
		auth := "Password"
		if conn.KeyPath != "" {
			auth = "SSH Key"
		}
		rows[i] = table.Row{
			fmt.Sprintf("%d", conn.ID),
			conn.Name,
			conn.User,
			conn.Host,
			fmt.Sprintf("%d", conn.Port),
			auth,
		}
	}
	return rows
}

func (m *listModel) updateConnections(connections []connectionItem) {
	m.connections = connections
	m.table.SetRows(makeRows(connections))
	if m.cursor >= len(connections) && len(connections) > 0 {
		m.cursor = len(connections) - 1
		m.table.SetCursor(m.cursor)
	}
}

func (m *listModel) setSize(width, height int) {
	m.width = width
	m.height = height

	// Adjust table height
	tableHeight := height - 10 // Leave space for title and help
	if tableHeight < 5 {
		tableHeight = 5
	}
	m.table.SetHeight(tableHeight)

	// Adjust column widths based on available width
	if width < 100 {
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Name", Width: 15},
			{Title: "Host", Width: 20},
			{Title: "Port", Width: 5},
		}
		m.table.SetColumns(columns)
	}
}

func (m listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Connect to selected server
			if len(m.connections) > 0 {
				m.cursor = m.table.Cursor()
				if m.cursor < len(m.connections) {
					// Find the connection in config
					connName := m.connections[m.cursor].Name
					return m, func() tea.Msg {
						// This will be handled by the parent model
						// We'll use a custom message
						cfg, err := loadConfig()
						if err != nil {
							return errorMsg{err: err}
						}
						conn, exists := cfg.Connections[connName]
						if !exists {
							return errorMsg{err: fmt.Errorf("connection not found")}
						}
						return connectMsg{Connection: conn}
					}
				}
			}

		case "a", "n":
			// Add new connection
			m.shouldAdd = true
			return m, nil

		case "e":
			// Edit selected connection
			if len(m.connections) > 0 {
				m.cursor = m.table.Cursor()
				m.shouldEdit = true
			}
			return m, nil

		case "d", "x":
			// Delete selected connection
			if len(m.connections) > 0 {
				m.cursor = m.table.Cursor()
				m.shouldDelete = true
			}
			return m, nil

		case "r":
			// Refresh connection list
			return m, func() tea.Msg {
				return successMsg{}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	m.cursor = m.table.Cursor()

	return m, cmd
}

func (m listModel) View() string {
	var s strings.Builder

	// Title
	title := titleStyle.Render("🔐 SSH Connection Manager")
	s.WriteString(title)
	s.WriteString("\n\n")

	// Connection count
	countText := fmt.Sprintf("Total Connections: %d", len(m.connections))
	s.WriteString(headerStyle.Render(countText))
	s.WriteString("\n\n")

	// Table
	if len(m.connections) == 0 {
		emptyText := "No connections found. Press 'a' to add one."
		s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render(emptyText))
	} else {
		s.WriteString(m.table.View())
	}

	s.WriteString("\n\n")

	// Help text
	help := []string{
		"↑/↓: Navigate",
		"Enter: Connect",
		"a: Add",
		"e: Edit",
		"d: Delete",
		"r: Refresh",
		"q: Quit",
	}
	helpText := strings.Join(help, " • ")
	s.WriteString(helpStyle.Render(helpText))

	return borderStyle.Render(s.String())
}

// Helper function to load config
func loadConfig() (*models.AppConfig, error) {
	return config.GetConfig()
}
