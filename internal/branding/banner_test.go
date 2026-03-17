package branding

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestGetBannerString(t *testing.T) {
	banner := GetBannerString()

	if !strings.Contains(banner, "_") {
		t.Error("Banner should contain ASCII art with underscores")
	}

	if !strings.Contains(banner, "|") {
		t.Error("Banner should contain ASCII art with pipes")
	}

	if len(banner) == 0 {
		t.Error("Banner should not be empty")
	}
}

func TestGetVersionString(t *testing.T) {
	version := GetVersionString()

	if !strings.HasPrefix(version, "v") {
		t.Error("Version should start with 'v'")
	}

	if !strings.Contains(version, Version) {
		t.Errorf("Version string should contain version constant: %s", Version)
	}
}

func TestGetBrandInfo(t *testing.T) {
	brandInfo := GetBrandInfo()

	if !strings.Contains(brandInfo, "Version:") {
		t.Error("Brand info should contain 'Version:'")
	}

	if !strings.Contains(brandInfo, "Description:") {
		t.Error("Brand info should contain 'Description:'")
	}

	if !strings.Contains(brandInfo, Tagline) {
		t.Errorf("Brand info should contain tagline: %s", Tagline)
	}
}

func TestDisplayBanner(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayBanner()

	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify output contains banner components
	if !strings.Contains(output, "_") {
		t.Error("Banner output should contain ASCII art")
	}

	if !strings.Contains(output, Version) {
		t.Error("Banner output should contain version")
	}

	if !strings.Contains(output, Tagline) {
		t.Error("Banner output should contain tagline")
	}
}

func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	// Version should follow semantic versioning format
	parts := strings.Split(Version, ".")
	if len(parts) < 3 {
		t.Errorf("Version should be in semantic format (major.minor.patch): %s", Version)
	}
}

func TestTaglineConstant(t *testing.T) {
	if Tagline == "" {
		t.Error("Tagline constant should not be empty")
	}

	if len(Tagline) < 5 {
		t.Error("Tagline should be a reasonable length")
	}
}
