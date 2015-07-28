package main

import (
    "fmt"
    "time"
    "github.com/go-gl/mathgl/mgl32"
    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event/config"
    "golang.org/x/mobile/event/lifecycle"
    "golang.org/x/mobile/event/paint"
    "golang.org/x/mobile/event/touch"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/geom"
    "golang.org/x/mobile/gl"
    "bytes"
    "encoding/binary"
    "image"
    image_draw "image/draw"
    _ "image/png"
    "io/ioutil"
    "golang.org/x/mobile/asset"
    "golang.org/x/mobile/exp/gl/glutil"
)

type Shape struct {
    buf     gl.Buffer
    texture gl.Texture
}

type Shader struct {
    program      gl.Program
    vertCoord    gl.Attrib
    vertTexCoord gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    model        gl.Uniform
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
    gl.BufferData(gl.ARRAY_BUFFER, EncodeObject(cubeData), gl.STATIC_DRAW)

    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")
    e.shader.vertTexCoord = gl.GetAttribLocation(e.shader.program, "vertTexCoord")

    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.model = gl.GetUniformLocation(e.shader.program, "model")

    e.shape.texture, err = LoadTexture("gopher.png")
    if err != nil {
        panic(fmt.Sprintln("LoadTexture failed:", err))
    }

    e.started = time.Now()
}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    gl.DeleteBuffer(e.shape.buf)
}


func (e *Engine) Draw(c config.Event) {
    since := time.Now().Sub(e.started)
    //gl.Enable()

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

    m = mgl32.HomogRotate3D(float32(since.Seconds()), mgl32.Vec3{0, 1, 0})
    gl.UniformMatrix4fv(e.shader.model, m[:])

    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)

    coordsPerVertex := 3
    texCoordsPerVertex := 2
    vertexCount := len(cubeData) / (coordsPerVertex + texCoordsPerVertex)

    gl.EnableVertexAttribArray(e.shader.vertCoord)
    gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 20, 0) // 4 bytes in float, 5 values per vertex

    gl.EnableVertexAttribArray(e.shader.vertTexCoord)
    gl.VertexAttribPointer(e.shader.vertTexCoord, texCoordsPerVertex, gl.FLOAT, false, 20, 12)

    gl.BindTexture(gl.TEXTURE_2D, e.shape.texture)

    gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

    gl.DisableVertexAttribArray(e.shader.vertCoord)

    debug.DrawFPS(c)
}

var cubeData = []float32{
    //  X, Y, Z, U, V
    // Bottom
    -1.0, -1.0, -1.0, 0.0, 0.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, -1.0, 1.0, 1.0, 1.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,

    // Top
    -1.0, 1.0, -1.0, 0.0, 0.0,
    -1.0, 1.0, 1.0, 0.0, 1.0,
    1.0, 1.0, -1.0, 1.0, 0.0,
    1.0, 1.0, -1.0, 1.0, 0.0,
    -1.0, 1.0, 1.0, 0.0, 1.0,
    1.0, 1.0, 1.0, 1.0, 1.0,

    // Front
    -1.0, -1.0, 1.0, 1.0, 0.0,
    1.0, -1.0, 1.0, 0.0, 0.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,
    1.0, -1.0, 1.0, 0.0, 0.0,
    1.0, 1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,

    // Back
    -1.0, -1.0, -1.0, 0.0, 0.0,
    -1.0, 1.0, -1.0, 0.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    -1.0, 1.0, -1.0, 0.0, 1.0,
    1.0, 1.0, -1.0, 1.0, 1.0,

    // Left
    -1.0, -1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, -1.0, 1.0, 0.0,
    -1.0, -1.0, -1.0, 0.0, 0.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,
    -1.0, 1.0, -1.0, 1.0, 0.0,

    // Right
    1.0, -1.0, 1.0, 1.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, 1.0, -1.0, 0.0, 0.0,
    1.0, -1.0, 1.0, 1.0, 1.0,
    1.0, 1.0, -1.0, 0.0, 0.0,
    1.0, 1.0, 1.0, 0.0, 1.0,
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
                a.EndPaint()
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

// LoadTexture reads and decodes an image from the asset repository and creates
// a texture object based on the full dimensions of the image.
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
