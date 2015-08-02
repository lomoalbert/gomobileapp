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
    "golang.org/x/mobile/geom"
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
    useuv       bool
    tex       gl.Texture
    uvcoord     gl.Buffer
}

type Shape struct {
    Objs     []Obj
}

type Shader struct {
    models       map[string]*wavefront.Object
    program      gl.Program
    vertCoord    gl.Attrib
    vertTexCoord gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    modelx        gl.Uniform
    modely        gl.Uniform
    color        gl.Uniform
    useuv       gl.Uniform
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchLoc geom.Point
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

    e.shader.models, err = wavefront.Read("spiritframe.obj")
    check(err)

    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")
    e.shader.vertTexCoord = gl.GetAttribLocation(e.shader.program, "vertTexCoord")
    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.modelx = gl.GetUniformLocation(e.shader.program, "modelx")
    e.shader.modely = gl.GetUniformLocation(e.shader.program, "modely")
    e.shader.color = gl.GetUniformLocation(e.shader.program, "color")
    e.shader.useuv= gl.GetUniformLocation(e.shader.program, "useuv")

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
            //贴图文件
            var useuv bool
            tex, err := LoadTexture(group.Material.Texturefile)
            if err !=nil {
                useuv = false
            }else {
                useuv = true
            }
            e.shape.Objs = append(e.shape.Objs, Obj{vcount:vertexCount, coord:databuf, color:color, useuv:useuv, uvcoord:uvbuf,tex:tex})

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

    m := mgl32.Perspective(0.785, float32(c.WidthPt/c.HeightPt), 0.1, 10.0)
    gl.UniformMatrix4fv(e.shader.projection, m[:])

    eye := mgl32.Vec3{0,0 , 0.01}
    center := mgl32.Vec3{0, 0, 0}
    up := mgl32.Vec3{0, 1, 0}

    m = mgl32.LookAtV(eye, center, up)
    gl.UniformMatrix4fv(e.shader.view, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.X/c.WidthPt-0.5)*6.28, mgl32.Vec3{0, 1, 0})
    gl.UniformMatrix4fv(e.shader.modelx, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.Y/c.HeightPt-0.5)*3.14, mgl32.Vec3{1, 0, 0})
    gl.UniformMatrix4fv(e.shader.modely, m[:])

    coordsPerVertex := 3
    for _, obj := range e.shape.Objs {
        gl.BindBuffer(gl.ARRAY_BUFFER, obj.coord)
        gl.EnableVertexAttribArray(e.shader.vertCoord)

        gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 12, 0)

        if obj.useuv==true{
            gl.Uniform1i(e.shader.useuv,1)
            texCoordsPerVertex := 2
            gl.BindBuffer(gl.ARRAY_BUFFER, obj.uvcoord)
            gl.EnableVertexAttribArray(e.shader.vertTexCoord)
            gl.VertexAttribPointer(e.shader.vertTexCoord, texCoordsPerVertex, gl.FLOAT, false, 8, 0)

            gl.BindTexture(gl.TEXTURE_2D, obj.tex)
        }else{
            gl.Uniform1i(e.shader.useuv,0)
            gl.Uniform4f(e.shader.color, obj.color[0], obj.color[1], obj.color[2], obj.color[3])
        }
        gl.DrawArrays(gl.TRIANGLES, 0, obj.vcount)
        if obj.useuv{
            gl.DisableVertexAttribArray(e.shader.vertTexCoord)
        }
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
                e.touchLoc = geom.Point{c.WidthPt / 2, c.HeightPt / 2}
                case paint.Event:
                e.Draw(c)
                a.EndPaint(eve)
                case touch.Event:
                e.touchLoc = eve.Loc
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