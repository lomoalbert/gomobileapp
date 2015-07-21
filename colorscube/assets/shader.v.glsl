#version 100

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertCoord;


void main() {
    gl_Position = projection * view * model * vec4(vertCoord, 1);
}
