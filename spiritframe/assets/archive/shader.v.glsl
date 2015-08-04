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

#version 400
layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;

out vec4 Position;
out vec3 Normal;


uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 ProjectionMatrix;
uniform mat4 MVP;


void getEyeSpace(out vec3 norm, out vec4 position)
{
    norm =  normalize(NormalMatrix * VertexNormal);
    position = ModelViewMatrix * vec4(VertexPosition, 1.0);
}


void main()
{
    getEyeSpace(Normal, Position);
    gl_Position = MVP * vec4( VertexPosition, 1.0);
}