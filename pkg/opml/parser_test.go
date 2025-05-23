package opml_test

import (
	"strings"
	"testing"

	"github.com/napcs/opml2pptx/pkg/opml"
)

const sampleOPML = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>Sample Deck</title>
  </head>
  <body>
    <outline text="Sample Deck">
      <outline text="Slide One">
        <outline text="Bullet A"></outline>
        <outline text="Bullet B"></outline>
      </outline>
      <outline text="Slide Two" _note="Speaker notes here.">
        <outline text="Point 1"></outline>
      </outline>
    </outline>
  </body>
</opml>`

func TestParseOPML(t *testing.T) {
	r := strings.NewReader(sampleOPML)
	pres, err := opml.ParseOPML(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pres.Title != "Sample Deck" {
		t.Errorf("expected title 'Sample Deck', got %q", pres.Title)
	}

	if len(pres.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(pres.Slides))
	}

	slide1 := pres.Slides[0]
	if slide1.Title != "Sample Deck" {
		t.Errorf("expected slide 1 title 'Sample Deck', got %q", slide1.Title)
	}
	if len(slide1.Bullets) != 0 {
		t.Errorf("expected no bullets on slide 1, got %+v", slide1.Bullets)
	}

	slide2 := pres.Slides[1]
	if slide2.Title != "Slide One" {
		t.Errorf("expected slide 2 title 'Slide One', got %q", slide2.Title)
	}
	if len(slide2.Bullets) != 2 ||
		slide2.Bullets[0].Text != "Bullet A" ||
		slide2.Bullets[1].Text != "Bullet B" {
		t.Errorf("unexpected bullets in slide 2: %+v", slide2.Bullets)
	}

	slide3 := pres.Slides[2]
	if slide3.Title != "Slide Two" {
		t.Errorf("expected slide 3 title 'Slide Two', got %q", slide3.Title)
	}
	if len(slide3.Bullets) != 1 || slide3.Bullets[0].Text != "Point 1" {
		t.Errorf("unexpected bullets in slide 3: %+v", slide3.Bullets)
	}
	if slide3.Notes != "Speaker notes here." {
		t.Errorf("expected notes on slide 3 to be 'Speaker notes here.', got %q", slide3.Notes)
	}
}
