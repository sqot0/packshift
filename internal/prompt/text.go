package prompt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textModel struct {
	input textinput.Model
	value string
	label string
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m textModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.label,
		m.input.View(),
	) + "\n"
}

func Text(label string, placeholder string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = 20
	ti.CharLimit = 156
	ti.Focus()

	p := tea.NewProgram(textModel{input: ti, label: label}, tea.WithOutput(os.Stdout))
	model, err := p.Run()
	if err != nil {
		return "", err
	}
	if model.(textModel).value == "" {
		return placeholder, nil
	}
	return model.(textModel).value, nil
}
