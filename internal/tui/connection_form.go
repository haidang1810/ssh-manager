package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sm/internal/config"
	"sm/internal/models"
	"sm/internal/utils"
)

type formMode int

const (
	formModeAdd formMode = iota
	formModeEdit
)

type formModel struct {
	mode      formMode
	inputs    []textinput.Model
	focused   int
	cfg       *models.AppConfig
	cancelled bool
	err       error

	// Original connection (for edit mode)
	originalName string
	originalID   int
}

const (
	inputName = iota
	inputHost
	inputUser
	inputPort
	inputKeyPath
	inputPassword
)

func newFormModel(mode formMode, conn connectionItem, cfg *models.AppConfig) formModel {
	m := formModel{
		mode:    mode,
		inputs:  make([]textinput.Model, 6),
		focused: 0,
		cfg:     cfg,
	}

	// Name
	m.inputs[inputName] = textinput.New()
	m.inputs[inputName].Placeholder = "my-server"
	m.inputs[inputName].Focus()
	m.inputs[inputName].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputName].TextStyle = inputStyle

	// Host
	m.inputs[inputHost] = textinput.New()
	m.inputs[inputHost].Placeholder = "example.com or 192.168.1.1"
	m.inputs[inputHost].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputHost].TextStyle = inputStyle

	// User
	m.inputs[inputUser] = textinput.New()
	m.inputs[inputUser].Placeholder = "root"
	m.inputs[inputUser].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputUser].TextStyle = inputStyle

	// Port
	m.inputs[inputPort] = textinput.New()
	m.inputs[inputPort].Placeholder = "22"
	m.inputs[inputPort].CharLimit = 5
	m.inputs[inputPort].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputPort].TextStyle = inputStyle

	// Key Path
	m.inputs[inputKeyPath] = textinput.New()
	m.inputs[inputKeyPath].Placeholder = "~/.ssh/id_rsa (optional)"
	m.inputs[inputKeyPath].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputKeyPath].TextStyle = inputStyle

	// Password
	m.inputs[inputPassword] = textinput.New()
	m.inputs[inputPassword].Placeholder = "password (optional)"
	m.inputs[inputPassword].EchoMode = textinput.EchoPassword
	m.inputs[inputPassword].EchoCharacter = '•'
	m.inputs[inputPassword].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
	m.inputs[inputPassword].TextStyle = inputStyle

	// Fill in values for edit mode
	if mode == formModeEdit {
		m.originalName = conn.Name
		m.originalID = conn.ID
		m.inputs[inputName].SetValue(conn.Name)
		m.inputs[inputHost].SetValue(conn.Host)
		m.inputs[inputUser].SetValue(conn.User)
		m.inputs[inputPort].SetValue(fmt.Sprintf("%d", conn.Port))
		m.inputs[inputKeyPath].SetValue(conn.KeyPath)
		// Don't show password in edit mode
	}

	return m
}

func (m formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, nil

		case "tab", "down":
			m.nextInput()

		case "shift+tab", "up":
			m.prevInput()

		case "enter":
			if m.focused == len(m.inputs)-1 {
				// Last input, submit form
				return m, m.submitForm()
			}
			m.nextInput()
		}
	}

	// Update current input
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)

	return m, cmd
}

func (m *formModel) nextInput() {
	m.inputs[m.focused].Blur()
	m.focused++
	if m.focused >= len(m.inputs) {
		m.focused = 0
	}
	m.inputs[m.focused].Focus()
}

func (m *formModel) prevInput() {
	m.inputs[m.focused].Blur()
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	m.inputs[m.focused].Focus()
}

func (m *formModel) submitForm() tea.Cmd {
	return func() tea.Msg {
		// Validate inputs
		name := strings.TrimSpace(m.inputs[inputName].Value())
		host := strings.TrimSpace(m.inputs[inputHost].Value())
		user := strings.TrimSpace(m.inputs[inputUser].Value())
		portStr := strings.TrimSpace(m.inputs[inputPort].Value())
		keyPath := strings.TrimSpace(m.inputs[inputKeyPath].Value())
		password := m.inputs[inputPassword].Value()

		if name == "" {
			return errorMsg{err: fmt.Errorf("name is required")}
		}
		if host == "" {
			return errorMsg{err: fmt.Errorf("host is required")}
		}
		if user == "" {
			return errorMsg{err: fmt.Errorf("user is required")}
		}

		port := 22
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil || port < 1 || port > 65535 {
				return errorMsg{err: fmt.Errorf("invalid port number")}
			}
		}

		// Check if name already exists (for add mode)
		if m.mode == formModeAdd {
			if _, exists := m.cfg.Connections[name]; exists {
				return errorMsg{err: fmt.Errorf("connection with name '%s' already exists", name)}
			}
		}

		// Encrypt password if provided
		if password != "" {
			encrypted, err := utils.Encrypt(password)
			if err != nil {
				return errorMsg{err: fmt.Errorf("failed to encrypt password: %w", err)}
			}
			password = encrypted
		}

		// Create or update connection
		var conn models.Connection
		if m.mode == formModeAdd {
			conn = models.Connection{
				ID:        m.cfg.NextID,
				Name:      name,
				Host:      host,
				User:      user,
				Port:      port,
				KeyPath:   keyPath,
				Password:  password,
				CreatedAt: time.Now().Unix(),
			}
			m.cfg.Connections[name] = conn
			m.cfg.NextID++
		} else {
			// Edit mode - preserve existing connection data
			existingConn := m.cfg.Connections[m.originalName]

			// If name changed, delete old entry
			if name != m.originalName {
				delete(m.cfg.Connections, m.originalName)
			}

			conn = models.Connection{
				ID:        m.originalID,
				Name:      name,
				Host:      host,
				User:      user,
				Port:      port,
				KeyPath:   keyPath,
				CreatedAt: existingConn.CreatedAt,
				LastUsed:  existingConn.LastUsed,
			}

			// Only update password if a new one was entered
			if password != "" {
				conn.Password = password
			} else {
				conn.Password = existingConn.Password
			}

			m.cfg.Connections[name] = conn
		}

		// Save config
		if err := config.SaveConfig(m.cfg); err != nil {
			return errorMsg{err: fmt.Errorf("failed to save config: %w", err)}
		}

		return successMsg{}
	}
}

func (m formModel) View() string {
	var s strings.Builder

	// Title
	title := "Add New Connection"
	if m.mode == formModeEdit {
		title = "Edit Connection"
	}
	s.WriteString(titleStyle.Render(title))
	s.WriteString("\n\n")

	// Form fields
	fields := []string{
		"Name:",
		"Host:",
		"User:",
		"Port:",
		"SSH Key Path:",
		"Password:",
	}

	for i, field := range fields {
		label := labelStyle.Render(field)
		input := m.inputs[i].View()

		if i == m.focused {
			label = labelStyle.Copy().Foreground(primaryColor).Bold(true).Render(field)
		}

		s.WriteString(fmt.Sprintf("%s %s\n", label, input))
	}

	s.WriteString("\n")

	// Help text
	help := "Tab/↑↓: Navigate • Enter: Submit • Esc: Cancel"
	s.WriteString(helpStyle.Render(help))

	return borderStyle.Render(s.String())
}
