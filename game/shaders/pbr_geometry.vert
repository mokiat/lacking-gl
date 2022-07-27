layout(location = 0) in vec4 coordIn;
layout(location = 1) in vec3 normalIn;
#if defined(USES_TEX_COORD0)
layout(location = 3) in vec2 texCoordIn;
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
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	normalInOut = inverse(transpose(mat3(modelMatrix))) * normalIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * (modelMatrix * coordIn));
}