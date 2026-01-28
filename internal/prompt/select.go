package prompt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 10

var (
	titleStyle        = lipgloss.NewStyle().Bold(true)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return "" }

/* ---------- delegate ---------- */

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

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

	_, _ = fmt.Fprint(w, fn(str))
}

/* ---------- model ---------- */

type selectModel struct {
	list     list.Model
	choice   int
	quitting bool
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			m.choice = m.list.Index()
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectModel) View() string {
	return "\n" + m.list.View()
}

/* ---------- public API ---------- */

func Select(title string, options []string) (int, error) {
	items := make([]list.Item, len(options))
	for i, o := range options {
		items[i] = item(o)
	}

	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Title = title
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowFilter(false)

	p := tea.NewProgram(
		selectModel{list: l, choice: -1},
		tea.WithOutput(os.Stdout),
	)

	model, err := p.Run()
	if err != nil {
		return -1, err
	}

	m := model.(selectModel)
	if m.quitting {
		return -1, fmt.Errorf("selection cancelled")
	}

	return m.choice, nil
}
