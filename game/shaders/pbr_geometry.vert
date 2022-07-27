layout(location = 0) in vec4 coordIn;
layout(location = 1) in vec3 normalIn;
#if defined(USES_TEX_COORD0)
layout(location = 3) in vec2 texCoordIn;
#endif
#if defined(USES_BONES)
layout(location = 5) in vec4 weightsIn;
layout(location = 6) in uvec4 jointsIn;
#endif

layout (binding = 0, std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

layout (binding = 1, std140) uniform Model
{
	mat4 modelMatrixIn[256];
};

smooth out vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth out vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_TEX_COORD0)
	texCoordInOut = texCoordIn;
#endif
#if defined(USES_BONES)
	mat4 modelMatrixA = modelMatrixIn[0] * modelMatrixIn[jointsIn.x+uint(1)];
	mat4 modelMatrixB = modelMatrixIn[0] * modelMatrixIn[jointsIn.y+uint(1)];
	mat4 modelMatrixC = modelMatrixIn[0] * modelMatrixIn[jointsIn.z+uint(1)];
	mat4 modelMatrixD = modelMatrixIn[0] * modelMatrixIn[jointsIn.w+uint(1)];
	vec4 worldPosition =
		modelMatrixA * (weightsIn.x * coordIn) +
		modelMatrixB * (weightsIn.y * coordIn) +
		modelMatrixC * (weightsIn.z * coordIn) +
		modelMatrixD * (weightsIn.w * coordIn);
	vec3 worldNormal =
		inverse(transpose(mat3(modelMatrixA))) * (weightsIn.x * normalIn) +
		inverse(transpose(mat3(modelMatrixB))) * (weightsIn.y * normalIn) +
		inverse(transpose(mat3(modelMatrixC))) * (weightsIn.z * normalIn) +
		inverse(transpose(mat3(modelMatrixD))) * (weightsIn.w * normalIn);
#else
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	vec4 worldPosition = modelMatrix * coordIn;
	vec3 worldNormal = inverse(transpose(mat3(modelMatrix))) * normalIn;
#endif
	normalInOut = worldNormal;
	gl_Position = projectionMatrixIn * (viewMatrixIn * worldPosition);
}