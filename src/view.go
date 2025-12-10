package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- Génération du Payload ---
func (m Model) generatePayload() (string, string) {
	ip := m.Inputs[0].Value()
	if ip == "" {
		ip = "10.10.10.10"
	}

	port := m.Inputs[1].Value()
	if port == "" {
		port = "9001"
	}

	raw := m.SelectedShell.Command
	raw = strings.ReplaceAll(raw, "{shell}", "/bin/bash")
	raw = strings.ReplaceAll(raw, "{ip}", ip)
	raw = strings.ReplaceAll(raw, "{port}", port)
	payload := raw

	switch m.Encoding {
	case Base64:
		payload = base64.StdEncoding.EncodeToString([]byte(payload))
	case URL:
		payload = url.QueryEscape(payload)
	case DoubleURL:
		payload = url.QueryEscape(url.QueryEscape(payload))
	}

	currentListener := m.Listeners[m.ListenerIndex]
	listener := strings.ReplaceAll(currentListener.Template, "{port}", port)

	return payload, listener
}

// --- Rendu Visuel ---
func renderOption(label string, active bool) string {
	if active {
		return activeOptionStyle.Render(label)
	}
	return inactiveOptionStyle.Render(label)
}

func renderBlock(content string, isActive bool, contentWidth, contentHeight int) string {
	style := baseBlockStyle.Copy()
	borderColor := InactiveBorderColor
	if isActive {
		borderColor = FocusBorderColor
	}
	block := style.
		BorderForeground(borderColor).
		Width(contentWidth).
		Height(contentHeight).
		SetString(content)
	return block.Render()
}

// View génère l'interface utilisateur complète
func (m Model) View() string {
	appDecoration := 2 * GeneralSpacing
	blockDecorationWidth := 2 * (GeneralSpacing + 1)
	footerHeight := 1
	usableWidth := m.Width
	usableHeight := m.Height - appDecoration
	blockVerticalDecoration := 2 * (GeneralSpacing + 1)
	totalDecorationLoss := (2 * blockDecorationWidth) + GeneralSpacing
	remainingWidth := usableWidth - totalDecorationLoss
	if remainingWidth < 0 {
		remainingWidth = 0
	}
	halfContentWidth := remainingWidth / 2
	configBlockWidth := halfContentWidth - 8
	optionsBlockWidth := remainingWidth - configBlockWidth + 3
	if configBlockWidth < 1 {
		configBlockWidth = 1
	}
	if optionsBlockWidth < 1 {
		optionsBlockWidth = 1
	}

	m.ShellList.SetWidth(optionsBlockWidth - (2 * GeneralSpacing))

	encLine := fmt.Sprintf("%s  %s  %s  %s\n",
		renderOption("[N]one", m.Encoding == None),
		renderOption("[B]ase64", m.Encoding == Base64),
		renderOption("[U]rl", m.Encoding == URL),
		renderOption("[D]oubleUrl", m.Encoding == DoubleURL),
	)

	optionsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.ShellList.View(),
		labelStyle.Render("Encoding\n"),
		encLine,
	)

	listContentHeight := lipgloss.Height(m.ShellList.View())
	optionsLabelsHeight := lipgloss.Height(labelStyle.Render("Encoding"))
	encLineHeight := lipgloss.Height(encLine)
	targetContentHeight := listContentHeight + optionsLabelsHeight + encLineHeight + 1
	minStaticHeight := 5
	if targetContentHeight < minStaticHeight {
		targetContentHeight = minStaticHeight
	}
	topRowHeight := targetContentHeight + blockVerticalDecoration

	// Gestion affichage Interface
	currentIP := m.Inputs[0].Value()
	interfaceLabel := ""
	foundIface := false
	for _, iface := range m.Interfaces {
		if iface.IP == currentIP {
			interfaceLabel = lipgloss.NewStyle().Foreground(GreenColor).Render(fmt.Sprintf("󰈀 %s", iface.Name))
			foundIface = true
			break
		}
	}
	if !foundIface && len(m.Interfaces) > 0 {
		interfaceLabel = labelStyle.Render("(Ctrl+n: Switch Interface)")
	} else if foundIface && len(m.Interfaces) > 1 {
		interfaceLabel += labelStyle.Render("")
	}
	// Gestion affichage Sélecteur de Listener
	isListenerActive := m.ActiveBlock == 0 && m.InputIndex == 2
	lName := m.Listeners[m.ListenerIndex].Name
	arrowLeft := " "
	arrowRight := " "
	listenerStyle := lipgloss.NewStyle().Foreground(WhiteColor)

	if isListenerActive {
		listenerStyle = lipgloss.NewStyle().Foreground(FocusColor).Bold(true)
		arrowLeft = "< "
		arrowRight = " >"
	}

	listenerSelector := fmt.Sprintf("%s%s%s",
		lipgloss.NewStyle().Foreground(SubTextColor).Render(arrowLeft),
		listenerStyle.Render(lName),
		lipgloss.NewStyle().Foreground(SubTextColor).Render(arrowRight),
	)
	// Contenu du bloc configuration
	configContent := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("\nListener settings\n"),
		m.Inputs[0].View(),
		interfaceLabel,

		m.Inputs[1].View(),

		labelStyle.Render("\nTool Selection"),
		listenerSelector,
	)
	// Ajustement hauteur bloc config
	currentConfigHeight := lipgloss.Height(configContent)
	verticalPaddingToAdd := targetContentHeight - currentConfigHeight
	safePadding := verticalPaddingToAdd - 1
	if safePadding < 0 {
		safePadding = 0
	}

	paddedConfigContent := lipgloss.JoinVertical(
		lipgloss.Left,
		configContent,
		strings.Repeat("\n", safePadding),
	)
	// Rendu final des blocs
	configBlock := renderBlock(paddedConfigContent, m.ActiveBlock == 0, configBlockWidth, targetContentHeight)
	optionsBlock := renderBlock(optionsContent, m.ActiveBlock == 1, optionsBlockWidth, targetContentHeight)
	outputContentHeight := usableHeight - topRowHeight - GeneralSpacing - footerHeight
	if outputContentHeight < 1 {
		outputContentHeight = 1
	}
	payload, listener := m.generatePayload()
	outputContentWidth := usableWidth - blockDecorationWidth
	outputContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("\nVICTIM COMMAND :\n"),
		lipgloss.NewStyle().Foreground(FocusColor).Width(outputContentWidth).Render(payload),
		"\n",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("LISTENER COMMAND :\n"),
		lipgloss.NewStyle().Foreground(GreenColor).Render(listener),
		labelStyle.Render("(Ctrl+Y to copy)"),
	)
	outputFinalBlock := renderBlock(outputContent, false, outputContentWidth, outputContentHeight)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, configBlock, lipgloss.NewStyle().Width(3).Render(""), optionsBlock)
	verticalSpacer := lipgloss.NewStyle().Height(1).Width(usableWidth).Render("")
	helpBar := HelpStyle.Width(m.Width).Align(lipgloss.Center).Render("\n\n 󰌒  Tab : Switch Block   •     Ctrl+Y : Copy Payload   •   󰈆  q : Quit   •   󰀂  Ctrl+n : Switch Interface")

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, topRow, verticalSpacer, outputFinalBlock, helpBar))
}
