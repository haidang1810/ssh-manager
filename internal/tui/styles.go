package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("62")  // Bright blue
	secondaryColor = lipgloss.Color("170") // Purple
	successColor   = lipgloss.Color("42")  // Green
	dangerColor    = lipgloss.Color("196") // Red
	warningColor   = lipgloss.Color("214") // Orange
	mutedColor     = lipgloss.Color("241") // Gray

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			MarginBottom(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Status bar style
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(primaryColor).
			Padding(0, 1)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	// Success style
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(primaryColor).
			Bold(true)

	// Form label style
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Width(15).
			Align(lipgloss.Right)

	// Form input style
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// Border style
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	// Dialog style
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			Width(50)

	// Button styles
	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(primaryColor).
				Padding(0, 3).
				MarginRight(2)

	inactiveButtonStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Background(lipgloss.Color("236")).
				Padding(0, 3).
				MarginRight(2)
)
