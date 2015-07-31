package main

import (
    "fmt"
    "encoding/binary"
    "io/ioutil"
    "log"
    "github.com/go-gl/mathgl/mgl32"
    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event/config"
    "golang.org/x/mobile/event/lifecycle"
    "golang.org/x/mobile/event/paint"
    "golang.org/x/mobile/event/touch"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/geom"
    "golang.org/x/mobile/gl"
    "golang.org/x/mobile/asset"
    "golang.org/x/mobile/exp/gl/glutil"
    "golang.org/x/mobile/exp/f32"
)

type Shape struct {
    buf     gl.Buffer
    texture gl.Texture
}

type Shader struct {
    program      gl.Program
    vertCoord    gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    modelx        gl.Uniform
    modely        gl.Uniform
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchLoc geom.Point
}

func (e *Engine) Start() {
    var err error
    e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
    if err != nil {
        panic(fmt.Sprintln("LoadProgram failed:", err))
    }

    fmt.Println(len(cubeData))
    e.shape.buf = gl.CreateBuffer()
    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)
    gl.BufferData(gl.ARRAY_BUFFER, cubeData, gl.STATIC_DRAW)
    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")
    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.modelx = gl.GetUniformLocation(e.shader.program, "modelx")
    e.shader.modely = gl.GetUniformLocation(e.shader.program, "modely")

}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    gl.DeleteBuffer(e.shape.buf)
}


func (e *Engine) Draw(c config.Event) {
    var vertexCount=len(cubeData)
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.ClearColor(0, 0, 0, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    gl.UseProgram(e.shader.program)

    m := mgl32.Perspective(0.8, float32(c.WidthPt/c.HeightPt), 0.1, 10.0)
    gl.UniformMatrix4fv(e.shader.projection, m[:])

    eye := mgl32.Vec3{0, 2, 6}
    center := mgl32.Vec3{0, 0, 0}
    up := mgl32.Vec3{0, 1, 0}

    m = mgl32.LookAtV(eye, center, up)
    gl.UniformMatrix4fv(e.shader.view, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.X*5/c.WidthPt), mgl32.Vec3{0, 1, 0})
    gl.UniformMatrix4fv(e.shader.modelx, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.Y*5/c.HeightPt), mgl32.Vec3{1, 0, 0})
    gl.UniformMatrix4fv(e.shader.modely, m[:])

    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)
    gl.EnableVertexAttribArray(e.shader.vertCoord)
    gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 0, 0)

    gl.DrawArrays(gl.LINES, 0, vertexCount)

    gl.DisableVertexAttribArray(e.shader.vertCoord)

    debug.DrawFPS(c)
}

//var vdata = LoadOBJ()

//var cubeData = f32.Bytes(binary.LittleEndian,vdata...)

var cubeData = f32.Bytes(binary.LittleEndian,   //三角
0.5,0.5,0.5,
0.5,0.5,-0.5,

0.5,0.5,-0.5,
0.5,-0.5,-0.5,

0.5,-0.5,-0.5,
0.5,-0.5,0.5,

0.5,-0.5,0.5,
0.5,0.5,0.5,


-0.5,0.5,0.5,
-0.5,0.5,-0.5,

-0.5,0.5,-0.5,
-0.5,-0.5,-0.5,

-0.5,-0.5,-0.5,
-0.5,-0.5,0.5,

-0.5,-0.5,0.5,
-0.5,0.5,0.5,

0.5,0.5,0.5,
-0.5,0.5,0.5,

0.5,0.5,-0.5,
-0.5,0.5,-0.5,

0.5,-0.5,-0.5,
-0.5,-0.5,-0.5,

0.5,-0.5,0.5,
-0.5,-0.5,0.5,

)


const (
    coordsPerVertex = 3 //坐标属性个数 x y z
    //vertexCount     = 36 //总点数
    colorsPerVertex = 4 //颜色属性个数 r g b a
)

func main() {
    defer func() {
        if err := recover(); err!=nil {
            log.Println(err)
        }
    }()
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
                e.touchLoc = geom.Point{0, 0}
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
    defer  f.Close()
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

