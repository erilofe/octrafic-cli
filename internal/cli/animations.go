package cli

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// generateGradientText creates a fill animation where colors fill the text letter by letter
func generateGradientText(text string, frame int) string {
	// Blue gradient colors for animation
	colors := []lipgloss.Color{
		Theme.PrimaryDark,   // Blue
		Theme.Primary,       // Lighter blue
		Theme.PrimaryStrong, // Darker blue
	}

	textLength := len([]rune(text))
	if textLength == 0 {
		return text
	}

	// Determine which color is currently filling and how far it's progressed
	currentColorIndex := (frame / textLength) % len(colors)
	fillProgress := frame % textLength

	// Previous color (the one that filled before current)
	previousColorIndex := (currentColorIndex - 1 + len(colors)) % len(colors)

	currentColor := string(colors[currentColorIndex])
	previousColor := string(colors[previousColorIndex])

	result := ""
	i := 0
	for _, char := range text {
		var color string
		if i <= fillProgress {
			// This letter is being filled by current color
			color = currentColor
		} else {
			// This letter still has previous color
			color = previousColor
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		result += style.Render(string(char))
		i++
	}

	return result
}

// animationTick sends a tick message for animation
func animationTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return animationTickMsg(t)
	})
}
