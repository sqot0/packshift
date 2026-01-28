package prompt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passwordModel struct {
	input textinput.Model
	value string
	label string
}

func (m passwordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m passwordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.value = m.input.Value()
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m passwordModel) View() string {
	displayedText := ""

	for range len(m.input.Value()) {
		displayedText += "*"
	}

	return fmt.Sprintf(
		"%s\n> %s",
		m.label,
		displayedText,
	) + "\n"
}

func Password(label string) (string, error) {
	ti := textinput.New()
	ti.Width = 20
	ti.CharLimit = 156
	ti.Focus()

	p := tea.NewProgram(passwordModel{input: ti, label: label}, tea.WithOutput(os.Stdout))
	model, err := p.Run()
	if err != nil {
		return "", err
	}
	if model.(passwordModel).value == "" {
		return "", nil
	}
	return model.(passwordModel).value, nil
}
