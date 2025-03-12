package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	return nil
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
	s := "Which task have you done today?\n"

	for i, task := range m.Tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.SelectedTask[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s \n", cursor, checked, task)
	}

	if m.TextInput.Value() != "" {
		return fmt.Sprintf(
			"Which task would you like to add?\nctrl + a to save the task and ctrl + x to cancel\n%s\n", m.TextInput.View())
	}

	s += "Type ctrl+c or q to quit. \n"
	return s
}

func main() {
	prog := tea.NewProgram(initModel())
	_, err := prog.Run()
	if err != nil {
		fmt.Printf("Error while running the prog : %s", err)
		os.Exit(1)
	}
}
