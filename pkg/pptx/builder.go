package pptx

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/napcs/opml2pptx/pkg/opml"
)

//go:embed all:templates
var templatesFS embed.FS

type fileEntry struct {
	Name string
	Data []byte
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"mul": func(a, b int) int { return a * b },
	"eq":  func(a, b int) bool { return a == b },
	"xmlEscape": html.EscapeString,
}

// BuildPPTX generates a .pptx file from the Presentation model.
func BuildPPTX(p *opml.Presentation, outputPath string) error {
	files, err := generatePPTXFiles(p)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close()

	for _, f := range files {
		path := filepath.ToSlash(f.Name)
		
		// Create header with proper timestamp
		header := &zip.FileHeader{
			Name:               path,
			Method:             zip.Deflate,
			Modified:           time.Now(),
			UncompressedSize64: uint64(len(f.Data)),
		}
		header.SetMode(0644)
		
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("creating zip entry %s: %w", path, err)
		}

		if _, err := writer.Write(f.Data); err != nil {
			return fmt.Errorf("writing zip entry %s: %w", path, err)
		}
	}

	return nil
}

func generatePPTXFiles(p *opml.Presentation) ([]fileEntry, error) {
	var files []fileEntry
	
	sharedData := map[string]any{
		"Presentation": p,
		"Slides":       p.Slides,
		"Title":        p.Title,
	}

	// Core files
	coreFiles := map[string]string{
		"[Content_Types].xml":                          "templates/[Content_Types].xml.tmpl",
		"_rels/.rels":                                  "templates/_rels/.rels.tmpl",
		"ppt/presentation.xml":                         "templates/ppt/presentation.xml.tmpl",
		"ppt/_rels/presentation.xml.rels":              "templates/ppt/_rels/presentation.xml.rels.tmpl",
		"ppt/presProps.xml":                            "templates/ppt/presProps.xml.tmpl",
		"ppt/viewProps.xml":                            "templates/ppt/viewProps.xml.tmpl",
		"ppt/slideLayouts/slideLayout1.xml":            "templates/ppt/slideLayouts/slideLayout1.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout1.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout2.xml":            "templates/ppt/slideLayouts/slideLayout2.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout2.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout2.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout3.xml":            "templates/ppt/slideLayouts/slideLayout3.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout3.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout3.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout4.xml":            "templates/ppt/slideLayouts/slideLayout4.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout4.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout4.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout5.xml":            "templates/ppt/slideLayouts/slideLayout5.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout5.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout5.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout6.xml":            "templates/ppt/slideLayouts/slideLayout6.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout6.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout6.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout7.xml":            "templates/ppt/slideLayouts/slideLayout7.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout7.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout7.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout8.xml":            "templates/ppt/slideLayouts/slideLayout8.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout8.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout8.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout9.xml":            "templates/ppt/slideLayouts/slideLayout9.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout9.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout9.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout10.xml":           "templates/ppt/slideLayouts/slideLayout10.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout10.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout10.xml.rels.tmpl",
		"ppt/slideLayouts/slideLayout11.xml":           "templates/ppt/slideLayouts/slideLayout11.xml.tmpl",
		"ppt/slideLayouts/_rels/slideLayout11.xml.rels": "templates/ppt/slideLayouts/_rels/slideLayout11.xml.rels.tmpl",
		"ppt/slideMasters/slideMaster1.xml":            "templates/ppt/slideMasters/slideMaster1.xml.tmpl",
		"ppt/slideMasters/_rels/slideMaster1.xml.rels": "templates/ppt/slideMasters/_rels/slideMaster1.xml.rels.tmpl",
		"ppt/theme/theme1.xml":                         "templates/ppt/theme/theme1.xml.tmpl",
		"ppt/theme/theme2.xml":                         "templates/ppt/theme/theme2.xml.tmpl",
		"ppt/notesMasters/notesMaster1.xml":            "templates/ppt/notesMasters/notesMaster1.xml.tmpl",
		"ppt/notesMasters/_rels/notesMaster1.xml.rels": "templates/ppt/notesMasters/_rels/notesMaster1.xml.rels.tmpl",
	}

	for path, tmpl := range coreFiles {
		data, err := renderTemplate(tmpl, sharedData)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", tmpl, err)
		}
		files = append(files, fileEntry{Name: path, Data: data})
	}

	// Per-slide files
	for i, slide := range p.Slides {
		slideNum := i + 1
		
		// Determine slide template based on index (first slide is title slide)
		var slideTemplate, relsTemplate string
		if i == 0 {
			slideTemplate = "templates/ppt/slides/title_slide.xml.tmpl"
			relsTemplate = "templates/ppt/slides/_rels/title_slide.rels.tmpl"
		} else {
			slideTemplate = "templates/ppt/slides/content_slide.xml.tmpl"
			relsTemplate = "templates/ppt/slides/_rels/content_slide.rels.tmpl"
		}

		// Slide XML
		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNum)
		slideData, err := renderTemplate(slideTemplate, slide)
		if err != nil {
			return nil, fmt.Errorf("rendering slide %d: %w", slideNum, err)
		}
		files = append(files, fileEntry{Name: slidePath, Data: slideData})

		// Slide relationships
		relsPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNum)
		relsData, err := renderTemplate(relsTemplate, map[string]any{
			"Index": slideNum,
		})
		if err != nil {
			return nil, fmt.Errorf("rendering slide %d rels: %w", slideNum, err)
		}
		files = append(files, fileEntry{Name: relsPath, Data: relsData})

		// Notes slide
		notesPath := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", slideNum)
		notesData, err := renderTemplate("templates/ppt/notesSlides/notesSlide.xml.tmpl", slide)
		if err != nil {
			return nil, fmt.Errorf("rendering notes slide %d: %w", slideNum, err)
		}
		files = append(files, fileEntry{Name: notesPath, Data: notesData})

		// Notes slide relationships
		notesRelsPath := fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", slideNum)
		notesRelsData, err := renderTemplate("templates/ppt/notesSlides/_rels/notesSlide.xml.rels.tmpl", map[string]any{
			"Index": slideNum,
		})
		if err != nil {
			return nil, fmt.Errorf("rendering notes slide %d rels: %w", slideNum, err)
		}
		files = append(files, fileEntry{Name: notesRelsPath, Data: notesRelsData})
	}

	return files, nil
}

func renderTemplate(path string, data any) ([]byte, error) {
	tmplBytes, err := templatesFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", path, err)
	}

	tmpl, err := template.New(filepath.Base(path)).Funcs(funcMap).Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", path, err)
	}

	return buf.Bytes(), nil
}


