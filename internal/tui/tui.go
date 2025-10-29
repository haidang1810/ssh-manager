package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"sm/internal/models"
	"sm/internal/ssh"
)

// Launch starts the TUI application
func Launch() error {
	m, err := InitialModel()
	if err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	// Check if we need to connect to a server after exiting
	if model, ok := finalModel.(Model); ok {
		if model.connectOnExit != nil {
			fmt.Printf("Connecting to %s (%s@%s)...\n",
				model.connectOnExit.Name,
				model.connectOnExit.User,
				model.connectOnExit.Host)

			if err := ssh.Connect(model.connectOnExit); err != nil {
				return fmt.Errorf("ssh connection failed: %w", err)
			}

			fmt.Println("Connection closed.")
		}
	}

	return nil
}

// LaunchWithConnection starts the TUI and immediately attempts to connect to the specified connection
func LaunchWithConnection(conn *models.Connection) error {
	fmt.Printf("Connecting to %s (%s@%s)...\n", conn.Name, conn.User, conn.Host)

	if err := ssh.Connect(conn); err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}

	fmt.Println("Connection closed.")
	return nil
}
