#version 100
precision mediump float;
varying vec4 vColor;
uniform sampler2D tex;
uniform bool useuv;
varying vec2 fragTexCoord;

void main() {
	if (useuv){
		gl_FragColor = texture2D(tex, fragTexCoord);
	}else{
		gl_FragColor = vColor;
	}

}