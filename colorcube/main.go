// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
    "encoding/binary"
    "log"

    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/exp/f32"
    "golang.org/x/mobile/exp/gl/glutil"
    "golang.org/x/mobile/geom"
    "golang.org/x/mobile/gl"
)

var (
    program  gl.Program
    position gl.Attrib
    offset   gl.Uniform
    color    gl.Uniform
    buf      gl.Buffer

    green    float32
    touchLoc geom.Point
)

func main() {
    app.Run(app.Callbacks{
        Start:  start,
        Stop:   stop,
        Draw:   draw,  //持续触发,每次触发生成一帧图像
        Touch:  touch, //触摸屏幕,以及滑动时触发
        Config: config,//初始化及窗口大小调整位置调整时触发
    })
}

func start() {
    var err error
    program, err = glutil.CreateProgram(vertexShader, fragmentShader)
    if err != nil {
        log.Printf("error creating GL program: %v", err)
        return
    }

    buf = gl.CreateBuffer()
    gl.BindBuffer(gl.ARRAY_BUFFER, buf)
    gl.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

    position = gl.GetAttribLocation(program, "position")//获取位置对象
    color =  gl.GetUniformLocation(program, "color") // 获取颜色对象
    offset = gl.GetUniformLocation(program, "offset") // 获取偏移对象
    // fmt.Println(position.String(),color.String(),offset.String())//Attrib(0) Uniform(1) Uniform(0)
    // TODO(crawshaw): the debug package needs to put GL state init here
    // Can this be an event.Register call now??
}

//停止时触发,清理
func stop() {
    gl.DeleteProgram(program)
    gl.DeleteBuffer(buf)
}

func config(new, old event.Config) {
    log.Println(new,old)
    touchLoc = geom.Point{new.Width / 2, new.Height / 2}
}

func touch(t event.Touch, c event.Config) {
    touchLoc = t.Loc
    log.Println(t.Loc)
}

func draw(c event.Config) {
    //清场
    gl.ClearColor(1, 0, 0, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)

    //使用program
    gl.UseProgram(program)

    green += 0.01
    if green > 1 {
        green = 0
    }
    gl.Uniform4f(color, 0, green, 0, 1)//设置color对象值,设置4个浮点数.
    //offset有两个值X,Y,窗口左上角为(0,0),右下角为(1,1)
    gl.Uniform2f(offset, float32(touchLoc.X/c.Width), float32(touchLoc.Y/c.Height))//设置偏移量对象,设置2个浮点数;

    log.Println(offset, float32(touchLoc.X/c.Width), float32(touchLoc.Y/c.Height),touchLoc,c)
    gl.BindBuffer(gl.ARRAY_BUFFER, buf)
    gl.EnableVertexAttribArray(position)
    gl.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)//顶点属性
    gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
    gl.DisableVertexAttribArray(position)

    debug.DrawFPS(c)
}

var triangleData = f32.Bytes(binary.LittleEndian,   //三角
0.0, 1.3, 0.0, // top left
0.0, -1.3, 0.0, // bottom left
0.4, 0.0, 0.0, // bottom right
)

const (
    coordsPerVertex = 3 //点坐标
    vertexCount     = 3 //点数
)

//两类着色器编程使用GLSL(GL Shader Language，GL着色语言)，它是OpenGL的一部分。与C或Java不同，GLSL必须在运行时编译，这意味着每次启动程序，所有的着色器将重新编译。
//顶点(vertex)着色器，它将作用于每个顶点上
//vec2即2个值,vec4即4个值
const vertexShader = `#version 330
#define _DEBUG_VERSION
layout (std140) uniform Matrices {
	mat4 pvm;
} ;
in vec4 position;
#if defined(_DEBUG_VERSION)
out vec4 color;
#endif
void main()
{
#if defined(_DEBUG_VERSION)
	 color = position;
#endif
	gl_Position = pvm * position ;
} `

//片断（Fragment）着色器，它将作用于每一个采样点
const fragmentShader = `#version 330

#define _DEBUG_VERSION

#if defined(_DEBUG_VERSION)
in  vec4 color;
#else
uniform vec4 color = vec4(1.0, 0.0, 0.0, 1.0);
#endif

out vec4 outputF;

void main()
{
	outputF = color;
} `
