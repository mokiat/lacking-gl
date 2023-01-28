/*template "version.glsl"*/

layout(location = 0) out vec4 fbColor0Out;

layout(binding = 0) uniform sampler2D fbColor0TextureIn;
layout(binding = 1) uniform sampler2D fbColor1TextureIn;
layout(binding = 3) uniform sampler2D fbDepthTextureIn;
layout(binding = 4) uniform sampler2DShadow fbShadowTextureIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_light.glsl"*/

uniform vec3 lightDirectionIn;
uniform vec3 lightIntensityIn = vec3(1.0, 1.0, 1.0);

noperspective in vec2 texCoordInOut;

const float pi = 3.141592;

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

void main()
{
	vec3 ndcPosition = vec3(
		(texCoordInOut.x - 0.5) * 2.0,
		(texCoordInOut.y - 0.5) * 2.0,
		texture(fbDepthTextureIn, texCoordInOut).x * 2.0 - 1.0
	);
	vec3 clipPosition = vec3(
		ndcPosition.x / projectionMatrixIn[0][0],
		ndcPosition.y / projectionMatrixIn[1][1],
		-1.0
	);
	vec3 viewPosition = clipPosition * projectionMatrixIn[3][2] / (projectionMatrixIn[2][2] + ndcPosition.z);
	vec3 worldPosition = (cameraMatrixIn * vec4(viewPosition, 1.0)).xyz;
	vec3 cameraPosition = cameraMatrixIn[3].xyz;

	vec4 albedoMetalness = texture(fbColor0TextureIn, texCoordInOut);
	vec4 normalRoughness = texture(fbColor1TextureIn, texCoordInOut);
	vec3 baseColor = albedoMetalness.xyz;
	vec3 normal = normalize(normalRoughness.xyz);
	float metalness = albedoMetalness.w;
	float roughness = normalRoughness.w;

	vec3 refractedColor = baseColor * (1.0 - metalness);
	vec3 reflectedColor = mix(vec3(0.02), baseColor, metalness);

	vec3 hdr = calculateDirectionalHDR(directionalSetup(
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		normalize(lightDirectionIn),
		normal,
		lightIntensityIn
	));

	vec4 lightPosition = lightProjectionMatrixIn * lightViewMatrixIn * vec4(worldPosition, 1.0);
	float directness = clamp(abs(dot(normal, normalize(lightDirectionIn))), 0.0, 1.0);
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