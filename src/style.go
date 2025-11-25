// style_constants.go

package main

import "github.com/charmbracelet/lipgloss"

// --- Styles et Constantes Globales ---

// GeneralSpacing est la marge et le padding standard (doit être déclaré en premier)
var GeneralSpacing = 2

// Définitions des couleurs pour Focus/Inactive
var (
	FocusBorderColor    = lipgloss.Color("205") // Violet/Magenta pour l'état actif
	InactiveBorderColor = lipgloss.Color("250") // Gris clair pour l'état inactif
)

// --- Palette de Couleurs (Nerd Fonts friendly) ---
var (
	FocusColor    = lipgloss.Color("205") // Violet/Magenta
	InactiveColor = lipgloss.Color("250") // Gris clair
	GreenColor    = lipgloss.Color("43")  // Vert Matrix
	SubTextColor  = lipgloss.Color("241") // Gris foncé
	WhiteColor    = lipgloss.Color("#FFFDF5")
)
var (
	appStyle = lipgloss.NewStyle().Padding(GeneralSpacing)

	// titleStyle: Style standard pour un titre normal sur fond coloré (en haut du bloc).
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")). // Sera écrasé par la couleur du focus
			Bold(true).
			Padding(0, 1)

	// baseBlockStyle: Utilise la bordure normale et le padding interne.
	baseBlockStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, GeneralSpacing)

	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)

	// Styles pour les options d'encodage
	activeOptionStyle   = lipgloss.NewStyle().Foreground(FocusBorderColor).Bold(true)
	inactiveOptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// Pied de page (Barre de statut)
	HelpStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Italic(true)
)
