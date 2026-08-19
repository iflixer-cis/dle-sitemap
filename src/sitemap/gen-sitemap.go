package sitemap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SmSitemap struct {
	f            *os.File
	w            *bufio.Writer
	Lastmod      string
	Domain       string
	FileName     string
	FileFullName string
	maxRows      int
	rowsCount    int
	partNumber   int
	fileNames    []string
}

const maxSitemapRows = 50000

type SmSitemapRow struct {
	Lastmod    string
	Loc        string
	ChangeFreq string
	Priority   string
}

func (sf *SmSitemap) Init(domain, outDir, fileName string) (err error) {
	sf.Domain = domain
	sf.FileName = fileName
	sf.maxRows = maxSitemapRows
	sf.rowsCount = 0
	sf.partNumber = 1
	sf.fileNames = nil
	tz := time.FixedZone("EET", 2*60*60)
	sf.Lastmod = time.Date(2025, 10, 27, 3, 30, 3, 0, tz).Format(time.RFC3339)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	return sf.openNewPart(outDir)
}

func (sf *SmSitemap) Add(row SmSitemapRow) error {
	if sf.rowsCount >= sf.maxRows {
		if err := sf.rotatePart(); err != nil {
			return err
		}
	}

	if row.Lastmod == "" {
		row.Lastmod = sf.Lastmod
	}
	res := fmt.Sprintf(`<url>
		<loc>%s</loc>
		<changefreq>%s</changefreq>
		<lastmod>%s</lastmod>
		<priority>%s</priority>
	</url>`, row.Loc, row.ChangeFreq, row.Lastmod, row.Priority)
	fmt.Fprint(sf.w, sf.removeBadSymbols(res))
	sf.rowsCount++
	return nil
}

func (sf *SmSitemap) FileNames() []string {
	out := make([]string, len(sf.fileNames))
	copy(out, sf.fileNames)
	return out
}

func (sf *SmSitemap) rotatePart() error {
	outDir := filepath.Dir(sf.FileFullName)
	if err := sf.closeCurrentFile(); err != nil {
		return err
	}
	sf.partNumber++
	sf.rowsCount = 0
	return sf.openNewPart(outDir)
}

func (sf *SmSitemap) openNewPart(outDir string) error {
	partFileName := sf.partFileName()
	sf.FileFullName = filepath.Join(outDir, partFileName)
	f, err := os.Create(sf.FileFullName)
	if err != nil {
		return err
	}
	sf.f = f

	sf.w = bufio.NewWriter(sf.f)
	fmt.Fprintln(sf.w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprint(sf.w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	sf.fileNames = append(sf.fileNames, partFileName)
	return nil
}

func (sf *SmSitemap) partFileName() string {
	if sf.partNumber <= 1 {
		return sf.FileName
	}

	ext := filepath.Ext(sf.FileName)
	base := strings.TrimSuffix(sf.FileName, ext)
	return fmt.Sprintf("%s_%d%s", base, sf.partNumber, ext)
}

func (sf *SmSitemap) removeBadSymbols(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' || r == ' ' {
			return -1 // выкинуть символ
		}
		return r
	}, s)
}

func (sf *SmSitemap) Close() error {
	return sf.closeCurrentFile()
}

func (sf *SmSitemap) closeCurrentFile() error {
	if sf.f == nil || sf.w == nil {
		return nil
	}

	fmt.Fprint(sf.w, `</urlset>`)
	if err := sf.w.Flush(); err != nil {
		_ = sf.f.Close()
		sf.f = nil
		sf.w = nil
		return err
	}

	err := sf.f.Close()
	sf.f = nil
	sf.w = nil
	return err
}
