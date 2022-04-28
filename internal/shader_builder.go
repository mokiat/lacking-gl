package internal

import "strings"

func NewShaderSourceBuilder(template string) *ShaderSourceBuilder {
	return &ShaderSourceBuilder{
		version:  "460",
		features: []string{},
		template: template,
	}
}

type ShaderSourceBuilder struct {
	version  string
	features []string
	template string
}

func (b *ShaderSourceBuilder) SetVersion(version string) {
	b.version = version
}

func (b *ShaderSourceBuilder) AddFeature(feature string) {
	b.features = append(b.features, feature)
}

func (b *ShaderSourceBuilder) Build() string {
	var builder strings.Builder
	builder.WriteString("#version ")
	builder.WriteString(b.version)
	builder.WriteRune('\n')
	for _, feature := range b.features {
		builder.WriteString("#define ")
		builder.WriteString(feature)
		builder.WriteRune('\n')
	}
	builder.WriteString(b.template)
	builder.WriteRune('\n')
	return builder.String()
}
