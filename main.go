package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

// global var
var exit bool = false

// selected checkboxes
var selected map[int]struct{}

func initialModel() model {
	return model{
		choices: []string{"Find unique zips", "Find unique networks", "Find unique specialties"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Grocery List")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "p", "P":
			tea.Printf("Hello jssjbs\n")
			if len(m.selected) > 0 {
				selected = m.selected
				return m, tea.Quit
			}
		}

	}

	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress p or P to proceed.\n"
	s += "\nPress q to quit.\n"

	return s
}

func main() {
	// the first arguement is the filename
	filename := os.Args[1]

	// check if the provided filename exists
	exists, err := checkIfFileExists(filename)

	// check for error
	if err != nil && !exists {
		log.Fatal(err)
	}

	initModel := initialModel()

	p := tea.NewProgram(initModel, tea.WithOutput(os.Stdin))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// check if exit was called for
	if exit {
		return
	}

	// buffered channel, only 3 max inputs is possible
	values := make(chan string, len(selected))

	// []string{"Find unique zips", "Find unique networks", "Find unique specialties"},
	// the selected choices
	for i, _ := range selected {
		switch i {
		case 0:
			go searchUniqueZips(filename, values)
		case 1:
			go searchUniqueNetworks(filename, values)
		case 2:
			go searchUniqueSpecialties(filename, values)
		}
	}

	// Since it is a buffered channel
	for i := 0; i < len(selected); i++ {
		fmt.Println(<-values)
	}

}

func checkIfFileExists(filename string) (bool, error) {

	file, err := os.Open(filename)

	if err != nil {
		return false, err
	}

	file.Close()
	return true, nil
}