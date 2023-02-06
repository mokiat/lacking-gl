/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;
layout(binding = 4) uniform samplerCube reflectionTextureIn;
layout(binding = 5) uniform samplerCube refractionTextureIn;

/*template "ubo_camera.glsl"*/

/*template "math.glsl"*/

/*template "lighting.glsl"*/

struct ambientFresnelInput {
	vec3 reflectanceF0;
	vec3 normal;
	vec3 viewDirection;
	float roughness;
};

vec3 calculateAmbientFresnel(ambientFresnelInput i) {
	float normViewDot = clamp(dot(i.normal, i.viewDirection), 0.0, 1.0);
	return i.reflectanceF0 + (max(vec3(1.0 - i.roughness), i.reflectanceF0) - i.reflectanceF0) * pow(1.0 - normViewDot, 5);
}

struct ambientSetup {
	samplerCube reflectionTexture;
	samplerCube refractionTexture;
	float roughness;
	vec3 reflectedColor;
	vec3 refractedColor;
	vec3 viewDirection;
	vec3 normal;
};

vec3 calculateAmbientHDR(ambientSetup s) {
	vec3 fresnel = calculateAmbientFresnel(ambientFresnelInput(
		s.reflectedColor,
		s.normal,
		s.viewDirection,
		s.roughness
	));

	vec3 lightDirection = reflect(s.viewDirection, s.normal);
	vec3 reflectedLightIntensity = pow(mix(
			pow(texture(s.refractionTexture, lightDirection) / pi, vec4(0.25)),
			pow(texture(s.reflectionTexture, lightDirection), vec4(0.25)),
			pow(1.0 - s.roughness, 4)
		), vec4(4)).xyz;
	float geometry = calculateGeometry(geometryInput(
		s.roughness
	));
	vec3 reflectedHDR = fresnel * s.reflectedColor * reflectedLightIntensity * geometry;

	vec3 refractedLightIntensity = texture(s.refractionTexture, -s.normal).xyz;
	vec3 refractedHDR = (vec3(1.0) - fresnel) * s.refractedColor * refractedLightIntensity / pi;

	return (reflectedHDR + refractedHDR);
}

void main()
{
	vec2 screenCoord = getScreenUVCoords(viewportIn);
	vec3 ndcPosition = getScreenNDC(screenCoord, fbDepthTextureIn);
	vec3 viewPosition = getViewCoords(ndcPosition, projectionMatrixIn);
	vec3 worldPosition = getWorldCoords(viewPosition, cameraMatrixIn);
	vec3 cameraPosition = cameraMatrixIn[3].xyz;

	vec4 albedoMetalness = texture(fbColor0TextureIn, screenCoord);
	vec4 normalRoughness = texture(fbColor1TextureIn, screenCoord);
	vec3 baseColor = albedoMetalness.xyz;
	vec3 normal = normalize(normalRoughness.xyz);
	float metalness = albedoMetalness.w;
	float roughness = normalRoughness.w;

	vec3 refractedColor = baseColor * (1.0 - metalness);
	vec3 reflectedColor = mix(vec3(0.02), baseColor, metalness);

	vec3 hdr = calculateAmbientHDR(ambientSetup(
		reflectionTextureIn,
		refractionTextureIn,
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		normal
	));
	fbColor0Out = vec4(hdr, 1.0);
}