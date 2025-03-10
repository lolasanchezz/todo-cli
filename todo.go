package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type todoItem struct {
	Name      string
	Completed bool
}

type model struct {
	todos     []todoItem
	cursor    int
	selected  map[int]struct{}
	textInput textinput.Model
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	//Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22).
	Align(lipgloss.Center)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Print("hii")
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "new todo"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{
		todos:     readFromFile(),
		selected:  make(map[int]struct{}),
		cursor:    0,
		textInput: ti,
	}

}

func readFromFile() []todoItem {
	file, err := os.OpenFile("tasks.json", os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	er(err)
	defer file.Close()
	fileCont, err := io.ReadAll(file)
	er(err)
	var data []todoItem
	err = json.Unmarshal(fileCont, &data)
	er(err)

	return data
}

func (m model) Init() tea.Cmd {
	return nil
}

func er(Error error) {
	if Error != nil {
		log.Fatal(Error)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//qvar cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.todos) {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.todos) {
				m.todos = append(m.todos, todoItem{Name: m.textInput.Value()})
			}
			m.todos[m.cursor].Completed = !(m.todos[m.cursor].Completed)
		case "backspace":
			//m.todos[m.cursor].Name = "hello"
			if m.cursor < len(m.todos) {
				m.todos[m.cursor] = m.todos[len(m.todos)-1]
				m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
			}
		}

	}
	if m.cursor == len(m.todos) {
		m.textInput, _ = m.textInput.Update(msg)
	}

	return m, nil

}

func (m model) View() string {
	s := "todo \n" //hopefully we add todo page name here soon!
	//need another struct for all the placeholder values

	for i, todo := range m.todos {
		cursor := ""
		completed := ""
		if i == m.cursor {
			cursor = ">"
		}
		if todo.Completed {
			completed = "X"
		}
		s += fmt.Sprintf("%s [%s] %s \n", cursor, completed, todo.Name)
	}
	if m.cursor == len(m.todos) {
		s += fmt.Sprintf(m.textInput.View())
	}
	s = style.Render(s)
	return s

}
