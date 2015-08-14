package main

import (
    "fmt"
    "time"
    "github.com/go-gl/mathgl/mgl32"
    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event/size"
    "golang.org/x/mobile/event/lifecycle"
    "golang.org/x/mobile/event/paint"
    "golang.org/x/mobile/event/touch"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/geom"
    "golang.org/x/mobile/gl"
    "bytes"
    "encoding/binary"
    _ "image/png"
    "io/ioutil"
    "golang.org/x/mobile/asset"
    "golang.org/x/mobile/exp/gl/glutil"
    "golang.org/x/mobile/exp/f32"
)

type Shape struct {
    buf     gl.Buffer
    colorbuf     gl.Buffer
}

type Shader struct {
    program      gl.Program
    vertCoord    gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    modelx        gl.Uniform
    modely        gl.Uniform
    color        gl.Attrib
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchLoc geom.Point
    started  time.Time
}

func (e *Engine) Start() {
    var err error

    e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
    if err != nil {
        panic(fmt.Sprintln("LoadProgram failed:", err))
    }

    e.shape.buf = gl.CreateBuffer()
    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)
    gl.BufferData(gl.ARRAY_BUFFER, cubeData, gl.STATIC_DRAW)
    fmt.Println(len(cubeData))
    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")
    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.modelx = gl.GetUniformLocation(e.shader.program, "modelx")
    e.shader.modely = gl.GetUniformLocation(e.shader.program, "modely")

    e.shader.color = gl.GetAttribLocation(e.shader.program, "color")
    e.shape.colorbuf = gl.CreateBuffer()
    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.colorbuf)
    gl.BufferData(gl.ARRAY_BUFFER, colorData, gl.STATIC_DRAW)
    gl.VertexAttribPointer(e.shader.color, colorsPerVertex, gl.FLOAT, false, 4, 0) //更新color值
    gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)



    e.started = time.Now()
}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    gl.DeleteBuffer(e.shape.buf)
    gl.DeleteBuffer(e.shape.colorbuf)
}


func (e *Engine) Draw(c size.Event) {

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.ClearColor(0, 0, 0, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    gl.UseProgram(e.shader.program)

    m := mgl32.Perspective(0.785, float32(c.WidthPt/c.HeightPt), 0.1, 10.0)
    gl.UniformMatrix4fv(e.shader.projection, m[:])

    eye := mgl32.Vec3{3, 3, 3}
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

    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.colorbuf)
    gl.EnableVertexAttribArray(e.shader.color)
    gl.VertexAttribPointer(e.shader.color, colorsPerVertex, gl.FLOAT, false, 0, 0) //更新color值

    gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

    gl.DisableVertexAttribArray(e.shader.vertCoord)
    gl.DisableVertexAttribArray(e.shader.color)

    debug.DrawFPS(c)
}

var cubeData = f32.Bytes(binary.LittleEndian,
    -1.0, -1.0, -1.0,
    1.0, -1.0, -1.0,
    -1.0, -1.0, 1.0,
    1.0, -1.0, -1.0,
    1.0, -1.0, 1.0,
    -1.0, -1.0, 1.0,

    -1.0, 1.0, -1.0,
    -1.0, 1.0, 1.0,
    1.0, 1.0, -1.0,
    1.0, 1.0, -1.0,
    -1.0, 1.0, 1.0,
    1.0, 1.0, 1.0,

    -1.0, -1.0, 1.0,
    1.0, -1.0, 1.0,
    -1.0, 1.0, 1.0,
    1.0, -1.0, 1.0,
    1.0, 1.0, 1.0,
    -1.0, 1.0, 1.0,

    -1.0, -1.0, -1.0,
    -1.0, 1.0, -1.0,
    1.0, -1.0, -1.0,
    1.0, -1.0, -1.0,
    -1.0, 1.0, -1.0,
    1.0, 1.0, -1.0,

    -1.0, -1.0, 1.0,
    -1.0, 1.0, -1.0,
    -1.0, -1.0, -1.0,
    -1.0, -1.0, 1.0,
    -1.0, 1.0, 1.0,
    -1.0, 1.0, -1.0,

    1.0, -1.0, 1.0,
    1.0, -1.0, -1.0,
    1.0, 1.0, -1.0,
    1.0, -1.0, 1.0,
    1.0, 1.0, -1.0,
    1.0, 1.0, 1.0,
)

var colorData = f32.Bytes(binary.LittleEndian,   //过渡色
1,1,1,1,
1,1,0,1,
1,0,1,1,
1,1,0,1,
0,1,0,1,
1,0,1,1,
0,0,0,1,
1,0,0,1,
0,0,1,1,
0,0,1,1,
1,0,0,1,
0,1,1,1,
1,0,1,1,
0,1,0,1,
1,0,0,1,
0,1,0,1,
0,1,1,1,
1,0,0,1,
1,1,1,1,
0,0,0,1,
1,1,0,1,
1,1,0,1,
0,0,0,1,
0,0,1,1,
1,0,1,1,
0,0,0,1,
1,1,1,1,
1,0,1,1,
1,0,0,1,
0,0,0,1,
0,1,0,1,
1,1,0,1,
0,0,1,1,
0,1,0,1,
0,0,1,1,
0,1,1,1,

)


const (
    coordsPerVertex = 3 //坐标属性个数 x y z
    vertexCount     = 36 //总点数
    colorsPerVertex = 4 //颜色属性个数 r g b a
)

func main() {
    e := Engine{}
    app.Main(func(a app.App) {
        var c size.Event
        for eve := range a.Events() {
            switch eve := app.Filter(eve).(type) {
                case lifecycle.Event:
                switch eve.Crosses(lifecycle.StageVisible) {
                    case lifecycle.CrossOn:
                    e.Start()
                    case lifecycle.CrossOff:
                    e.Stop()
                }
                case size.Event:
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


// EncodeObject converts float32 vertices into a LittleEndian byte array.
func EncodeObject(vertices ...[]float32) []byte {
    buf := bytes.Buffer{}
    for _, v := range vertices {
        err := binary.Write(&buf, binary.LittleEndian, v)
        if err != nil {
            panic(fmt.Sprintln("binary.Write failed:", err))
        }
    }

    return buf.Bytes()
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