/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_light.glsl"*/

uniform vec3 lightIntensityIn = vec3(1.0, 1.0, 1.0);
uniform float lightRangeIn = 4.0;

/*template "math.glsl"*/

struct fresnelInput {
	vec3 reflectanceF0;
	vec3 halfDirection;
	vec3 lightDirection;
};

vec3 calculateFresnel(fresnelInput i) {
	float halfLightDot = clamp(dot(i.halfDirection, i.lightDirection), 0.0, 1.0);
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
	return (reflectedHDR + refractedHDR) * s.lightIntensity * clamp(dot(s.normal, s.lightDirection), 0.0, 1.0);
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

	vec3 lightDirection = lightMatrixIn[3].xyz - worldPosition;
	float lightDistance = length(lightDirection);
	float distAttenuation = getCappedDistanceAttenuation(lightDistance, lightRangeIn);

	vec3 hdr = calculateDirectionalHDR(directionalSetup(
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		normalize(lightDirection),
		normal,
		lightIntensityIn
	));
	fbColor0Out = vec4(hdr * distAttenuation, 1.0);
}