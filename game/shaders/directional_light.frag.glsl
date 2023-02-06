/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;
layout(binding = 4) uniform sampler2DShadow fbShadowTextureIn;

uniform vec3 lightIntensityIn = vec3(1.0, 1.0, 1.0);

/*template "ubo_camera.glsl"*/

/*template "ubo_light.glsl"*/

/*template "math.glsl"*/

/*template "lighting.glsl"*/

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

	vec3 lightDirection = normalize(lightMatrixIn[3].xyz);

	vec3 hdr = calculateDirectionalHDR(directionalSetup(
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		lightDirection,
		normal,
		lightIntensityIn
	));

	vec4 lightPosition = lightProjectionMatrixIn * lightViewMatrixIn * vec4(worldPosition, 1.0);
	float directness = clamp(abs(dot(normal, lightDirection)), 0.0, 1.0);
	lightPosition.xyz = lightPosition.xyz * 0.5 + 0.5;
	lightPosition.z /= lightPosition.w;
	lightPosition.z -= 0.0005;

	vec2 shift = 1.0 / vec2(textureSize(fbShadowTextureIn, 0));

	const vec3[] shifts = {
		vec3(0.0, 0.0, 0.0),
		vec3(-1.0, 0.0, 0.0),
		vec3(1.0, 0.0, 0.0),
		vec3(0.0, -1.0, 0.0),
		vec3(0.0, 1.0, 0.0),
		vec3(-1.0, -1.0, 0.0),
		vec3(1.0, -1.0, 0.0),
		vec3(-1.0, 1.0, 0.0),
		vec3(1.0, 1.0, 0.0),
	};

	float amount = 0.0;
	for (int i = 0; i < 9; i++) {
		float probability = texture(fbShadowTextureIn, lightPosition.xyz + shifts[i] * vec3(shift.x, shift.y, 1.0));
		amount = max(amount, probability);
	}

	float factor = (clamp(directness, 0.3, 0.5) - 0.3) / 0.2;
	amount = amount * factor;

	fbColor0Out = vec4(hdr * amount, 1.0);
}
