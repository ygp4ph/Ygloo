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

		case "tab", "shift+tab":
			if m.ActiveBlock == 0 {
				m.ActiveBlock = 1
				m.Inputs[0].Blur()
				m.Inputs[1].Blur()
			} else {
				m.ActiveBlock = 0
				// Retourne au dernier input actif ou au premier
				if m.InputIndex < 2 {
					m.Inputs[m.InputIndex].Focus()
				}
			}
			return m, nil

		case "ctrl+y":
			finalCmd, _ := m.generatePayload()
			clipboard.WriteAll(finalCmd)
			return m, nil
		}

		if m.ActiveBlock == 0 {
			// --- LOGIQUE BLOC CONFIG ---
			switch msg.String() {

			// Navigation verticale (IP -> Port -> Listener)
			case "up", "shift+tab":
				m.InputIndex--
				if m.InputIndex < 0 {
					m.InputIndex = 0
				}
			case "down", "enter", "tab":
				m.InputIndex++
				if m.InputIndex > 2 { // On a maintenant 3 champs (0, 1, 2)
					m.InputIndex = 2
				}

			// Actions spécifiques selon le champ
			case "ctrl+n":
				// Cycle IP seulement si on est sur l'Input 0
				if m.InputIndex == 0 && len(m.Interfaces) > 0 {
					m.CurrentInterface = (m.CurrentInterface + 1) % len(m.Interfaces)
					selectedIface := m.Interfaces[m.CurrentInterface]
					m.Inputs[0].SetValue(selectedIface.IP)
					m.Inputs[0].SetCursor(len(selectedIface.IP))
				}

			// Gestion du Sélecteur de Listener (Index 2)
			case "left", "h":
				if m.InputIndex == 2 {
					m.ListenerIndex--
					if m.ListenerIndex < 0 {
						m.ListenerIndex = len(m.Listeners) - 1
					}
				} else {
					// Propagation standard pour les inputs texte (déplacement curseur)
					m.updateInputs(msg)
				}
			case "right", "l":
				if m.InputIndex == 2 {
					m.ListenerIndex = (m.ListenerIndex + 1) % len(m.Listeners)
				} else {
					m.updateInputs(msg)
				}
			}

			// Focus Management
			// On désactive tout d'abord
			m.Inputs[0].Blur()
			m.Inputs[1].Blur()

			// On active celui qui correspond à l'index
			if m.InputIndex == 0 {
				m.Inputs[0].Focus()
			} else if m.InputIndex == 1 {
				m.Inputs[1].Focus()
			}
			// Si InputIndex == 2, aucun textinput n'a le focus, c'est le sélecteur custom qui est "actif" visuellement

		} else {
			// Bloc Liste (inchangé)
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
		m.ShellList.SetWidth((msg.Width / 2) - 10)
	}

	if m.ActiveBlock == 0 && m.InputIndex == 2 {
	} else {
		cmd := m.updateInputs(msg)
		cmds = append(cmds, cmd)
	}

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
