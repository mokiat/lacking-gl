layout (binding = 2, std140) uniform Material
{
	vec4 albedoColorIn;
	float alphaThresholdIn;
	float normalScaleIn;
	float metallicIn;
	float roughnessIn;
};