package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type todoItem struct {
	Name      string `json:"Name"`
	Completed bool   `json:"Completed"`
}
type allTodos struct {
	PageName string
	Todos    []todoItem
}
type model struct {
	todos           []allTodos
	cursor          int
	pageNum         int
	selected        allTodos
	textInput       textinput.Model
	currentPageName textinput.Model
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	//Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	//Width(22)
	Align(lipgloss.Left)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	ti := textinput.New()
	title := textinput.New()
	ti.Placeholder = "new todo"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	initData := readFromFile()
	title.Placeholder = initData[0].PageName

	title.CharLimit = 156
	title.Width = 20
	return model{
		todos:           initData,
		pageNum:         0,
		selected:        initData[0],
		cursor:          0,
		textInput:       ti,
		currentPageName: title,
	}

}

func readFromFile() []allTodos {
	file, err := os.OpenFile("tasks.json", os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	er(err)
	defer file.Close()
	fileCont, err := io.ReadAll(file)
	er(err)
	var data []allTodos
	err = json.Unmarshal(fileCont, &data)
	er(err)

	return data
}

func writeToFile(todos []allTodos) {

	//turn into json string
	data, err := json.Marshal(todos)

	er(err)
	err = os.WriteFile("tasks.json", data, os.FileMode(os.O_WRONLY))

	er(err)
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
		case "right":
			if m.pageNum < len(m.todos)-1 {
				m.pageNum += 1
				m.selected = m.todos[m.pageNum]
				m.currentPageName.SetValue(m.todos[m.pageNum].PageName)

			} else if m.pageNum == len(m.todos)-1 {
				m.pageNum = 0
				m.selected = m.todos[m.pageNum]
				m.currentPageName.SetValue(m.todos[m.pageNum].PageName)

			}

		case "left":
			if m.pageNum > 0 {
				m.pageNum -= 1
				m.selected = m.todos[m.pageNum]
				m.currentPageName.SetValue(m.todos[m.pageNum].PageName)
			} else if m.pageNum == 0 {
				m.pageNum = len(m.todos) - 1
				m.selected = m.todos[m.pageNum]
				m.currentPageName.SetValue(m.todos[m.pageNum].PageName)

			}
		case "ctrl+c", "esc":
			if m.cursor != len(m.selected.Todos) {
				return m, tea.Quit
			}
		case "q":
			if m.cursor != len(m.selected.Todos) || m.cursor != -1 {
				return m, tea.Quit
			}
		case "up":
			if m.cursor > -1 {
				m.cursor--
				fmt.Print(m.cursor)

			}
			if m.cursor == -1 {
				m.currentPageName.Focus()
			}

		case "down":
			if m.cursor < len(m.selected.Todos) {
				m.cursor++
				m.currentPageName.Blur()
			}
		case "enter":
			if m.cursor == len(m.selected.Todos) {
				m.selected.Todos = append(m.selected.Todos, todoItem{
					Name: m.textInput.Value(),
				})
				m.todos = append(m.todos, m.selected)
				m.textInput.Reset()
				m.todos[m.cursor] = m.selected
			} else if m.cursor > 0 {
				m.selected.Todos[m.cursor].Completed = !(m.selected.Todos[m.cursor].Completed)
				m.todos[m.cursor] = m.selected
			} else if m.cursor == -1 {
				m.selected.PageName = m.currentPageName.Value()
				m.currentPageName.SetValue(m.selected.PageName)

			}
		case "backspace":
			//m.todos[m.cursor].Name = "hello"
			//make SURE that the cursor is on a todo item
			if m.cursor < len(m.selected.Todos) && m.cursor != -1 {
				//we must reslice

				m.selected.Todos = slices.Delete(m.selected.Todos, m.cursor, m.cursor+1)
				//run the below everytime the todos are updates
				m.todos[m.pageNum] = m.selected
				//i hope this works...
			}

		}

	}
	if m.cursor == len(m.selected.Todos) {
		m.textInput, _ = m.textInput.Update(msg)
	}
	if m.cursor == -1 {
		m.currentPageName, _ = m.currentPageName.Update(msg)
		m.todos[m.pageNum].PageName = m.currentPageName.Value()
	}

	defer writeToFile(m.todos)
	return m, nil

}

func (m model) View() string {
	s := fmt.Sprintf("%s \n", m.currentPageName.View()) //hopefully we add todo page name here soon!
	//need another struct for all the placeholder values
	m.currentPageName.SetValue(m.todos[0].PageName)
	for i, todo := range m.selected.Todos {
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
	if m.cursor == len(m.selected.Todos) {
		s += fmt.Sprintf(m.textInput.View())
	}
	s = style.Render(s)
	return s

}
