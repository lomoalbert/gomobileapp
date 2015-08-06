#version 100
uniform mat4 u_projectionMatrix;
uniform mat4 u_viewMatrix;
uniform mat4 u_modelMatrix;
uniform mat4 u_modelMatrixy;
uniform mat3 u_normalMatrix;
uniform vec3 u_lightDirection;
uniform mat4 u_lightmatrix;

attribute vec4 a_vertex;
attribute vec3 a_normal;
attribute vec2 a_texCoord;

varying float v_intensity; //强度
varying vec2 v_texCoord;

void main(void)
{
	// Now the normal is in world space, as we pass the light in world space.

    // Intensity is lambert without emissive color. al is the ambient, hard coded light factor.
 	v_intensity = 0.3+0.7*abs(dot(a_normal, u_lightDirection*vec3(u_lightmatrix))); //dot(x , y) = x[0]⋅y[0]+x[1]⋅y[1]+...

	v_texCoord = a_texCoord;

	gl_Position = u_projectionMatrix*u_viewMatrix * u_modelMatrixy * u_modelMatrix * a_vertex;
}