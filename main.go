package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	TextInput    textinput.Model
	Tasks        []string
	cursor       int
	SelectedTask map[int]struct{}
}

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		Tasks:        []string{"Review du code", "correct the bugs", "help others"},
		SelectedTask: make(map[int]struct{}),
		TextInput:    ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Task manager Tui")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.TextInput.Value() != "" && msg.String() != "ctrl+a" && msg.String() != "ctrl+x" {
			m.TextInput, cmd = m.TextInput.Update(msg)
		} else {

			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "j":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "k":
				if m.cursor < (len(m.Tasks) - 1) {
					m.cursor++
				}
			case "enter", " ":
				if _, ok := m.SelectedTask[m.cursor]; ok {
					delete(m.SelectedTask, m.cursor)
				} else {
					m.SelectedTask[m.cursor] = struct{}{}
				}
			case "ctrl+x":
				m.TextInput.SetValue("")
			case "ctrl+a":
				m.Tasks = append(m.Tasks, m.TextInput.Value())
				m.TextInput.SetValue("")
			default:
				m.TextInput, cmd = m.TextInput.Update(msg)
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	vert := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#04B575")).Padding(0, 1)
	jaune := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ffff00")).Padding(0, 1).Align(lipgloss.Center)
	rouge := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#e9383f")).Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		Margin(1, 0).
		Foreground(lipgloss.Color("205"))

	taskStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("69")).
		Padding(0, 1)

	selectedTaskStyle := taskStyle.Copy().
		Foreground(lipgloss.Color("201"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)

	b.WriteString("\n" + vert.Render(" ") + jaune.Render("*") + rouge.Render(" ") + "\n")
	b.WriteString(titleStyle.Render("Which task have you done today?") + "\n\n")

	for i, task := range m.Tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.SelectedTask[i]; ok {
			checked = "x"
		}

		taskLine := fmt.Sprintf("%s [%s] %s", cursor, checked, task)
		if _, ok := m.SelectedTask[i]; ok {
			b.WriteString(selectedTaskStyle.Render(taskLine) + "\n")
		} else {
			b.WriteString(taskStyle.Render(taskLine) + "\n")
		}
	}

	if m.TextInput.Value() != "" {
		b.WriteString(fmt.Sprintf(
			"\n%s\n%s\n",
			titleStyle.Render("Which task would you like to add?"),
			m.TextInput.View(),
		))
		b.WriteString("\n" + helpStyle.Render("\n• ctrl+a: add task • ctrl+x: clear input") + "\n")
	} else {
		b.WriteString("\n" + titleStyle.Render("Type ctrl+c or q to quit.") + "\n")
	}

	b.WriteString("\n" + helpStyle.Render("Commands: ↑/↓: navigate • space: toggle \n") + "\n")

	return b.String()
}

func main() {
	prog := tea.NewProgram(initModel())
	_, err := prog.Run()
	if err != nil {
		fmt.Printf("Error while running the prog : %s", err)
		os.Exit(1)
	}
}
