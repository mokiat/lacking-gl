/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;
layout(binding = 4) uniform sampler2DShadow fbShadowTextureIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_light.glsl"*/

uniform vec3 lightIntensityIn = vec3(1.0, 1.0, 1.0);

/*template "math.glsl"*/

struct fresnelInput {
	vec3 reflectanceF0;
	vec3 halfDirection;
	vec3 lightDirection;
};

vec3 calculateFresnel(fresnelInput i) {
	float halfLightDot = clamp(abs(dot(i.halfDirection, i.lightDirection)), 0.0, 1.0);
	return i.reflectanceF0 + (1.0 - i.reflectanceF0) * pow(1.0 - halfLightDot, 5);
}

struct distributionInput {
	float roughness;
	vec3 normal;
	vec3 halfDirection;
};

float calculateDistribution(distributionInput i) {
	float sqrRough = i.roughness * i.roughness;
	float halfNormDot = dot(i.normal, i.halfDirection);
	float denom = halfNormDot * halfNormDot * (sqrRough - 1.0) + 1.0;
	return sqrRough / (pi * denom * denom);
}

struct geometryInput {
	float roughness;
};

float calculateGeometry(geometryInput i) {
	// TODO: Use better model
	return 1.0 / 4.0;
}

struct directionalSetup {
	float roughness;
	vec3 reflectedColor;
	vec3 refractedColor;
	vec3 viewDirection;
	vec3 lightDirection;
	vec3 normal;
	vec3 lightIntensity;
};

vec3 calculateDirectionalHDR(directionalSetup s) {
	vec3 halfDirection = normalize(s.lightDirection + s.viewDirection);
	vec3 fresnel = calculateFresnel(fresnelInput(
		s.reflectedColor,
		halfDirection,
		s.lightDirection
	));
	float distributionFactor = calculateDistribution(distributionInput(
		s.roughness,
		s.normal,
		halfDirection
	));
	float geometryFactor = calculateGeometry(geometryInput(
		s.roughness
	));
	vec3 reflectedHDR = fresnel * distributionFactor * geometryFactor;
	vec3 refractedHDR = (vec3(1.0) - fresnel) * s.refractedColor / pi;
	return (reflectedHDR + refractedHDR) * s.lightIntensity * clamp(abs(dot(s.normal, s.lightDirection)), 0.0, 1.0);
}

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
