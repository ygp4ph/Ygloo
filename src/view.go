package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- Logique Métier ---

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

	// Encoding management
	switch m.Encoding {
	case Base64:
		payload = base64.StdEncoding.EncodeToString([]byte(payload))
	case URL:
		payload = url.QueryEscape(payload)
	case DoubleURL:
		payload = url.QueryEscape(url.QueryEscape(payload))
	}

	listener := fmt.Sprintf("nc -lvnp %s", port)
	return payload, listener
}

// --- Helpers de style ---

func renderOption(label string, active bool) string {
	if active {
		return activeOptionStyle.Render(label)
	}
	return inactiveOptionStyle.Render(label)
}

// Rendu unifié du bloc
func renderBlock(title string, content string, isActive bool, contentWidth, contentHeight int) string {

	style := baseBlockStyle.Copy()

	// Choix de la couleur en fonction de l'état actif
	borderColor := InactiveBorderColor
	if isActive {
		borderColor = FocusBorderColor
	}

	block := style.
		BorderForeground(borderColor). // Applique la couleur à la bordure
		Width(contentWidth).
		Height(contentHeight).
		SetString(content)

	// Retourne le bloc
	return block.Render()
}

// --- View ---
func (m Model) View() string {
	// --- CALCUL DES DIMENSIONS STABLES ---

	appDecoration := 2 * GeneralSpacing
	blockDecorationWidth := 2 * (GeneralSpacing + 1) // Décoration d'un seul bloc (Bordures + Padding)

	// On définit la hauteur du footer
	footerHeight := 1

	usableWidth := m.Width
	usableHeight := m.Height - appDecoration

	blockVerticalDecoration := 2 * (GeneralSpacing + 1) // Bordure verticale + Padding

	// 1. CALCUL DES LARGEURS
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

	// 2. --- CALCUL DYNAMIQUE DE LA HAUTEUR CIBLE (TARGET HEIGHT) ---

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

	// Calculer la hauteur interne ciblée du contenu
	listContentHeight := lipgloss.Height(m.ShellList.View())
	optionsLabelsHeight := lipgloss.Height(labelStyle.Render("Encoding"))
	encLineHeight := lipgloss.Height(encLine)

	targetContentHeight := listContentHeight + optionsLabelsHeight + encLineHeight + 1

	minStaticHeight := 5
	if targetContentHeight < minStaticHeight {
		targetContentHeight = minStaticHeight
	}

	topRowHeight := targetContentHeight + blockVerticalDecoration

	// --- Bloc 1 : Configuration (MODIFIÉ POUR LES INTERFACES) ---

	// 1. On récupère l'IP actuelle du champ texte
	currentIP := m.Inputs[0].Value()
	interfaceLabel := ""
	foundIface := false

	// 2. On compare avec nos interfaces détectées
	for _, iface := range m.Interfaces {
		if iface.IP == currentIP {
			// Si on trouve une correspondance, on affiche le nom en vert (ex: "󰈀 eth0")
			interfaceLabel = lipgloss.NewStyle().Foreground(GreenColor).Render(fmt.Sprintf("󰾲 %s", iface.Name))
			foundIface = true
			break
		}
	}

	// 3. Gestion du texte d'aide si on ne trouve pas ou pour rappeler le raccourci
	if !foundIface && len(m.Interfaces) > 0 {
		interfaceLabel = labelStyle.Render("(Ctrl+n: Switch Interface)")
	} else if !foundIface {
		interfaceLabel = labelStyle.Render("(No network interfaces detected)")
	}

	configContent := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("Listenner settings\n"),
		m.Inputs[0].View(),
		interfaceLabel,
		labelStyle.Render("\n"),
		m.Inputs[1].View(),
		labelStyle.Render("(Default: 9001)"),
	)

	currentConfigHeight := lipgloss.Height(configContent)
	verticalPaddingToAdd := targetContentHeight - currentConfigHeight

	// CORRECTION ICI : On retire 1 au padding pour éviter le saut de ligne fantôme
	safePadding := verticalPaddingToAdd - 1
	if safePadding < 0 {
		safePadding = 0
	}

	paddedConfigContent := lipgloss.JoinVertical(
		lipgloss.Left,
		configContent,
		strings.Repeat("\n", safePadding),
	)

	// On passe quand même targetContentHeight au renderBlock,
	// Lipgloss comblera proprement l'espace restant sans dépasser.
	configBlock := renderBlock("Configuration", paddedConfigContent, m.ActiveBlock == 0, configBlockWidth, targetContentHeight)
	// --- Bloc 2 : Payloads & Options ---
	optionsBlock := renderBlock("Payloads & Options", optionsContent, m.ActiveBlock == 1, optionsBlockWidth, targetContentHeight)

	// 3. CALCUL DE LA HAUTEUR DU BLOC DU BAS (Output Large)
	outputContentHeight := usableHeight - topRowHeight - GeneralSpacing - footerHeight

	if outputContentHeight < 1 {
		outputContentHeight = 1
	}

	// --- Bloc 3 : Output ---

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

	outputFinalBlock := renderBlock("Output (Large Area)", outputContent, false, outputContentWidth, outputContentHeight)

	// --- Assemblage Final ---

	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		configBlock,
		lipgloss.NewStyle().Width(3).Render(""),
		optionsBlock,
	)

	verticalSpacer := lipgloss.NewStyle().Height(1).Width(usableWidth).Render("")

	helpBar := HelpStyle.Render("\n\n 󰌒  Tab: Switch Block   •     Ctrl+y: Copy Payload   •   󰈆  q: Quit   •   󰾲  Ctrl+n: Switch Interface")

	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			topRow,
			verticalSpacer,
			outputFinalBlock,
			helpBar,
		),
	)
}
