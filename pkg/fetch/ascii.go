package fetch

import (
	"bufio"
	"math"
	"os"
	"strings"

	asci "github.com/minhnh/fetch/internal/ascii"
)

type AsciiArt struct {
	Lines          []string
	MaxCleanLen    int
	MaxOriginalLen int
}

func processLines(scanner *bufio.Scanner) ([]string, int, int) {
	clean := math.MinInt
	original := math.MinInt
	lines := []string{}

	for scanner.Scan() {
		asc := scanner.Text()
		lines = append(lines, asc)

		original = Max(original, len(asc))
		for placeholder := range asci.PlaceHolder {
			asc = strings.ReplaceAll(asc, placeholder, "")
		}
		clean = Max(len(asc), clean)
	}

	return lines, clean, original
}

func NewAsciiArt(filePath string) (*AsciiArt, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines, clean, original := processLines(scanner)

	return &AsciiArt{
		Lines:          lines,
		MaxCleanLen:    clean,
		MaxOriginalLen: original,
	}, nil
}

func DefaultArt(sys string) *AsciiArt {
	scanner := bufio.NewScanner(strings.NewReader(asci.Art[sys]))
	lines, clean, original := processLines(scanner)

	return &AsciiArt{
		Lines:          lines,
		MaxCleanLen:    clean,
		MaxOriginalLen: original,
	}
}
