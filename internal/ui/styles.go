package ui

import "github.com/charmbracelet/lipgloss"

// Color palette - dark neutral base with cyan accent
var (
	ColorBg        = lipgloss.Color("234") // deep background for active items/cards
	ColorPanel     = lipgloss.Color("235") // card background
	ColorPanelAlt  = lipgloss.Color("236") // subtle contrast
	ColorBorder    = lipgloss.Color("239") // hairline divider/border
	ColorMuted     = lipgloss.Color("244") // muted text
	ColorMuted2    = lipgloss.Color("242") // extra muted
	ColorText      = lipgloss.Color("252") // primary text
	ColorCaption   = lipgloss.Color("245") // caption/subtitle
	ColorAccent    = lipgloss.Color("45")  // cyan accent
	ColorAccentAlt = lipgloss.Color("51")  // cyan-light
	ColorWarn      = lipgloss.Color("178") // amber
	ColorError     = lipgloss.Color("203") // soft red
	ColorSuccess   = lipgloss.Color("114") // soft green
	ColorInfo      = lipgloss.Color("110") // blue-cyan
)

// Spacing and layout
const (
	SpaceX = 2
	SpaceY = 1
)

// Header styles
var (
	AppTitle = lipgloss.NewStyle().
		Foreground(ColorText).
		Bold(true)

	AppCaption = lipgloss.NewStyle().
		Foreground(ColorCaption)

	TitleDivider = lipgloss.NewStyle().
		Foreground(ColorBorder)
)

// Card/panel styles
var (
	Card = lipgloss.NewStyle().
		Padding(1, 2).
		Background(ColorPanel).
		Foreground(ColorText).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)

	CardAlt = lipgloss.NewStyle().
		Padding(0, 2).
		Background(ColorPanelAlt).
		Foreground(ColorText).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder)
)

// Menu styles
var (
	MenuPointer = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	MenuItem = lipgloss.NewStyle().
		Foreground(ColorMuted)

	MenuItemActive = lipgloss.NewStyle().
		Foreground(ColorText).
		Bold(true)

	CurrentPointer = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	// Style for current checkpoint (bold green)
	CurrentCheckpointStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)
)

// Status styles
var (
	InfoStyle = lipgloss.NewStyle().
		Foreground(ColorInfo).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	LoadingTextStyle = lipgloss.NewStyle().
		Foreground(ColorAccent)

	HelpStyle = lipgloss.NewStyle().
		Foreground(ColorMuted2)

	DisabledStyle = lipgloss.NewStyle().
		Foreground(ColorBorder).  // Use border color (239) - darker than muted
		Faint(true)               // Make it even more faded

	DisabledReasonStyle = lipgloss.NewStyle().
		Foreground(ColorError)    // Red color for disabled reasons

	Hairline = lipgloss.NewStyle().
		Foreground(ColorBorder)
)