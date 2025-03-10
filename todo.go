package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type todoItem struct {
	Name      string
	Completed bool
}

type model struct {
	todos    []todoItem
	cursor   int
	selected map[int]struct{}
}

func main() {
	fmt.Print("hii")
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {

	return model{
		todos:    readFromFile(),
		selected: make(map[int]struct{}),
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

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
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
			m.todos[m.cursor].Completed = !(m.todos[m.cursor].Completed)
		}

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, completed, todo.Name)
	}
	return s

}
