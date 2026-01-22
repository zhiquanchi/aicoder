package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
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

type model struct {
	list     list.Model
	choice   string
	quitting bool
	view     string // "menu", "tools", "config", "projects"
}

func (m model) Init() tea.Cmd {
	return nil
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

		case "q":
			if m.view != "menu" {
				// Go back to main menu
				m.view = "menu"
				m.list.Title = "🤖 AICoder TUI - Main Menu"
				items := []list.Item{
					item("View Tool Status"),
					item("Configure API Keys"),
					item("Manage Projects"),
					item("Launch AI Tool"),
					item("Exit"),
				}
				m.list.SetItems(items)
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
				
				switch m.view {
				case "menu":
					switch m.choice {
					case "View Tool Status":
						m.view = "tools"
						m.list.Title = "📊 Tool Status (Press 'q' to go back)"
						items := []list.Item{
							item("Claude - Checking..."),
							item("Gemini - Checking..."),
							item("Codex - Checking..."),
							item("OpenCode - Checking..."),
							item("CodeBuddy - Checking..."),
							item("Qoder - Checking..."),
						}
						m.list.SetItems(items)
						return m, nil
					case "Configure API Keys":
						m.view = "config"
						m.list.Title = "🔑 Configure API Keys (Press 'q' to go back)"
						items := []list.Item{
							item("Configure Claude API Key"),
							item("Configure Gemini API Key"),
							item("Configure Codex API Key"),
							item("Configure OpenCode API Key"),
							item("Configure CodeBuddy API Key"),
							item("Configure Qoder API Key"),
						}
						m.list.SetItems(items)
						return m, nil
					case "Manage Projects":
						m.view = "projects"
						m.list.Title = "📂 Manage Projects (Press 'q' to go back)"
						items := []list.Item{
							item("List Projects"),
							item("Add New Project"),
							item("Set Active Project"),
						}
						m.list.SetItems(items)
						return m, nil
					case "Launch AI Tool":
						m.view = "launch"
						m.list.Title = "🚀 Launch AI Tool (Press 'q' to go back)"
						items := []list.Item{
							item("Launch Claude"),
							item("Launch Gemini"),
							item("Launch Codex"),
							item("Launch OpenCode"),
							item("Launch CodeBuddy"),
							item("Launch Qoder"),
						}
						m.list.SetItems(items)
						return m, nil
					case "Exit":
						m.quitting = true
						return m, tea.Quit
					}
				case "tools", "config", "projects", "launch":
					// For now, just show a message
					// This will be expanded with actual functionality
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" && m.view != "menu" {
		return quitTextStyle.Render(fmt.Sprintf("Selected: %s\n\nThis feature is under construction. Press Ctrl+C to exit.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Thanks for using AICoder TUI! Goodbye! 👋")
	}
	return "\n" + m.list.View()
}

// RunTUI starts the TUI application
func RunTUI() error {
	items := []list.Item{
		item("View Tool Status"),
		item("Configure API Keys"),
		item("Manage Projects"),
		item("Launch AI Tool"),
		item("Exit"),
	}

	const defaultWidth = 80

	l := list.New(items, itemDelegate{}, defaultWidth, 14)
	l.Title = "🤖 AICoder TUI - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, view: "menu"}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}

	return nil
}
