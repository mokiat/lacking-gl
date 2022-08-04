layout(location = 0) in vec4 coordIn;
#if defined(USES_BONES)
layout(location = 5) in vec4 weightsIn;
layout(location = 6) in uvec4 jointsIn;
#endif

layout (binding = 3, std140) uniform Light
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
};

#if defined(USES_BONES)
layout (binding = 1, std140) uniform Model
{
	mat4 modelMatrixIn;
	mat4 boneMatrixIn[255];
};
#else
layout (binding = 1, std140) uniform Model
{
	mat4 modelMatrixIn[256];
};
#endif

void main()
{
#if defined(USES_BONES)
	mat4 modelMatrixA = modelMatrixIn * boneMatrixIn[jointsIn.x];
	mat4 modelMatrixB = modelMatrixIn * boneMatrixIn[jointsIn.y];
	mat4 modelMatrixC = modelMatrixIn * boneMatrixIn[jointsIn.z];
	mat4 modelMatrixD = modelMatrixIn * boneMatrixIn[jointsIn.w];
	vec4 worldPosition =
		modelMatrixA * (weightsIn.x * coordIn) +
		modelMatrixB * (weightsIn.y * coordIn) +
		modelMatrixC * (weightsIn.z * coordIn) +
		modelMatrixD * (weightsIn.w * coordIn);
#else
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	vec4 worldPosition = modelMatrix * coordIn;
#endif
  gl_Position = projectionMatrixIn * (viewMatrixIn * worldPosition);
}
