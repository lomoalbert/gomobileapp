package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	image_draw "image/draw"
	_ "image/png"
	"io/ioutil"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

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
