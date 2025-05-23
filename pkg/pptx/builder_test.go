package pptx_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/napcs/opml2pptx/pkg/opml"
	"github.com/napcs/opml2pptx/pkg/pptx"
)

func createTestPresentation() *opml.Presentation {
	return &opml.Presentation{
		Title: "Test Presentation & <Special> Characters",
		Slides: []opml.Slide{
			{
				Title:   "Test Presentation & <Special> Characters",
				Bullets: []string{},
				Notes:   "",
			},
			{
				Title:   "Content Slide with <XML> & \"Quotes\"",
				Bullets: []string{"Bullet & Point 1", "Bullet <Two>", "Point \"Three\""},
				Notes:   "These are speaker notes with <special> & characters.",
			},
			{
				Title:   "Another Content Slide",
				Bullets: []string{"Simple bullet", "Another point"},
				Notes:   "",
			},
		},
	}
}

func TestBuildPPTX(t *testing.T) {
	pres := createTestPresentation()
	
	// Create temporary file
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("PPTX file was not created")
	}
}

func TestPPTXStructure(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Open and verify ZIP structure
	reader, err := zip.OpenReader(tmpFile)
	if err != nil {
		t.Fatalf("Failed to open PPTX as ZIP: %v", err)
	}
	defer reader.Close()
	
	// Track found files
	foundFiles := make(map[string]bool)
	for _, file := range reader.File {
		foundFiles[file.Name] = true
	}
	
	// Required core files
	requiredFiles := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"ppt/presentation.xml",
		"ppt/_rels/presentation.xml.rels",
		"ppt/slideMasters/slideMaster1.xml",
		"ppt/theme/theme1.xml",
	}
	
	for _, required := range requiredFiles {
		if !foundFiles[required] {
			t.Errorf("Required file missing: %s", required)
		}
	}
	
	// Check slide files
	expectedSlides := 3
	for i := 1; i <= expectedSlides; i++ {
		slideFile := fmt.Sprintf("ppt/slides/slide%d.xml", i)
		notesFile := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", i)
		
		if !foundFiles[slideFile] {
			t.Errorf("Missing slide file: %s", slideFile)
		}
		if !foundFiles[notesFile] {
			t.Errorf("Missing notes file: %s", notesFile)
		}
	}
}

func TestTitleSlideIsFirst(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	content := extractFileFromPPTX(t, tmpFile, "ppt/slides/slide1.xml")
	
	// Title slide should contain the presentation title
	if !strings.Contains(content, "Test Presentation") {
		t.Error("First slide should contain presentation title")
	}
	
	// Title slide should not contain bullet points
	if strings.Contains(content, "<a:buChar") {
		t.Error("Title slide should not contain bullet points")
	}
}

func TestContentSlidesHaveBullets(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Check second slide (first content slide)
	content := extractFileFromPPTX(t, tmpFile, "ppt/slides/slide2.xml")
	
	// Should contain bullet points
	if !strings.Contains(content, "Bullet &amp; Point 1") {
		t.Error("Content slide should contain first bullet point")
	}
	if !strings.Contains(content, "Bullet &lt;Two&gt;") {
		t.Error("Content slide should contain second bullet point")
	}
	if !strings.Contains(content, "Point &#34;Three&#34;") {
		t.Error("Content slide should contain third bullet point")
	}
}

func TestSpecialCharacterEscaping(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Check title slide
	titleContent := extractFileFromPPTX(t, tmpFile, "ppt/slides/slide1.xml")
	
	// XML special characters should be escaped
	if !strings.Contains(titleContent, "&amp;") {
		t.Error("Ampersand should be escaped as &amp;")
	}
	if !strings.Contains(titleContent, "&lt;Special&gt;") {
		t.Error("Angle brackets should be escaped")
	}
	
	// Check content slide
	contentSlide := extractFileFromPPTX(t, tmpFile, "ppt/slides/slide2.xml")
	
	// Should not contain unescaped special characters
	if strings.Contains(contentSlide, "<XML>") || strings.Contains(contentSlide, "\"Quotes\"") {
		t.Error("Special characters in content should be escaped")
	}
	
	// Should contain properly escaped versions
	if !strings.Contains(contentSlide, "&lt;XML&gt;") {
		t.Error("XML tags should be escaped in slide content")
	}
	if !strings.Contains(contentSlide, "&#34;Quotes&#34;") {
		t.Error("Quotes should be escaped in slide content")
	}
}

func TestNotesSlideContent(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Check notes for slide 2 (which has notes)
	notesContent := extractFileFromPPTX(t, tmpFile, "ppt/notesSlides/notesSlide2.xml")
	
	// Should contain the speaker notes
	if !strings.Contains(notesContent, "These are speaker notes") {
		t.Error("Notes slide should contain speaker notes text")
	}
	
	// Special characters in notes should be escaped
	if !strings.Contains(notesContent, "&lt;special&gt;") {
		t.Error("Special characters in notes should be escaped")
	}
	if !strings.Contains(notesContent, "&amp;") {
		t.Error("Ampersand in notes should be escaped")
	}
}

func TestEmptyNotesSlide(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Check notes for slide 3 (which has empty notes)
	notesContent := extractFileFromPPTX(t, tmpFile, "ppt/notesSlides/notesSlide3.xml")
	
	// Should still be a valid XML file but with no meaningful text content
	if strings.Contains(notesContent, "Simple bullet") {
		t.Error("Notes slide should not contain slide bullet content")
	}
}

func TestTemplateExecution(t *testing.T) {
	pres := createTestPresentation()
	tmpFile := filepath.Join(t.TempDir(), "test.pptx")
	
	err := pptx.BuildPPTX(pres, tmpFile)
	if err != nil {
		t.Fatalf("BuildPPTX failed: %v", err)
	}
	
	// Verify that templates were executed properly by checking for template placeholders
	contentTypes := extractFileFromPPTX(t, tmpFile, "[Content_Types].xml")
	
	// Should not contain template syntax
	if strings.Contains(contentTypes, "{{") || strings.Contains(contentTypes, "}}") {
		t.Error("Content types file contains unprocessed template syntax")
	}
	
	// Should contain expected MIME types
	if !strings.Contains(contentTypes, "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml") {
		t.Error("Content types should contain presentation MIME type")
	}
}

// Helper function to extract a file from the PPTX ZIP
func extractFileFromPPTX(t *testing.T, pptxPath, fileName string) string {
	reader, err := zip.OpenReader(pptxPath)
	if err != nil {
		t.Fatalf("Failed to open PPTX: %v", err)
	}
	defer reader.Close()
	
	for _, file := range reader.File {
		if file.Name == fileName {
			rc, err := file.Open()
			if err != nil {
				t.Fatalf("Failed to open file %s: %v", fileName, err)
			}
			defer rc.Close()
			
			var buf bytes.Buffer
			_, err = buf.ReadFrom(rc)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", fileName, err)
			}
			
			return buf.String()
		}
	}
	
	t.Fatalf("File %s not found in PPTX", fileName)
	return ""
}