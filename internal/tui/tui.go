package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/bansikah22/kswp/internal/cleaner"
	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/pkg/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type model struct {
	list     list.Model
	items    []item
	client   kubernetes.Client
	deleting bool
}

type item struct {
	title    string
	desc     string
	resource models.Resource
	selected bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewModel(resources []models.Resource, client kubernetes.Client) model {
	var items []item
	for _, res := range resources {
		items = append(items, item{
			title:    fmt.Sprintf("%s/%s", res.Namespace, res.Name),
			desc:     fmt.Sprintf("Kind: %s, Reason: %s, Age: %s", res.Kind, res.Reason, res.Age.Round(time.Second).String()),
			resource: res,
		})
	}
	listItems := make([]list.Item, len(items))
	for i, itm := range items {
		listItems[i] = itm
	}
	delegate := newItemDelegate()
	resourceList := list.New(listItems, delegate, 0, 0)
	resourceList.Title = "Unused Resources"
	resourceList.Styles.Title = titleStyle
	return model{
		list:   resourceList,
		items:  items,
		client: client,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("space"))):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				i.selected = !i.selected
				m.list.SetItem(m.list.Index(), i)
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("d", "enter"))):
			m.deleting = true
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
			if m.deleting {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					err := cleaner.DeleteResource(m.client.Clientset(), i.resource)
					if err != nil {
						// handle error
					}
					m.list.RemoveItem(m.list.Index())
				}
				m.deleting = false
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			if m.deleting {
				m.deleting = false
			}
		}
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.deleting {
		return appStyle.Render(statusMessageStyle.Render("Are you sure you want to delete this resource? (y/n)"))
	}
	return appStyle.Render(m.list.View())
}

func StartTUI(resources []models.Resource, client kubernetes.Client) {
	m := NewModel(resources, client)
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
