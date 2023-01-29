// getScreenUVCoords returns the coordinates on the screen as though there is
// a UV mapping on top (meaning {0.0, 0.0} bottom left and {1.0, 1.0} top right).
vec2 getScreenUVCoords(vec4 viewport)
{
	return (gl_FragCoord.xy - viewport.xy) / viewport.zw;
}

// getScreenNDC converts screen UV coordinates to NDC
// (Normalized Device Coordinates).
vec3 getScreenNDC(vec2 uv, sampler2D depthTexture)
{
  vec3 nonNormized = vec3(uv.x, uv.y, texture(depthTexture, uv).x);
  return nonNormized * 2.0 - vec3(1.0);
}

// getViewCoords converst the NDC coords into view coordinates.
vec3 getViewCoords(vec3 ndc, mat4 projectionMatrix)
{
  vec3 clipCoords = vec3(
		ndc.x / projectionMatrix[0][0],
		ndc.y / projectionMatrix[1][1],
		-1.0
	);
  float scale = projectionMatrix[3][2] / (projectionMatrix[2][2] + ndc.z);
  return clipCoords * scale;
}

// getWorldCoords converts the specified view coords into world coordinates
// depending on the camera positioning.
vec3 getWorldCoords(vec3 viewCoords, mat4 cameraMatrix)
{
  return (cameraMatrix * vec4(viewCoords, 1.0)).xyz;
}

// getCappedDistanceAttenuation calculates the attenuation depending on the
// distance with an upper bound on the maximum distance.
float getCappedDistanceAttenuation(float dist, float maxDist)
{
  float sqrDist = dist * dist;
  float gradient = 1.0 - dist / maxDist;
  return clamp(gradient, 0.0, 1.0) / (1.0 + sqrDist);
}

// getConeAttenuation calculates the attenuation for a cone-shaped light
// source depending on the light direction.
float getConeAttenuation(float angle, float outerAngle, float innerAngle) {
  float hardAttenuation = 1.0 - step(outerAngle, angle);
  float softAttenuation = clamp((outerAngle - angle) / (outerAngle - innerAngle + 0.01), 0.0, 1.0);
  return hardAttenuation * softAttenuation;
}
