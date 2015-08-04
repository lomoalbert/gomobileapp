#version 100
precision lowp float;

uniform sampler2D u_texture;

varying float v_intensity;
varying vec2 v_texCoord;

void main(void)
{
	gl_FragColor = vec4(0.8,0.8,0.8,1) * v_intensity;
}