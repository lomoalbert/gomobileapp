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

#version 400

in vec4 Position;
in vec3 Normal;

struct SpotLightInfo{
    vec4 position;
    vec3 direction;
    vec3 intensity;
    float exponent;
    float cutoff;
};

struct MaterialInfo{
    vec3 Ka;
    vec3 Kd;
    vec3 Ks;
    float Shininess;
};

uniform SpotLightInfo Spot;
uniform MaterialInfo Material;

vec3 adsSpotLight(vec4 position, vec3 norm)
{
    vec3 s = normalize(vec3(Spot.position - position));
    float angle = acos(dot(-s, normalize(Spot.direction)));
    float cutoff = radians(clamp(Spot.cutoff, 0.0, 90.0));
    vec3 ambient = Spot.intensity * Material.Ka;

    if(angle < cutoff){
        float spotFactor = pow(dot(-s, normalize(Spot.direction)), Spot.exponent);
        vec3 v = normalize(vec3(-position));
        vec3 h = normalize(v + s);
        return ambient + spotFactor * Spot.intensity * (Material.Kd * max(dot(s, norm),0.0)
              + Material.Ks * pow(max(dot(h,norm), 0.0),Material.Shininess));
    }
    else
    {
        return ambient;
    }

}

void main(void)
{
    gl_FragColor = vec4(adsSpotLight(Position, Normal), 1.0);
    //gl_FragColor = vec4(1.0,1.0,0.5, 1.0);
}