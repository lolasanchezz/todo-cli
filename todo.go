package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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
	initialModel()

}

func initialModel() model {

	return model{
		todos: readFromFile(),
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

func er(Error error) {
	if Error != nil {
		log.Fatal(Error)
	}
}
