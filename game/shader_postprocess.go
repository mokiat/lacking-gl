package game

import (
	"fmt"

	"github.com/mokiat/lacking-gl/internal"
	"github.com/mokiat/lacking/game/graphics/renderapi/plugin"
)

func newPostprocessingShaderSet(mapping plugin.ToneMapping) plugin.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(tonePostprocessingVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(tonePostprocessingFragmentShader)
	switch mapping {
	case plugin.ReinhardToneMapping:
		fsBuilder.AddFeature("MODE_REINHARD")
	case plugin.ExponentialToneMapping:
		fsBuilder.AddFeature("MODE_EXPONENTIAL")
	default:
		panic(fmt.Errorf("unknown tone mapping mode: %s", mapping))
	}
	return plugin.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const tonePostprocessingVertexShader = `
layout(location = 0) in vec2 coordIn;

noperspective out vec2 texCoordInOut;

void main()
{
	texCoordInOut = (coordIn + 1.0) / 2.0;
	gl_Position = vec4(coordIn, 0.0, 1.0);
}
`

const tonePostprocessingFragmentShader = `
layout(location = 0) out vec4 fbColor0Out;

uniform sampler2D fbColor0TextureIn;
uniform float exposureIn = 1.0;

noperspective in vec2 texCoordInOut;

void main()
{
	vec3 hdr = texture(fbColor0TextureIn, texCoordInOut).xyz;
	vec3 exposedHDR = hdr * exposureIn;
	#if defined(MODE_REINHARD)
	vec3 ldr = exposedHDR / (exposedHDR + vec3(1.0));
	#endif
	#if defined(MODE_EXPONENTIAL)
	vec3 ldr = vec3(1.0) - exp2(-exposedHDR);
	#endif
	fbColor0Out = vec4(ldr, 1.0);
}
`
