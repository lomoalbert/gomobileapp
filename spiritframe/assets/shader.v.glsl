#version 100

uniform mat4 projection;
uniform mat4 view;
uniform mat4 modelx;
uniform mat4 modely;

attribute vec3 vertCoord;
attribute vec2 vertTexCoord;
varying vec2 fragTexCoord;

uniform vec4 color;
varying vec4 vColor;
void main() {
	vColor = color;
	fragTexCoord = vertTexCoord;
    gl_Position = projection * view * modely* modelx * vec4(vertCoord, 1);

}
