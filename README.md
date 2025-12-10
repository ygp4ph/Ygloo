<div align="center"><h1>Ygloo</h1></div>

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey?style=flat-square)

Un générateur de reverse shells interactif en TUI ma foi charmant
</div>


## Fonctionnalités

- **Interface TUI moderne** - Navigation intuitive avec Bubble Tea
- **30+ reverse shells** - Bash, Python, PHP, PowerShell, Netcat, et plus
- **Détection automatique d'IP** - Switch entre interfaces réseau (Ctrl+N)
- **Encodage multiple** - None, Base64, URL, Double URL
- **Copie instantanée** - Ctrl+Y pour copier dans le presse-papiers
- **Bind shells inclus** - Python, Netcat, Perl
- **Configuration JSON** - Facile à étendre et personnaliser

## Installation

```bash
# Clone le repository
git clone https://github.com/yourusername/revshell-tui.git
cd revshell-tui

# Build
cd src
go build -o revshell-tui

# Lance l'application
./revshell-tui
```
## Démonstration

![alt text](image.png)


### Dépendances

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [Bubbles](https://github.com/charmbracelet/bubbles) - Composants TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [clipboard](https://github.com/atotto/clipboard) - Gestion presse-papiers