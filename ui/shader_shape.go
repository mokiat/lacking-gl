package ui

import (
	"github.com/mokiat/lacking-gl/internal"
	"github.com/mokiat/lacking/ui"
)

func newShapeShaders() ui.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(shapeMaterialVertexShaderTemplate)
	fsBuilder := internal.NewShaderSourceBuilder(shapeMaterialFragmentShaderTemplate)
	return ui.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

func newShapeBlankShaders() ui.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(shapeBlankMaterialVertexShaderTemplate)
	fsBuilder := internal.NewShaderSourceBuilder(shapeBlankMaterialFragmentShaderTemplate)
	return ui.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const shapeMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;

uniform mat4 transformMatrixIn;
uniform mat4 textureTransformMatrixIn;
uniform mat4 projectionMatrixIn;
uniform vec4 clipDistancesIn;

noperspective out vec2 texCoordInOut;

out gl_PerVertex
{
  vec4 gl_Position;
  float gl_ClipDistance[4];
};

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	texCoordInOut = (textureTransformMatrixIn * vec4(positionIn, 0.0, 1.0)).xy;
	gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

uniform sampler2D textureIn;
uniform vec4 colorIn = vec4(1.0, 1.0, 1.0, 1.0);

noperspective in vec2 texCoordInOut;

void main()
{
	fragmentColor = texture(textureIn, texCoordInOut) * colorIn;
}
`

const shapeBlankMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;

uniform mat4 transformMatrixIn;
uniform mat4 projectionMatrixIn;
uniform vec4 clipDistancesIn;

out gl_PerVertex
{
  vec4 gl_Position;
  float gl_ClipDistance[4];
};

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeBlankMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

void main()
{
	fragmentColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`
