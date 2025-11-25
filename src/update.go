package main

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Navigation entre les blocs
		case "tab", "shift+tab":
			if m.ActiveBlock == 0 {
				m.ActiveBlock = 1
				m.Inputs[0].Blur()
				m.Inputs[1].Blur()
			} else {
				m.ActiveBlock = 0
				m.Inputs[m.InputIndex].Focus()
			}
			return m, nil

		case "ctrl+y":
			finalCmd, _ := m.generatePayload()
			clipboard.WriteAll(finalCmd)
			return m, nil
		}

		// Gestion par bloc actif
		if m.ActiveBlock == 0 {
			// Bloc Inputs
			switch msg.String() {
			case "up":
				m.InputIndex = 0
			case "down", "enter":
				m.InputIndex = 1
			}

			if m.InputIndex == 0 {
				m.Inputs[0].Focus()
				m.Inputs[1].Blur()
			} else {
				m.Inputs[0].Blur()
				m.Inputs[1].Focus()
			}
		} else {
			// Bloc Liste & Options
			switch msg.String() {
			case "n", "N":
				m.Encoding = None
			case "b", "B":
				m.Encoding = Base64
			case "u", "U":
				m.Encoding = URL
			case "d", "D":
				m.Encoding = DoubleURL
			default:
				var cmd tea.Cmd
				m.ShellList, cmd = m.ShellList.Update(msg)
				if i, ok := m.ShellList.SelectedItem().(Shell); ok {
					m.SelectedShell = i
				}
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		// Ajustement approximatif, le vrai calcul est dans View
		m.ShellList.SetWidth((msg.Width / 2) - 10)
	}

	// Mise à jour des inputs (toujours nécessaire pour le curseur/texte)
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.Inputs {
		var cmd tea.Cmd
		m.Inputs[i], cmd = m.Inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}