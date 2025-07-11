package opml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Outline struct {
	Text  string    `xml:"text,attr"`
	Note  string    `xml:"_note,attr"`
	Child []Outline `xml:"outline"`
}

type OPML struct {
	Head struct {
		Title string `xml:"title"`
	} `xml:"head"`
	Body struct {
		Outlines []Outline `xml:"outline"`
	} `xml:"body"`
}

type Bullet struct {
	Text     string
	Level    int
	Children []Bullet
}

type Slide struct {
	Title   string
	Bullets []Bullet
	Notes   string
}

type Presentation struct {
	Title  string
	Slides []Slide
}

// ParseOPMLFile reads an OPML file and parses it into a Presentation struct.
func ParseOPMLFile(path string) (*Presentation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open OPML file: %w", err)
	}
	defer file.Close()

	return ParseOPML(file)
}

// ParseOPML parses an OPML stream into a Presentation.
func ParseOPML(r io.Reader) (*Presentation, error) {
	var opml OPML
	if err := xml.NewDecoder(r).Decode(&opml); err != nil {
		return nil, fmt.Errorf("failed to parse OPML: %w", err)
	}

	if len(opml.Body.Outlines) == 0 {
		return nil, fmt.Errorf("OPML has no top-level nodes")
	}

	root := opml.Body.Outlines[0]
	presentation := &Presentation{
		Title: root.Text,
	}

	// The first node becomes the title slide
	presentation.Slides = append(presentation.Slides, Slide{
		Title: root.Text,
	})

	for _, child := range root.Child {
		slide := Slide{
			Title: child.Text,
			Notes: collectNote(child),
		}
		slide.Bullets = convertOutlinesToBullets(child.Child, 0)
		presentation.Slides = append(presentation.Slides, slide)
	}

	return presentation, nil
}

func collectNote(node Outline) string {
	if node.Note != "" {
		return node.Note
	}
	for _, child := range node.Child {
		if child.Note != "" {
			return child.Note
		}
	}
	return ""
}

// convertOutlinesToBullets recursively converts OPML outlines to hierarchical bullets
func convertOutlinesToBullets(outlines []Outline, level int) []Bullet {
	var bullets []Bullet
	for _, outline := range outlines {
		bullet := Bullet{
			Text:     outline.Text,
			Level:    level,
			Children: convertOutlinesToBullets(outline.Child, level+1),
		}
		bullets = append(bullets, bullet)
	}
	return bullets
}
