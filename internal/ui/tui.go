package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("76")).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	boxStyle      = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2)
)

// --- Commit Message Picker TUI ---

type commitPickerModel struct {
	messages []string
	cursor   int
	chosen   int
	done     bool
	quitting bool
}

type commitPickerResult struct {
	chosen int
}

func NewCommitPicker(messages []string) commitPickerModel {
	return commitPickerModel{
		messages: messages,
		chosen:   -1,
	}
}

func (m commitPickerModel) Init() tea.Cmd {
	return nil
}

func (m commitPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.messages)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = m.cursor
			m.done = true
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m commitPickerModel) View() string {
	if m.done {
		return successStyle.Render("Selected: ") + m.messages[m.chosen] + "\n"
	}
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Pick a commit message:") + "\n\n")

	for i, msg := range m.messages {
		cursor := "  "
		style := normalStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedStyle
		}

		// Indent multi-line messages
		lines := strings.Split(msg, "\n")
		b.WriteString(cursor + style.Render(lines[0]) + "\n")
		for _, line := range lines[1:] {
			b.WriteString("    " + dimStyle.Render(line) + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(dimStyle.Render("↑/↓ navigate • enter select • q quit"))
	return b.String()
}

// Result returns the chosen index, or -1 if cancelled.
func (m commitPickerModel) Result() int {
	if m.quitting {
		return -1
	}
	return m.chosen
}

// RunCommitPicker runs the TUI picker and returns the chosen index.
func RunCommitPicker(messages []string) (int, error) {
	model := NewCommitPicker(messages)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return -1, err
	}
	return finalModel.(commitPickerModel).Result(), nil
}

// --- Streaming Output TUI ---

type streamModel struct {
	content  strings.Builder
	done     bool
	quitting bool
	title    string
}

type streamTokenMsg string
type streamDoneMsg struct{}

func NewStreamModel(title string) streamModel {
	return streamModel{title: title}
}

func (m streamModel) Init() tea.Cmd {
	return nil
}

func (m streamModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case streamTokenMsg:
		m.content.WriteString(string(msg))
		return m, nil
	case streamDoneMsg:
		m.done = true
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m streamModel) View() string {
	var b strings.Builder
	if m.title != "" {
		b.WriteString(titleStyle.Render(m.title) + "\n\n")
	}

	content := m.content.String()
	if content == "" && !m.done {
		b.WriteString(dimStyle.Render("Generating..."))
	} else {
		b.WriteString(boxStyle.Render(content))
	}

	if m.done {
		b.WriteString("\n")
	}

	return b.String()
}

// --- Diff Preview ---

// FormatDiffPreview formats a diff for display with colors.
func FormatDiffPreview(diff string) string {
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)

	var b strings.Builder
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- "):
			b.WriteString(headerStyle.Render(line) + "\n")
		case strings.HasPrefix(line, "+"):
			b.WriteString(addStyle.Render(line) + "\n")
		case strings.HasPrefix(line, "-"):
			b.WriteString(delStyle.Render(line) + "\n")
		case strings.HasPrefix(line, "@@"):
			b.WriteString(dimStyle.Render(line) + "\n")
		case strings.HasPrefix(line, "diff "):
			b.WriteString(headerStyle.Render(line) + "\n")
		default:
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

// FormatCostInfo formats token count and cost estimate for display.
func FormatCostInfo(tokenCount int, costEstimate string) string {
	if costEstimate == "" {
		return dimStyle.Render(fmt.Sprintf("(%d tokens)", tokenCount))
	}
	return dimStyle.Render(fmt.Sprintf("(%s)", costEstimate))
}
