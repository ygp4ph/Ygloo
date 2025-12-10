package main

import "github.com/charmbracelet/lipgloss"

// --- Styles et Constantes Globales ---

var GeneralSpacing = 2

// --- Palette de Couleurs Dynamique (ANSI) ---
// En utilisant les strings "0" à "15", Lipgloss utilise la palette du terminal de l'utilisateur.

var (
	// 5 = Magenta (souvent utilisé pour la sélection/focus)
	FocusBorderColor = lipgloss.Color("5")

	// 8 = Bright Black (Gris foncé, idéal pour les bordures inactives)
	InactiveBorderColor = lipgloss.Color("8")
)

var (
	FocusColor    = lipgloss.Color("5")  // Magenta standard
	InactiveColor = lipgloss.Color("8")  // Gris standard
	GreenColor    = lipgloss.Color("2")  // Vert standard
	SubTextColor  = lipgloss.Color("8")  // Gris pour les textes secondaires
	WhiteColor    = lipgloss.Color("15") // Blanc brillant
)

var (
	appStyle = lipgloss.NewStyle().Padding(GeneralSpacing)

	// titleStyle: Utilise le vert du terminal.
	// On met le texte en noir (0) pour assurer le contraste sur le fond vert (2).
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(GreenColor).
			Bold(true).
			Padding(0, 1)

	// baseBlockStyle: Bordures épaisses + Couleurs ANSI
	baseBlockStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(0, GeneralSpacing)

	labelStyle = lipgloss.NewStyle().Foreground(SubTextColor).MarginTop(1)

	// Styles pour les options d'encodage
	activeOptionStyle   = lipgloss.NewStyle().Foreground(FocusColor).Bold(true)
	inactiveOptionStyle = lipgloss.NewStyle().Foreground(InactiveColor)

	// Pied de page (Barre de statut)
	HelpStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Italic(true)
)
