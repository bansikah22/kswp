package branding

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	// Version represents the application version
	Version = "0.1.0"
	// Tagline is the application tagline
	Tagline = "Kubernetes Cluster Hygiene Tool"
)

// DisplayBanner displays the KSWP ASCII banner with version and tagline
func DisplayBanner() {
	banner := `
  _                       
 | |                      
 | | _______      ___ __  
 | |/ / __\ \ /\ / / '_ \ 
 |   <\__ \\ V  V /| |_) |
 |_|\_\___/ \_/\_/ | .__/ 
                   | |    
                   |_|    
`

	// Create styling for the banner
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")). // Cyan color
		Bold(true)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // Gray color
		Italic(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("48")). // Green color
		Bold(true)

	// Print banner
	fmt.Println(style.Render(banner))

	// Print version and tagline
	fmt.Println("  " + versionStyle.Render(fmt.Sprintf("v%s", Version)))
	fmt.Println("  " + taglineStyle.Render(Tagline))
	fmt.Println()
}

// GetBannerString returns the ASCII banner as a string (useful for testing)
func GetBannerString() string {
	return `
  _                       
 | |                      
 | | _______      ___ __  
 | |/ / __\ \ /\ / / '_ \ 
 |   <\__ \\ V  V /| |_) |
 |_|\_\___/ \_/\_/ | .__/ 
                   | |    
                   |_|  
`
}

// GetVersionString returns the version information
func GetVersionString() string {
	return fmt.Sprintf("v%s", Version)
}

// GetBrandInfo returns text info about the application branding
func GetBrandInfo() string {
	return strings.Join([]string{
		GetBannerString(),
		fmt.Sprintf("Version: %s", GetVersionString()),
		fmt.Sprintf("Description: %s", Tagline),
	}, "\n")
}
