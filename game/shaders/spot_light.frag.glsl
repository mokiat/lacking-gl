/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_light.glsl"*/

uniform vec3 lightIntensityIn = vec3(1.0, 1.0, 1.0);
uniform float lightRangeIn = 4.0;
uniform float lightOuterCosIn = 0.3;
uniform float lightInnerCosIn = 0.3;

const float pi = 3.141592;

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

void main()
{
	vec2 screenCoord = vec2(
		gl_FragCoord.x / viewportIn.z,
		gl_FragCoord.y / viewportIn.w
	);

	vec3 ndcPosition = vec3(
		(screenCoord.x - 0.5) * 2.0,
		(screenCoord.y - 0.5) * 2.0,
		texture(fbDepthTextureIn, screenCoord).x * 2.0 - 1.0
	);
	vec3 clipPosition = vec3(
		ndcPosition.x / projectionMatrixIn[0][0],
		ndcPosition.y / projectionMatrixIn[1][1],
		-1.0
	);
	vec3 viewPosition = clipPosition * projectionMatrixIn[3][2] / (projectionMatrixIn[2][2] + ndcPosition.z);
	vec3 worldPosition = (cameraMatrixIn * vec4(viewPosition, 1.0)).xyz;
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
	float lightDistanceSqr = dot(lightDirection, lightDirection);

	float attenuation = max(min(1.0 - sqrt(lightDistanceSqr) / lightRangeIn, 1.0), 0.0) / (1.0 + lightDistanceSqr);

	float coneSoftness = (0.5 +  0.5 * step(lightInnerCosIn, dot(normalize(lightDirection), normalize(lightMatrixIn[1].xyz))));
	float coneAttentuation = step(lightOuterCosIn, dot(normalize(lightDirection), normalize(lightMatrixIn[1].xyz)));

	vec3 hdr = calculateDirectionalHDR(directionalSetup(
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		normalize(lightDirection),
		normal,
		lightIntensityIn
	));
	fbColor0Out = vec4(hdr * attenuation * coneAttentuation * coneSoftness, 1.0);
}