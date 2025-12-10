package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Structures ---

type Shell struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Meta    []string `json:"meta"`
}

func (s Shell) Title() string { return s.Name }

func (s Shell) Description() string {
	if len(s.Meta) > 0 {
		return "Meta: " + strings.Join(s.Meta, ", ")
	}
	return ""
}
func (s Shell) FilterValue() string { return s.Name }

type EncodingType int

const (
	None EncodingType = iota
	Base64
	URL
	DoubleURL
)

// NOUVELLE STRUCTURE : Pour stocker le nom de la carte et son IP
type NetworkInterface struct {
	Name string
	IP   string
}

type Model struct {
	Inputs []textinput.Model

	ActiveBlock int
	InputIndex  int

	ShellList     list.Model
	Shells        []Shell
	SelectedShell Shell
	Encoding      EncodingType

	// Modification ici : On stocke des objets NetworkInterface au lieu de strings
	Interfaces       []NetworkInterface
	CurrentInterface int

	Width  int
	Height int
}

// --- Helpers ---

// getNetworkInterfaces scanne les cartes et retourne le couple (Nom, IP)
func getNetworkInterfaces() []NetworkInterface {
	var results []NetworkInterface

	ifaces, err := net.Interfaces()
	if err != nil {
		return []NetworkInterface{{Name: "lo", IP: "127.0.0.1"}}
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// On garde IPv4 et on ignore loopback pour l'affichage principal (sauf si rien d'autre)
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				results = append(results, NetworkInterface{
					Name: i.Name,
					IP:   ip.String(),
				})
			}
		}
	}

	// Fallback si rien trouvé (offline)
	if len(results) == 0 {
		results = append(results, NetworkInterface{Name: "local", IP: "127.0.0.1"})
	}
	return results
}

// --- Initialisation ---

func initialModel() Model {
	greenPromptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("43"))
	data, err := os.ReadFile("shells.json")
	if err != nil {
		fmt.Printf("Error reading shells.json: %v\n", err)
		os.Exit(1)
	}

	var shells []Shell
	if err := json.Unmarshal(data, &shells); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	// 2. Détection des Interfaces
	netIfaces := getNetworkInterfaces()
	defaultIface := netIfaces[0]

	// 3. Initialisation des inputs
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "10.10.10.10"
	inputs[0].SetValue(defaultIface.IP) // Valeur par défaut
	inputs[0].Focus()
	inputs[0].Prompt = "Listener IP: "
	inputs[0].CharLimit = 15
	inputs[0].PromptStyle = greenPromptStyle

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "9001"
	inputs[1].Prompt = "Listener Port: "
	inputs[1].CharLimit = 5
	inputs[1].PromptStyle = greenPromptStyle

	items := make([]list.Item, len(shells))
	for i, s := range shells {
		items[i] = s
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(FocusBorderColor).BorderForeground(FocusBorderColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(FocusBorderColor).BorderForeground(FocusBorderColor)

	shellList := list.New(items, delegate, 0, 0)
	shellList.SetShowHelp(false)
	shellList.SetShowTitle(false)
	shellList.SetHeight(8)

	return Model{
		Inputs:           inputs,
		ActiveBlock:      0,
		InputIndex:       0,
		ShellList:        shellList,
		Shells:           shells,
		SelectedShell:    shells[0],
		Encoding:         None,
		Interfaces:       netIfaces, // Stockage des interfaces
		CurrentInterface: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
