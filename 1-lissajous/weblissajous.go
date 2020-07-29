package main

import (
  "image"
  "image/color"
  "image/gif"
  "io"
  "math"
  "math/rand"
  "net/http"
)

var palette = []color.Color{color.Black, color.RGBA{0x00, 0xff, 0x00, 0xff}}

const (
  whiteIndex = iota
  blackIndex
)

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    lissajous(w)
  })
  http.ListenAndServe("localhost:8000", nil)
}

func lissajous(out io.Writer) {
  const (
    cycles            = 5
    angularResolution = 0.00001
    canvasSize        = 500 //[-size..+size] -- so a size of 100 results in 201 pixels across. 100 on each side, plus a 0th pixel.
    nFrames           = 64
    delay             = 1
  )

  yOscillatorFrequency := rand.Float64() * 3.0
  animation := gif.GIF{LoopCount: nFrames}
  phaseDifference := 0.0

  for i := 0; i < nFrames; i++ {
    rect := image.Rect(0, 0, 2*canvasSize+1, 2*canvasSize+1)
    img := image.NewPaletted(rect, palette)
    for t := 0.0; t < cycles*2*math.Pi; t += angularResolution {
      x := math.Sin(t)
      y := math.Sin(t*yOscillatorFrequency + phaseDifference)
      img.SetColorIndex(canvasSize+int(x*canvasSize+0.5), canvasSize+int(y*canvasSize+0.5),
        blackIndex)
    }
    phaseDifference += 0.1
    animation.Delay = append(animation.Delay, delay)
    animation.Image = append(animation.Image, img)
  }
  gif.EncodeAll(out, &animation)
}
