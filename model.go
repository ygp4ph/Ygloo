package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Styles et Variables Globales ---

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

type Model struct {
	Inputs []textinput.Model

	ActiveBlock int
	InputIndex  int

	ShellList     list.Model
	Shells        []Shell
	SelectedShell Shell
	Encoding      EncodingType

	Width  int
	Height int
}

// --- Initialisation ---

func initialModel() Model {
	// 1. Charger le JSON depuis le fichier shells.json
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

	// 2. Initialisation des inputs IP/Port
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "10.10.10.10"
	inputs[0].Focus()
	inputs[0].Prompt = "Listener IP: "
	inputs[0].CharLimit = 15
	inputs[0].PromptStyle = greenPromptStyle

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "9001"
	inputs[1].Prompt = "Listener Port: "
	inputs[1].CharLimit = 5
	inputs[1].PromptStyle = greenPromptStyle

	// 3. Initialisation de la liste des shells
	items := make([]list.Item, len(shells))
	for i, s := range shells {
		items[i] = s
	}

	delegate := list.NewDefaultDelegate()
	// Applique la couleur focus aux éléments sélectionnés de la liste
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(FocusBorderColor).BorderForeground(FocusBorderColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(FocusBorderColor).BorderForeground(FocusBorderColor)

	shellList := list.New(items, delegate, 0, 0)
	shellList.SetShowHelp(false)
	shellList.SetShowTitle(false)
	shellList.SetHeight(8)

	return Model{
		Inputs:        inputs,
		ActiveBlock:   0,
		InputIndex:    0,
		ShellList:     shellList,
		Shells:        shells,
		SelectedShell: shells[0],
		Encoding:      None,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
