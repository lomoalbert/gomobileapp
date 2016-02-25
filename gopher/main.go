package main

import (
"fmt"
"time"
    "github.com/lomoalbert/wavefront"
"github.com/go-gl/mathgl/mgl32"
"golang.org/x/mobile/app"
"golang.org/x/mobile/event/size"
"golang.org/x/mobile/event/lifecycle"
"golang.org/x/mobile/event/paint"
"golang.org/x/mobile/event/touch"
"golang.org/x/mobile/exp/app/debug"
"golang.org/x/mobile/geom"
"golang.org/x/mobile/gl"
"encoding/binary"
_ "image/png"
"io/ioutil"
"golang.org/x/mobile/asset"
"golang.org/x/mobile/exp/gl/glutil"
"golang.org/x/mobile/exp/f32"
)

type Buf struct{
    vcount    int
    coord     gl.Buffer
    color     []float32
}

type Shape struct {
    bufs     []Buf
}

type Shader struct {
    models       map[string]*wavefront.Object
    program      gl.Program
    vertCoord    gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    modelx        gl.Uniform
    modely        gl.Uniform
    color        gl.Uniform
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchLoc geom.Point
    started  time.Time
}

func check(err error){
    if err != nil{
        panic(err.Error())
    }
}

func (e *Engine) Start() {
    var err error

    e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
    if err != nil {
        panic(fmt.Sprintln("LoadProgram failed:", err))
    }

    e.shader.models,err = wavefront.Read("gopher.obj")
    check(err)

    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")
    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.modelx = gl.GetUniformLocation(e.shader.program, "modelx")
    e.shader.modely = gl.GetUniformLocation(e.shader.program, "modely")
    e.shader.color = gl.GetUniformLocation(e.shader.program, "color")

    for _,model := range e.shader.models{
        for _,group := range model.Groups{
            data:=f32.Bytes(binary.LittleEndian,group.Vertexes...)
            color:=group.Material.Ambient
            vertexCount := len(data)

            databuf := gl.CreateBuffer()
            e.shape.bufs = append(e.shape.bufs,Buf{vertexCount,databuf,color})

            gl.BindBuffer(gl.ARRAY_BUFFER, databuf)
            gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW)


        }

    }
}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    for _,buf := range e.shape.bufs{
        gl.DeleteBuffer(buf.coord)
    }
}


func (e *Engine) Draw(c size.Event) {

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.ClearColor(0.2, 0.2, 0.2, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    gl.UseProgram(e.shader.program)

    m := mgl32.Perspective(0.785, float32(c.WidthPt/c.HeightPt), 0.1, 10.0)
    gl.UniformMatrix4fv(e.shader.projection, m[:])

    eye := mgl32.Vec3{0, 0, 5}
    center := mgl32.Vec3{0, 0, 0}
    up := mgl32.Vec3{0, 1, 0}

    m = mgl32.LookAtV(eye, center, up)
    gl.UniformMatrix4fv(e.shader.view, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.X*10/c.WidthPt), mgl32.Vec3{0, 1, 0})
    gl.UniformMatrix4fv(e.shader.modelx, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.Y*10/c.HeightPt), mgl32.Vec3{1, 0, 0})
    gl.UniformMatrix4fv(e.shader.modely, m[:])

    coordsPerVertex :=3
    for _,buf := range e.shape.bufs{
        gl.BindBuffer(gl.ARRAY_BUFFER, buf.coord)
        gl.EnableVertexAttribArray(e.shader.vertCoord)
        gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 0, 0)
        gl.Uniform4f(e.shader.color,buf.color[0],buf.color[1],buf.color[2],buf.color[3])
        gl.DrawArrays(gl.TRIANGLES, 0, buf.vcount)

        gl.DisableVertexAttribArray(e.shader.vertCoord)
    }


    debug.DrawFPS(c)
}


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
                e.touchLoc = geom.Point{geom.Pt(eve.X),geom.Pt(eve.Y)}
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