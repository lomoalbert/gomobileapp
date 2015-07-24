package main

import (
    "fmt"
    "strings"
    "strconv"
)


var v [][3]float32
var f []float32

func checkerr(err error){
    if err!=nil{
        fmt.Println(err.Error())
        panic(err)
    }
}

func getv(line string) [3]float32 {
    elems := strings.Split(line," ")
    v1,err:=strconv.ParseFloat(elems[1],32)
    checkerr(err)
    v2,err:=strconv.ParseFloat(elems[2],32)
    checkerr(err)
    v3,err:=strconv.ParseFloat(elems[3],32)
    checkerr(err)
    return [3]float32{float32(v1),float32(v2),float32(v3),}
}

func getf(line string) []float32 {
    elems := strings.Split(line," ")
    fv := []float32{}
    for _,elem := range elems[1:]{
        vs := strings.Split(elem,"/")[0]
        vi,err := strconv.Atoi(vs)
        checkerr(err)
        for _,ft := range v[vi-1]{
            fv=append(fv,ft)
        }
    }
    return fv
}


func LoadOBJ(filepath string) []float32{
    data,err:= loadAsset(filepath)
    checkerr(err)
    lines:=strings.Split(string(data),"\n")
    for _,line := range lines{
        switch strings.Split(line," ")[0]{
            case "v":
            v=append(v,getv(line))
            case "f":
            f=append(f,getf(line)...)
        }
    }
    return f
}