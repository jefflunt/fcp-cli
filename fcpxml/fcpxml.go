package fcpxml

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Params holds all values needed to render an FCPXML document.
type Params struct {
	ProjectName     string
	LibraryPath     string
	TitleCardPath   string
	IntroVOPath     string
	DevlogPath      string
	NormAudioPath   string
	FPS             int
	IntroDurTicks   int64
	DevlogDurTicks  int64
	TransitionTicks int64
}

// templateData is the internal view model passed to the FCPXML template.
type templateData struct {
	Version        string
	FormatName     string
	FrameDuration  string
	LibraryURI     string
	ProjectName    string
	TitleCardURI   string
	IntroVOURI     string
	DevlogURI      string
	NormAudioURI   string
	TotalDur       string
	IntroDur       string
	DevlogDur      string
	TransitionDur  string
	DevlogOffset   string
}

// tickDur formats an FCP frame-count duration as "[n]/[fps]s".
func tickDur(ticks int64, fps int) string {
	return fmt.Sprintf("%d/%ds", ticks, fps)
}

// fileURI converts an absolute path to a file:// URI.
func fileURI(absPath string) string {
	return "file://" + absPath
}

var fcpxmlTemplate = template.Must(template.New("fcpxml").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<fcpxml version="{{.Version}}">
    <resources>
        <format id="r1" name="{{.FormatName}}" frameDuration="{{.FrameDuration}}" width="3840" height="2160"/>
        <asset id="r2" name="title_card" src="{{.TitleCardURI}}" duration="{{.IntroDur}}"/>
        <asset id="r3" name="intro_vo" src="{{.IntroVOURI}}" duration="{{.IntroDur}}"/>
        <asset id="r4" name="devlog_video" src="{{.DevlogURI}}" duration="{{.DevlogDur}}"/>
        <asset id="r5" name="norm_audio" src="{{.NormAudioURI}}" duration="{{.DevlogDur}}"/>
    </resources>
    <library location="{{.LibraryURI}}">
        <event name="{{.ProjectName}}">
            <project name="Assembled_Timeline">
                <sequence format="r1" duration="{{.TotalDur}}">
                    <spine>
                        <video ref="r2" offset="0/1s" duration="{{.IntroDur}}">
                            <audio lane="-1" ref="r3" offset="0/1s" duration="{{.IntroDur}}" role="dialogue"/>
                        </video>
                        <transition name="Cross Dissolve" offset="{{.IntroDur}}" duration="{{.TransitionDur}}"/>
                        <video ref="r4" offset="{{.DevlogOffset}}" duration="{{.DevlogDur}}">
                            <audio lane="-1" ref="r5" offset="{{.DevlogOffset}}" duration="{{.DevlogDur}}" role="dialogue"/>
                        </video>
                    </spine>
                </sequence>
            </project>
        </event>
    </library>
</fcpxml>
`))

// Generate writes an FCPXML v1.10 document to outPath describing the assembled
// devlog timeline defined by p.
func Generate(p Params, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	data := templateData{
		Version:       "1.10",
		FormatName:    fmt.Sprintf("FFVideoFormat3840x2160p%d", p.FPS),
		FrameDuration: fmt.Sprintf("1/%ds", p.FPS),
		LibraryURI:    fileURI(p.LibraryPath),
		ProjectName:   p.ProjectName,
		TitleCardURI:  fileURI(p.TitleCardPath),
		IntroVOURI:    fileURI(p.IntroVOPath),
		DevlogURI:     fileURI(p.DevlogPath),
		NormAudioURI:  fileURI(p.NormAudioPath),
		TotalDur:      tickDur(p.IntroDurTicks+p.TransitionTicks+p.DevlogDurTicks, p.FPS),
		IntroDur:      tickDur(p.IntroDurTicks, p.FPS),
		DevlogDur:     tickDur(p.DevlogDurTicks, p.FPS),
		TransitionDur: tickDur(p.TransitionTicks, p.FPS),
		DevlogOffset:  tickDur(p.IntroDurTicks+p.TransitionTicks, p.FPS),
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating FCPXML file %s: %w", outPath, err)
	}
	defer f.Close()

	if err := fcpxmlTemplate.Execute(f, data); err != nil {
		return fmt.Errorf("rendering FCPXML template: %w", err)
	}

	return nil
}
