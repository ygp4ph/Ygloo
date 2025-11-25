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

// Rendu unifié du bloc avec titre NORMAL et bordure standard (CORRIGÉ).
func renderBlock(title string, content string, isActive bool, contentWidth, contentHeight int) string {

	style := baseBlockStyle.Copy()

	// Choix de la couleur en fonction de l'état actif
	borderColor := InactiveBorderColor
	if isActive {
		borderColor = FocusBorderColor
	}

	// 2. Le style Lip Gloss prend la bordure normale et définit la taille et la couleur.
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

	usableWidth := m.Width // - appDecoration
	usableHeight := m.Height - appDecoration

	blockVerticalDecoration := 2 * (GeneralSpacing + 1) // Bordure verticale + Padding
	// 1. CALCUL DES LARGEURS (Identique à la version précédente pour le responsive)
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

	// D'abord, configurer la liste de droite pour obtenir sa hauteur naturelle
	// La liste utilise optionsBlockWidth pour son rendu.
	m.ShellList.SetWidth(optionsBlockWidth - (2 * GeneralSpacing))
	// Note: m.ShellList.SetHeight est fait dans initialModel, nous laissons la liste calculer sa hauteur.

	encLine := fmt.Sprintf("%s  %s  %s  %s",
		renderOption("[N]one", m.Encoding == None),
		renderOption("[B]ase64", m.Encoding == Base64),
		renderOption("[U]rl", m.Encoding == URL),
		renderOption("[D]oubleUrl", m.Encoding == DoubleURL),
	)

	optionsContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.ShellList.View(),            // Hauteur dépendante de m.ShellList.Height()
		labelStyle.Render("Encoding"), // 1 ligne + marge
		encLine,                       // 1 ligne
	)

	// Calculer la hauteur interne ciblée du contenu (Bloc de Droite)
	// Hauteur de la liste (fixée dans model.go, mais peut changer) + Hauteur des labels/options
	listContentHeight := lipgloss.Height(m.ShellList.View())
	optionsLabelsHeight := lipgloss.Height(labelStyle.Render("Encoding"))
	encLineHeight := lipgloss.Height(encLine)

	// Hauteur totale du contenu interne du bloc Options (Liste + Label + Ligne options)
	targetContentHeight := listContentHeight + optionsLabelsHeight + encLineHeight + 1 // +1 pour la marge/séparation

	// La hauteur réelle du contenu du bloc doit être au moins la hauteur statique du bloc de gauche (environ 5)
	minStaticHeight := 5 // (IP, Port, 2 Labels, Default)
	if targetContentHeight < minStaticHeight {
		targetContentHeight = minStaticHeight
	}

	// La hauteur totale du bloc (avec décorations)
	topRowHeight := targetContentHeight + blockVerticalDecoration
	// --- Bloc 1 : Configuration (Ajustement du padding) ---

	// Contenu statique du bloc de gauche (IP, Port, Labels)
	configContent := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("Listenner settings\n"),
		m.Inputs[0].View(),
		labelStyle.Render("\n"),
		m.Inputs[1].View(),
		labelStyle.Render("(Default: 9001)"),
	)

	// Calculer l'espace nécessaire pour combler la différence de hauteur
	// Hauteur cible - Hauteur actuelle du contenu de gauche
	currentConfigHeight := lipgloss.Height(configContent)

	// Padding vertical à ajouter pour atteindre la targetContentHeight
	verticalPaddingToAdd := targetContentHeight - currentConfigHeight

	// On ajoute le padding au bas du contenu.
	paddedConfigContent := lipgloss.JoinVertical(
		lipgloss.Left,
		configContent,
		strings.Repeat("\n", verticalPaddingToAdd-1), // Ajout des lignes vides pour augmenter la hauteur
	)

	configBlock := renderBlock("Configuration", paddedConfigContent, m.ActiveBlock == 0, configBlockWidth, targetContentHeight)

	// --- Bloc 2 : Payloads & Options ---
	// La hauteur du contenu est implicitement targetContentHeight
	optionsBlock := renderBlock("Payloads & Options", optionsContent, m.ActiveBlock == 1, optionsBlockWidth, targetContentHeight)

	// ... (Le reste du code, y compris le Bloc 3 (Output), utilise maintenant targetContentHeight et topRowHeight) ...

	// 3. CALCUL DE LA HAUTEUR DU BLOC DU BAS (Output Large)
	outputContentHeight := usableHeight - topRowHeight - GeneralSpacing //- blockVerticalDecoration

	if outputContentHeight < 1 {
		outputContentHeight = 1
	}
	// --- Bloc 3 : Output (Bas - Grande Zone T-Shape) ---

	payload, listener := m.generatePayload()

	// Largeur du contenu du bloc du bas
	outputContentWidth := usableWidth - blockDecorationWidth

	outputContent := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("\nVICTIM COMMAND :\n"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Width(outputContentWidth).Render(payload),
		"\n",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A0C4FF")).Render("LISTENER COMMAND :\n"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Render(listener),
		labelStyle.Render("(Ctrl+Y to copy)"),
	)

	// Le bloc Output n'a jamais le focus actif (false)
	outputFinalBlock := renderBlock("Output (Large Area)", outputContent, false, outputContentWidth, outputContentHeight)

	// --- Assemblage Final ---

	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		configBlock,
		lipgloss.NewStyle().Width(3).Render(""), // Espace entre les blocs
		optionsBlock,
	)

	verticalSpacer := lipgloss.NewStyle().Height(1).Width(usableWidth).Render("")

	// 3. Jointure Verticale finale (T-Shape)
	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			topRow,
			verticalSpacer,
			outputFinalBlock,
		),
	)
}
