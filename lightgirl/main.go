package main

import (
    "fmt"
    "time"
    "github.com/vzever/wavefront"
    "github.com/go-gl/mathgl/mgl32"
    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event/config"
    "golang.org/x/mobile/event/lifecycle"
    "golang.org/x/mobile/event/paint"
    "golang.org/x/mobile/event/touch"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/gl"
    "encoding/binary"
    "image"
    image_draw "image/draw"
    _ "image/png"
    "io/ioutil"
    "golang.org/x/mobile/asset"
    "golang.org/x/mobile/exp/gl/glutil"
    "golang.org/x/mobile/exp/f32"
)

type Obj struct {
    vcount    int
    coord     gl.Buffer
    color     []float32
    tex       gl.Texture
    uvcoord     gl.Buffer
    normal      gl.Buffer
}

type Shape struct {
    Objs     []Obj
}

type Shader struct {
    models       map[string]*wavefront.Object
    program      gl.Program

    vertCoord       gl.Attrib
    normal          gl.Attrib
    texcoord        gl.Attrib

    projectionmatrix   gl.Uniform
    viewmatrix          gl.Uniform
    modelmatrix         gl.Uniform
    normalmatrix        gl.Uniform
    lightdir            gl.Uniform
    lightmatrix         gl.Uniform
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchx  float32
    touchy  float32
    started  time.Time
}

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func (e *Engine) Start() {
    var err error

    e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
    if err != nil {
        panic(fmt.Sprintln("LoadProgram failed:", err))
    }

    e.shader.models, err = wavefront.Read("girl.obj")
    check(err)


    e.shader.projectionmatrix =   gl.GetUniformLocation(e.shader.program, "u_projectionMatrix")
    e.shader.viewmatrix =   gl.GetUniformLocation(e.shader.program, "u_viewMatrix")
    e.shader.modelmatrix =       gl.GetUniformLocation(e.shader.program, "u_modelMatrix")
    e.shader.normalmatrix =       gl.GetUniformLocation(e.shader.program, "u_normalMatrix")
    e.shader.lightdir =     gl.GetUniformLocation(e.shader.program, "u_lightDirection")
    e.shader.lightmatrix =     gl.GetUniformLocation(e.shader.program, "u_lightmatrix")

    e.shader.vertCoord =    gl.GetAttribLocation(e.shader.program, "a_vertex")
    e.shader.normal =    gl.GetAttribLocation(e.shader.program, "a_normal")
    e.shader.texcoord = gl.GetAttribLocation(e.shader.program, "a_texCoord")


    for _, model := range e.shader.models {
        for _, group := range model.Groups {
            //颜色
            color := group.Material.Ambient
            //顶点
            data := f32.Bytes(binary.LittleEndian, group.Vertexes...)
            vertexCount := len(group.Vertexes)/3
            databuf := gl.CreateBuffer()
            gl.BindBuffer(gl.ARRAY_BUFFER, databuf)
            gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW)
            //UV坐标
            textcoords := f32.Bytes(binary.LittleEndian, group.Textcoords...)
            uvbuf := gl.CreateBuffer()
            gl.BindBuffer(gl.ARRAY_BUFFER, uvbuf)
            gl.BufferData(gl.ARRAY_BUFFER, textcoords, gl.STATIC_DRAW)
            //发现坐标
            normals := f32.Bytes(binary.LittleEndian, group.Normals...)
            normalbuf := gl.CreateBuffer()
            gl.BindBuffer(gl.ARRAY_BUFFER, normalbuf)
            gl.BufferData(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW)

            tex, _ := LoadTexture(group.Material.Texturefile)
            e.shape.Objs = append(e.shape.Objs, Obj{vcount:vertexCount, coord:databuf, color:color,  uvcoord:uvbuf, tex:tex,normal:normalbuf})

        }

    }
}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    for _, buf := range e.shape.Objs {
        gl.DeleteBuffer(buf.coord)
    }
}


func (e *Engine) Draw(c config.Event) {

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.ClearColor(0.5, 0.5, 0.5, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    gl.UseProgram(e.shader.program)

    gl.Uniform3fv(e.shader.lightdir,[]float32{0.5,0.5,0.5})

    m := mgl32.Perspective(0.785, float32(c.WidthPt/c.HeightPt), 0.1, 10.0)
    gl.UniformMatrix4fv(e.shader.projectionmatrix, m[:])

    eye := mgl32.Vec3{5,4, 5}
    center := mgl32.Vec3{0, 0, 0}
    up := mgl32.Vec3{0, 1, 0}

    m = mgl32.LookAtV(eye, center, up)
    gl.UniformMatrix4fv(e.shader.viewmatrix, m[:])

    m = mgl32.HomogRotate3D((e.touchx/float32(c.WidthPx)-0.5)*6.28, mgl32.Vec3{0, 1, 0})
    gl.UniformMatrix4fv(e.shader.modelmatrix, m[:])

    m = mgl32.HomogRotate3D((e.touchx/float32(c.WidthPx)-0.5)*6.28, mgl32.Vec3{0, -1, 0})
    gl.UniformMatrix4fv(e.shader.lightmatrix, m[:])



    coordsPerVertex := 3
    for _, obj := range e.shape.Objs {
        gl.BindBuffer(gl.ARRAY_BUFFER, obj.coord)
        gl.EnableVertexAttribArray(e.shader.vertCoord)
        gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 12, 0)

        texCoordsPerVertex := 2
        gl.BindBuffer(gl.ARRAY_BUFFER, obj.uvcoord)
        gl.EnableVertexAttribArray(e.shader.texcoord)
        gl.VertexAttribPointer(e.shader.texcoord, texCoordsPerVertex, gl.FLOAT, false, 8, 0)

        gl.BindBuffer(gl.ARRAY_BUFFER,obj.normal)
        gl.EnableVertexAttribArray(e.shader.normal)
        gl.VertexAttribPointer(e.shader.normal, 3, gl.FLOAT, false, 12, 0)

        gl.BindTexture(gl.TEXTURE_2D, obj.tex)

        gl.DrawArrays(gl.TRIANGLES, 0, obj.vcount)

        gl.DisableVertexAttribArray(e.shader.texcoord)
        gl.DisableVertexAttribArray(e.shader.normal)
        gl.DisableVertexAttribArray(e.shader.vertCoord)
    }

    debug.DrawFPS(c)
}


func main() {
    e := Engine{}

    app.Main(func(a app.App) {
        var c config.Event
        for eve := range a.Events() {
            switch eve := app.Filter(eve).(type) {
                case lifecycle.Event:
                switch eve.Crosses(lifecycle.StageVisible) {
                    case lifecycle.CrossOn:
                    e.Start()
                    case lifecycle.CrossOff:
                    e.Stop()
                }
                case config.Event:
                c = eve
                e.touchx = float32(c.WidthPt / 2)
                e.touchy = float32(c.HeightPt / 2)
                case paint.Event:
                e.Draw(c)
                a.EndPaint(eve)
                case touch.Event:
                e.touchx = eve.X
                e.touchy = eve.Y
            }
        }
    })
}



func loadAsset(name string) ([]byte, error) {
    f, err := asset.Open(name)
    if err != nil {
        return nil, err
    }
    return ioutil.ReadAll(f)
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(vertexAsset, fragmentAsset string) (p gl.Program, err error) {
    vertexSrc, err := loadAsset(vertexAsset)
    if err != nil {
        return
    }

    fragmentSrc, err := loadAsset(fragmentAsset)
    if err != nil {
        return
    }

    p, err = glutil.CreateProgram(string(vertexSrc), string(fragmentSrc))
    return
}


func LoadTexture(name string) (tex gl.Texture, err error) {
    imgFile, err := asset.Open(name)
    if err != nil {
        return
    }
    img, _, err := image.Decode(imgFile)
    if err != nil {
        return
    }

    rgba := image.NewRGBA(img.Bounds())
    image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

    tex = gl.CreateTexture()
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, tex)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
    gl.TexImage2D(
    gl.TEXTURE_2D,
    0,
    rgba.Rect.Size().X,
    rgba.Rect.Size().Y,
    gl.RGBA,
    gl.UNSIGNED_BYTE,
    rgba.Pix)

    return
}