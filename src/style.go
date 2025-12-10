package main

import "github.com/charmbracelet/lipgloss"

// --- Styles ---
var GeneralSpacing = 2

var (
	FocusBorderColor    = lipgloss.Color("5")
	InactiveBorderColor = lipgloss.Color("8")
)

// --- Couleurs issues du terminal du user ---
var (
	FocusColor    = lipgloss.Color("5")
	InactiveColor = lipgloss.Color("8")
	GreenColor    = lipgloss.Color("2")
	SubTextColor  = lipgloss.Color("8")
	WhiteColor    = lipgloss.Color("15")
)

var (
	appStyle = lipgloss.NewStyle().Padding(GeneralSpacing)
	// Styles de base pour les blocs avec bordures
	baseBlockStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Padding(0, GeneralSpacing)

	labelStyle = lipgloss.NewStyle().Foreground(SubTextColor).MarginTop(1)

	// Styles pour les options actives/inactives dans le sélecteur
	activeOptionStyle   = lipgloss.NewStyle().Foreground(FocusColor).Bold(true)
	inactiveOptionStyle = lipgloss.NewStyle().Foreground(InactiveColor)

	// Style pour l'aide contextuelle
	HelpStyle = lipgloss.NewStyle().
			Foreground(SubTextColor).
			Italic(true)
)
