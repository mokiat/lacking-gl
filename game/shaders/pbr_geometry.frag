layout(location = 0) out vec4 fbColor0Out;
layout(location = 1) out vec4 fbColor1Out;

#if defined(USES_ALBEDO_TEXTURE)
layout(binding = 0) uniform sampler2D albedoTwoDTextureIn;
#endif
uniform vec4 albedoColorIn = vec4(0.5, 0.0, 0.5, 1.0);

uniform float metalnessIn = 0.0;
uniform float roughnessIn = 0.8;
uniform float alphaThresholdIn = 0.5;

smooth in vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth in vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_ALBEDO_TEXTURE) && defined(USES_TEX_COORD0)
	vec4 color = texture(albedoTwoDTextureIn, texCoordInOut);
	if (color.a < alphaThresholdIn) {
		discard;
	}
#else
	vec4 color = albedoColorIn;
#endif

	fbColor0Out = vec4(color.xyz, metalnessIn);
	fbColor1Out = vec4(normalize(normalInOut), roughnessIn);
}
