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
	headerStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).MarginLeft(2).MarginBottom(1)
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

type toolStatusMsg map[string]string

type model struct {
	list         list.Model
	choice       string
	quitting     bool
	view         string // "menu", "tools", "config", "projects"
	toolStatuses map[string]string
	checker      *ToolChecker
	showingForm  bool
	form         *formModel
}

func (m model) Init() tea.Cmd {
	// Load tool statuses asynchronously
	return func() tea.Msg {
		return toolStatusMsg(m.checker.GetAllToolStatuses())
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle form if it's showing
	if m.showingForm && m.form != nil {
		updatedForm, cmd := m.form.Update(msg)
		if fm, ok := updatedForm.(formModel); ok {
			m.form = &fm
			if m.form.submitted {
				// Form submitted - go back to menu
				m.showingForm = false
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
			} else if m.form.cancelled {
				// Form cancelled - go back to config menu
				m.showingForm = false
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
			}
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case toolStatusMsg:
		m.toolStatuses = map[string]string(msg)
		// Update tool status view if we're viewing it
		if m.view == "tools" {
			m.updateToolStatusItems()
		}
		return m, nil

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
						m.updateToolStatusItems()
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
				case "config":
					// Handle API key configuration
					toolName := ""
					switch m.choice {
					case "Configure Claude API Key":
						toolName = "Claude"
					case "Configure Gemini API Key":
						toolName = "Gemini"
					case "Configure Codex API Key":
						toolName = "Codex"
					case "Configure OpenCode API Key":
						toolName = "OpenCode"
					case "Configure CodeBuddy API Key":
						toolName = "CodeBuddy"
					case "Configure Qoder API Key":
						toolName = "Qoder"
					}
					if toolName != "" {
						// Show form for API key configuration
						formModel := initialFormModel(toolName)
						m.form = &formModel
						m.showingForm = true
						return m, m.form.Init()
					}
				case "launch":
					// Handle tool launch - for now just show a message
					return m, tea.Quit
				case "tools", "projects":
					// For now, just show a message
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) updateToolStatusItems() {
	items := []list.Item{}
	tools := []string{"Claude", "Gemini", "Codex", "OpenCode", "CodeBuddy", "Qoder", "IFlow", "Kilo"}
	
	for _, tool := range tools {
		status := m.toolStatuses[tool]
		if status == "" {
			status = "Checking..."
		}
		items = append(items, item(fmt.Sprintf("%s - %s", tool, status)))
	}
	
	m.list.SetItems(items)
}

func (m model) View() string {
	if m.showingForm && m.form != nil {
		return m.form.View()
	}

	if m.choice != "" && m.view != "menu" && m.view != "tools" && m.view != "config" && m.view != "projects" && m.view != "launch" {
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

	checker := NewToolChecker()
	m := model{
		list:         l,
		view:         "menu",
		checker:      checker,
		toolStatuses: make(map[string]string),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}

	return nil
}
