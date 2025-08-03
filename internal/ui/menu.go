package ui

import (
	"fmt"
	"strings"
	"vibe-check/internal/models"

	"github.com/charmbracelet/lipgloss"
)

// RenderMenu renders the main menu
func RenderMenu(m models.AppModel) string {
	var s strings.Builder

	// Header
	title := lipgloss.JoinHorizontal(lipgloss.Left,
		AppTitle.Render("vibe-check"),
		"  ",
		AppCaption.Render("Git Checkpoint System"),
	)
	divider := TitleDivider.Render(strings.Repeat("─", 40))

	var menu strings.Builder
	
	// Menu items
	for i, choice := range m.MenuChoices {
		prefix := "  "
		itemStyle := MenuItem
		
		// Check if item is disabled
		if m.DisabledMenuItems[i] {
			itemStyle = DisabledStyle // Use very dark style for disabled items - no cursor on disabled
			reason := m.DisabledReasons[i]
			if reason != "" {
				// Show item with red reason
				line := fmt.Sprintf("  %s ", choice)
				reasonLine := DisabledReasonStyle.Render(reason)
				menu.WriteString(itemStyle.Render(line) + reasonLine)
			} else {
				line := fmt.Sprintf("  %s", choice)
				menu.WriteString(itemStyle.Render(line))
			}
		} else {
			if i == m.MenuCursor {
				prefix = MenuPointer.Render("› ")
				itemStyle = MenuItemActive
			}
			line := fmt.Sprintf("%s%s", prefix, choice)
			menu.WriteString(itemStyle.Render(line))
		}
		menu.WriteString("\n")
	}
	
	// Footer
	footer := HelpStyle.Render("↑/↓ navigate • Enter select • q quit")
	dividerLine := Hairline.Render(strings.Repeat("─", 40))
	
	body := strings.TrimRight(menu.String(), "\n") + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n" + TitleDivider.Render(divider) + "\n\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}

// RenderCheckpointCreation renders the checkpoint creation menu
func RenderCheckpointCreation(m models.AppModel) string {
	var s strings.Builder

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoStyle.Render("Create Checkpoint"),
		"  ",
		AppCaption.Render("Choose checkpoint type"),
	)

	var menu strings.Builder
	
	for i, choice := range m.CheckpointOptions {
		prefix := "  "
		itemStyle := MenuItem
		
		if i == m.CheckpointOptionsCursor {
			prefix = MenuPointer.Render("› ")
			itemStyle = MenuItemActive
		}
		
		line := fmt.Sprintf("%s%s", prefix, choice)
		menu.WriteString(itemStyle.Render(line))
		menu.WriteString("\n")
	}
	
	footer := HelpStyle.Render("↑/↓ navigate • Enter select • Esc back")
	dividerLine := Hairline.Render(strings.Repeat("─", 40))
	
	body := strings.TrimRight(menu.String(), "\n") + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}

func RenderFinalizeOptions(m models.AppModel) string {
	var s strings.Builder

	title := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoStyle.Render("Finalize and Push"),
		"  ",
		AppCaption.Render("Choose finalize option"),
	)

	var menu strings.Builder
	
	for i, choice := range m.FinalizeOptions {
		prefix := "  "
		itemStyle := MenuItem
		
		if i == m.FinalizeOptionsCursor {
			prefix = MenuPointer.Render("› ")
			itemStyle = MenuItemActive
		}
		
		line := fmt.Sprintf("%s%s", prefix, choice)
		menu.WriteString(itemStyle.Render(line))
		menu.WriteString("\n")
	}
	
	footer := HelpStyle.Render("↑/↓ navigate • Enter select • Esc back")
	dividerLine := Hairline.Render(strings.Repeat("─", 40))
	
	body := strings.TrimRight(menu.String(), "\n") + "\n" + dividerLine + "\n" + footer
	
	s.WriteString(CardAlt.Render(title) + "\n")
	s.WriteString(Card.Render(body))
	
	return s.String()
}