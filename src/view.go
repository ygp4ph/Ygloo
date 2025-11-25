package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Déclarations des variables et constantes globales définies dans model.go.

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
	borderColor := InactiveColor
	if isActive {
		borderColor = FocusColor
	}

	// NOTE: Si vous avez mis à jour lipgloss vers v0.10.0+, vous pouvez utiliser BorderTitle ici.
	// Sinon, cette version fonctionne avec les anciennes versions (titre géré par le contenu ou ignoré pour l'instant).
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

	usableWidth := m.Width // - appDecoration
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

	encLine := fmt.Sprintf("%s  %s  %s  %s",
		renderOption("[N]one", m.Encoding == None),
		renderOption("[B]ase64", m.Encoding == Base64),
		renderOption("[U]rl", m.Encoding == URL),
		renderOption("[D]oubleUrl", m.Encoding == DoubleURL),
	)

	optionsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.ShellList.View(),
		labelStyle.Render("Encoding"),
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

	// --- Bloc 1 : Configuration (Ajustement du padding) ---

	configContent := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("Listenner settings\n"),
		m.Inputs[0].View(),
		labelStyle.Render("\n"),
		m.Inputs[1].View(),
		labelStyle.Render("(Default: 9001)"),
	)

	currentConfigHeight := lipgloss.Height(configContent)
	verticalPaddingToAdd := targetContentHeight - currentConfigHeight

	paddedConfigContent := lipgloss.JoinVertical(
		lipgloss.Left,
		configContent,
		strings.Repeat("\n", verticalPaddingToAdd-1),
	)

	configBlock := renderBlock("Configuration", paddedConfigContent, m.ActiveBlock == 0, configBlockWidth, targetContentHeight)

	// --- Bloc 2 : Payloads & Options ---
	optionsBlock := renderBlock("Payloads & Options", optionsContent, m.ActiveBlock == 1, optionsBlockWidth, targetContentHeight)

	// 3. CALCUL DE LA HAUTEUR DU BLOC DU BAS (Output Large)
	// IMPORTANT : On soustrait footerHeight ici pour laisser la place en bas
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

	// 4. foooter
	// Assurez-vous que HelpStyle est défini dans style.go, sinon remplacez par labelStyle
	helpBar := HelpStyle.Render("\n\n 󰌒  Tab: Switch Block   •     Ctrl+Y: Copy Payload   •   󰅚  q: Quit")

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
