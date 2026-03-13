package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type delegateKeyMap struct {
	choose key.Binding
	delete key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.delete,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.delete,
		},
	}
}

func newItemDelegate() list.ItemDelegate {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle select"),
		),
		delete: key.NewBinding(
			key.WithKeys("d", "enter"),
			key.WithHelp("d/enter", "delete"),
		),
	}
}

func (d *delegateKeyMap) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	checkbox := "[ ]"
	if i.selected {
		checkbox = "[x]"
	}

	str := fmt.Sprintf("%s %s", checkbox, i.title)

	if index == m.Index() {
		_, _ = fmt.Fprint(w, selectedItemStyle.Render("> "+str))
	} else {
		_, _ = fmt.Fprint(w, itemStyle.Render(str))
	}
}

func (d *delegateKeyMap) Height() int {
	return 1
}

func (d *delegateKeyMap) Spacing() int {
	return 0
}

func (d *delegateKeyMap) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
