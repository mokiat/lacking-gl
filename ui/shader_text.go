package ui

import (
	"github.com/mokiat/lacking-gl/internal"
	"github.com/mokiat/lacking/ui"
)

func newTextShaders() ui.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(textMaterialVertexShaderTemplate)
	fsBuilder := internal.NewShaderSourceBuilder(textMaterialFragmentShaderTemplate)
	return ui.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const textMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;
layout(location = 1) in vec2 texCoordIn;

uniform mat4 transformMatrixIn;
uniform mat4 projectionMatrixIn;
uniform mat4 clipMatrixIn;

out gl_PerVertex
{
  vec4 gl_Position;
  float gl_ClipDistance[4];
};

noperspective out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	vec4 clipValues = clipMatrixIn * screenPosition;
	gl_ClipDistance[0] = clipValues.x;
	gl_ClipDistance[1] = clipValues.y;
	gl_ClipDistance[2] = clipValues.z;
	gl_ClipDistance[3] = clipValues.w;
	
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const textMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

uniform sampler2D textureIn;
uniform vec4 colorIn = vec4(1.0, 1.0, 1.0, 1.0);

noperspective in vec2 texCoordInOut;

void main()
{
	float amount = texture(textureIn, texCoordInOut).x;
	fragmentColor = vec4(amount) * colorIn;
}
`
