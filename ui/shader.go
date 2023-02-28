package ui

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/mokiat/lacking/ui"
)

//go:embed shaders/*
var sources embed.FS

var rootTemplate = template.Must(template.
	New("root").
	Delims("/*", "*/").
	ParseFS(sources, "shaders/*.glsl"),
)

func find(name string) *template.Template {
	result := rootTemplate.Lookup(name)
	if result == nil {
		panic(fmt.Errorf("template %q not found", name))
	}
	return result
}

var buffer = new(bytes.Buffer)

func runTemplate(tmpl *template.Template, data any) string {
	buffer.Reset()
	if err := tmpl.Execute(buffer, data); err != nil {
		panic(fmt.Errorf("template exec error: %w", err))
	}
	return buffer.String()
}

var (
	tmplShadedShapeVertexShader   = find("shaded_shape.vert.glsl")
	tmplShadedShapeFragmentShader = find("shaded_shape.frag.glsl")

	tmplBlankShapeVertexShader   = find("blank_shape.vert.glsl")
	tmplBlankShapeFragmentShader = find("blank_shape.frag.glsl")

	tmplContourVertexShader   = find("contour.vert.glsl")
	tmplContourFragmentShader = find("contour.frag.glsl")

	tmplTextVertexShader   = find("text.vert.glsl")
	tmplTextFragmentShader = find("text.frag.glsl")
)

func NewShaderCollection() ui.ShaderCollection {
	return ui.ShaderCollection{
		ShapeShadedSet: newShadedShapeShaderSet,
		ShapeBlankSet:  newBlankShapeShaderSet,
		ContourSet:     newContourShaderSet,
		TextSet:        newTextShaderSet,
	}
}

func newShadedShapeShaderSet() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader:   runTemplate(tmplShadedShapeVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplShadedShapeFragmentShader, struct{}{}),
	}
}

func newBlankShapeShaderSet() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader:   runTemplate(tmplBlankShapeVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplBlankShapeFragmentShader, struct{}{}),
	}
}

func newContourShaderSet() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader:   runTemplate(tmplContourVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplContourFragmentShader, struct{}{}),
	}
}

func newTextShaderSet() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader:   runTemplate(tmplTextVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplTextFragmentShader, struct{}{}),
	}
}
