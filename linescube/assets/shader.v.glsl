#version 100

uniform mat4 projection;
uniform mat4 view;
uniform mat4 modelx;
uniform mat4 modely;

attribute vec3 vertCoord;


void main() {
    gl_Position = projection * view * modely* modelx * vec4(vertCoord, 1);
}
