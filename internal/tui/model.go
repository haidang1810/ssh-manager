package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"sm/internal/config"
	"sm/internal/models"
	"sm/internal/ssh"
)

// ViewState represents the current view/screen
type ViewState int

const (
	ViewList ViewState = iota
	ViewAddForm
	ViewEditForm
	ViewDeleteConfirm
)

// Model represents the main TUI model
type Model struct {
	state       ViewState
	config      *models.AppConfig
	connections []connectionItem // Sorted list for display

	// Sub-models
	listModel   listModel
	formModel   formModel
	dialogModel dialogModel

	width  int
	height int
	err    error

	// Connection to connect
	connectOnExit *models.Connection
}

// connectionItem is a helper struct for display
type connectionItem struct {
	ID      int
	Name    string
	Host    string
	User    string
	Port    int
	KeyPath string
}

// InitialModel creates the initial TUI model
func InitialModel() (Model, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load config: %w", err)
	}

	m := Model{
		state:  ViewList,
		config: cfg,
	}

	m.refreshConnections()
	m.listModel = newListModel(m.connections)

	return m, nil
}

// refreshConnections updates the sorted connection list
func (m *Model) refreshConnections() {
	m.connections = make([]connectionItem, 0, len(m.config.Connections))
	for _, conn := range m.config.Connections {
		m.connections = append(m.connections, connectionItem{
			ID:      conn.ID,
			Name:    conn.Name,
			Host:    conn.Host,
			User:    conn.User,
			Port:    conn.Port,
			KeyPath: conn.KeyPath,
		})
	}
	// Sort by ID
	for i := 0; i < len(m.connections)-1; i++ {
		for j := i + 1; j < len(m.connections); j++ {
			if m.connections[i].ID > m.connections[j].ID {
				m.connections[i], m.connections[j] = m.connections[j], m.connections[i]
			}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == ViewList {
				return m, tea.Quit
			}
			// In other views, go back to list
			m.state = ViewList
			return m, nil

		case "esc":
			if m.state != ViewList {
				m.state = ViewList
				m.err = nil
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listModel.setSize(msg.Width, msg.Height)

	case connectMsg:
		// Connect command - save connection and quit
		m.connectOnExit = &msg.Connection
		return m, tea.Quit

	case errorMsg:
		m.err = msg.err
		return m, nil

	case successMsg:
		// Refresh config and go back to list
		cfg, err := config.GetConfig()
		if err != nil {
			m.err = err
		} else {
			m.config = cfg
			m.refreshConnections()
			m.listModel.updateConnections(m.connections)
			m.state = ViewList
			m.err = nil
		}
		return m, nil
	}

	// Delegate to sub-models based on state
	var cmd tea.Cmd
	switch m.state {
	case ViewList:
		m.listModel, cmd = m.listModel.Update(msg)
		if m.listModel.shouldAdd {
			m.listModel.shouldAdd = false
			m.state = ViewAddForm
			m.formModel = newFormModel(formModeAdd, connectionItem{}, m.config)
		} else if m.listModel.shouldEdit {
			m.listModel.shouldEdit = false
			if m.listModel.cursor < len(m.connections) {
				m.state = ViewEditForm
				m.formModel = newFormModel(formModeEdit, m.connections[m.listModel.cursor], m.config)
			}
		} else if m.listModel.shouldDelete {
			m.listModel.shouldDelete = false
			if m.listModel.cursor < len(m.connections) {
				m.state = ViewDeleteConfirm
				m.dialogModel = newDialogModel(m.connections[m.listModel.cursor].Name)
			}
		}

	case ViewAddForm, ViewEditForm:
		m.formModel, cmd = m.formModel.Update(msg)
		if m.formModel.cancelled {
			m.state = ViewList
		}

	case ViewDeleteConfirm:
		m.dialogModel, cmd = m.dialogModel.Update(msg)
		if m.dialogModel.confirmed {
			// Delete connection
			if m.listModel.cursor < len(m.connections) {
				connName := m.connections[m.listModel.cursor].Name
				delete(m.config.Connections, connName)
				if err := config.SaveConfig(m.config); err != nil {
					m.err = fmt.Errorf("failed to delete connection: %w", err)
				} else {
					m.refreshConnections()
					m.listModel.updateConnections(m.connections)
					if m.listModel.cursor >= len(m.connections) && m.listModel.cursor > 0 {
						m.listModel.cursor--
					}
				}
			}
			m.state = ViewList
		} else if m.dialogModel.cancelled {
			m.state = ViewList
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string
	switch m.state {
	case ViewList:
		content = m.listModel.View()
	case ViewAddForm, ViewEditForm:
		content = m.formModel.View()
	case ViewDeleteConfirm:
		content = m.dialogModel.View()
	}

	if m.err != nil {
		content += "\n\n" + errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	return content
}

// Custom messages
type connectMsg struct {
	Connection models.Connection
}

type errorMsg struct {
	err error
}

type successMsg struct{}

// ConnectToServer attempts to connect via SSH
func ConnectToServer(conn models.Connection) tea.Cmd {
	return func() tea.Msg {
		if err := ssh.Connect(&conn); err != nil {
			return errorMsg{err: fmt.Errorf("connection failed: %w", err)}
		}
		return successMsg{}
	}
}
