package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	inputPDF  = "test1.pdf"
	outputDir = "output"
	dpi       = 150
	format    = "png"
)

func pdfToImages() error {
	if _, err := exec.LookPath("pdftoppm"); err != nil {
		return fmt.Errorf("pdftoppm not found (install poppler: brew install poppler): %w", err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPDF), filepath.Ext(inputPDF))
	outPrefix := filepath.Join(outputDir, baseName)

	args := []string{"-r", fmt.Sprintf("%d", dpi)}
	ext := format
	switch format {
	case "png":
		args = append(args, "-png")
	case "jpeg", "jpg":
		args = append(args, "-jpeg")
		ext = "jpg"
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	args = append(args, inputPDF, outPrefix)

	cmd := exec.Command("pdftoppm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pdftoppm: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(outputDir, baseName+"-*."+ext))
	if err != nil {
		return fmt.Errorf("list output files: %w", err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	fmt.Printf("converted %d page(s) to %s\n", len(files), outputDir)
	return nil
}

func main() {
	if err := pdfToImages(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
