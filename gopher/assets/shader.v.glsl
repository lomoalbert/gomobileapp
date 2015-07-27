#version 100

uniform mat4 projection;
uniform mat4 view;
uniform mat4 modelx;
uniform mat4 modely;

attribute vec3 vertCoord;


attribute vec4 color;
varying vec4 vColor;

void main() {
    gl_Position = projection * view * modelx* modely * vec4(vertCoord, 1);

	vColor = color;
}
