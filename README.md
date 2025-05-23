# opml2pptx

A command-line tool that converts OPML (Outline Processor Markup Language) files into PowerPoint presentations (.pptx). Transform your hierarchical outlines into structured slides with titles, bullet points, and speaker notes.

## Features

- **Hierarchical Conversion**: Converts OPML outlines to PowerPoint slides with unlimited nesting levels
- **Smart Structure**: Root outline becomes title slide, child outlines become content slides
- **Sub-bullets**: Supports multi-level bullet points with proper indentation and visual hierarchy
- **Speaker Notes**: Preserves `_note` attributes as PowerPoint speaker notes
- **Cross-Platform**: Builds for macOS, Linux, and Windows

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the releases page.

#### macOS users

macOS will flag this app if you download it.

Copy the release to your Path

```bash
mv ~/downloads/opml2pptx-0.0.1-darwin-arm64 ~/bin/opml2pptx
```

Then make it executable and remove the quarantine bit:

```bash
chmod +x opml2pptx
xattr -d com.apple.quarantine opml2pptx
```

If you are uncomfortable doing this, review the source code and build your own binary which will run fine.

### Build from Source

Requirements:
- Go 1.19 or later

```bash
git clone https://github.com/napcs/opml2pptx.git
cd opml2pptx
make all
```

Built binaries will be available in the `dist/` directory.

## Usage

Basic Usage

```bash
opml2pptx -input input.opml -output output.pptx
```

Command Line Options

- `-input`: Path to the input OPML file (required)
- `-output`: Path for the output PowerPoint file (required)
- `-version`: Show version information

Example:

```bash
./opml2pptx -input my-outline.opml -output presentation.pptx
```

## OPML Structure

The tool expects OPML files with the following structure:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.0">
  <head>
    <title>My Presentation</title>
  </head>
  <body>
    <outline text="Presentation Title">
      <outline text="First Slide Title">
        <outline text="First bullet point">
          <outline text="Sub-bullet point">
            <outline text="Sub-sub-bullet point"></outline>
          </outline>
        </outline>
        <outline text="Second bullet point"></outline>
      </outline>
      <outline text="Second Slide Title" _note="Speaker notes go here">
        <outline text="Another bullet point"></outline>
      </outline>
    </outline>
  </body>
</opml>
```

### Conversion Rules

1. **Title Slide**: The root outline becomes the presentation title slide
2. **Content Slides**: Each child of the root becomes a separate content slide
3. **Bullet Points**: Nested outlines within slides become hierarchical bullet points
4. **Speaker Notes**: Use the `_note` attribute to add speaker notes to slides
5. **Visual Hierarchy**: Different bullet characters (-, •, ◦) for different nesting levels

## Development

Build for the current platform:

```bash
go build -o opml2pptx ./cmd/opml2pptx
```

Build for all platforms:

```bash
make all
```

Run tests:

```bash
make test
```

Or directly with Go:

```bash
go test ./...
```

Project Structure

```
opml2pptx/
├── cmd/opml2pptx/          # Main application entry point
├── pkg/opml/               # OPML parsing logic
├── pkg/pptx/               # PowerPoint generation engine
│   └── templates/          # PowerPoint XML templates
├── dist/                   # Build output directory
└── Makefile                # Build automation
```

- **OPML Parser** (`pkg/opml`): Parses OPML files into presentation data structures
- **PPTX Builder** (`pkg/pptx`): Generates PowerPoint files using XML templates
- **Template System**: Uses Go templates to create valid PowerPoint XML structure


## Contributing

This is very early software that works for my needs. But if you find a bug, I'd love your help:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Submit a pull request

## Support

For issues and feature requests, please use the GitHub issue tracker.

## License

MIT License - see [LICENSE](LICENSE) for details.
