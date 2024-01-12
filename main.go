package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("192"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(2)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// global var to store selected options
var selected map[string]struct{}

var options = []list.Item{
	item("Find unique zips"),
	item("Find unique networks"),
	item("Find unique specialties"),
}

type model struct {
	list     list.Model
	selected map[string]struct{}
	quitting bool
}

func initialModel() model {

	items := options

	const defaultWidth = 25

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the options."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	return model{
		list:     l,
		selected: make(map[string]struct{}),
		quitting: false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Vericred Helper")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter", " ":

			i, ok := m.list.SelectedItem().(item)
			if ok {
				val := string(i)
				val = strings.ReplaceAll(val, "[x]", " ")
				val = strings.Trim(val, " ")

				_, exists := m.selected[val]
				if exists {
					delete(m.selected, val)
				} else {
					m.selected[val] = struct{}{}
				}

			}

		case "p", "P":
			if len(m.selected) > 0 {
				selected = m.selected
				return m, tea.Quit
			}

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "\n\n"

	for i, val := range m.list.Items() {
		x, _ := val.(item)
		temp := string(x)
		temp = strings.ReplaceAll(temp, "[x]", " ")
		temp = strings.Trim(temp, " ")

		_, ok := m.selected[temp]
		if ok {
			m.list.SetItem(i, list.Item(item("[x] "+temp)))
		} else {
			m.list.SetItem(i, list.Item(item(temp)))

		}
	}

	return s + m.list.View() + "\n\nPress p or P to proceed"
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

	helpStyle.SetString("P - to proceed with the selections\n")

	initModel := initialModel()

	p := tea.NewProgram(initModel, tea.WithOutput(os.Stdin))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// buffered channel, only 3 max inputs is possible
	values := make(chan string, len(selected))

	// []string{"Find unique zips", "Find unique networks", "Find unique specialties"},
	// the selected choices
	for i, _ := range selected {
		switch i {
		case "Find unique specialties":
			go searchUniqueSpecialties(filename, values)
		case "Find unique networks":
			go searchUniqueNetworks(filename, values)
		case "Find unique zips":
			go searchUniqueZips(filename, values)
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
