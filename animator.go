package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	ball1 := createBall(600, 200, 2, 3, 1, 100)
	ball2 := createBall(100, 500, 3, -2, 1, 100)
	ball3 := createBall(300, 200, 4, 1, 1, 100)
	ball4 := createBall(450, 200, 1, 4, 1, 100)
	ball5 := createBall(750, 250, -3, 3, 1, 100)
	ball6 := createBall(200, 250, 2, 2, 1, 100)
	model := createModel(1024, 768, []Ball{*ball1, *ball2, *ball3, *ball4, *ball5, *ball6})

	width, height, err := getImageConfig("ball.png")
	if err != nil {
		panic(err)
	}

	fmt.Println(width, height)

	// Create window for drawing on
	// Kinda weird fitting the window onto the screen
	// Bounces seem to happen slightly off screen, but is just because of how window is drawn
	cfg := pixelgl.WindowConfig{
		Title:  "Bouncing Balls",
		Bounds: pixel.R(0, 0, model.width, model.height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	// load picture data for creating the sprite
	pic, err := loadPicture("ball.png")
	if err != nil {
		panic(err)
	}

	// create sprite from picture data
	sprite := pixel.NewSprite(pic, pic.Bounds())

	angle := 0.0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		angle += 3 * dt
		win.Clear(colornames.Whitesmoke)
		for _, ball := range model.balls {
			mat := pixel.IM
			//scaleX := (float32(ball.radius * 2.0)) / float32(width)
			//scaleY := (float32(ball.radius * 2.0)) / float32(height)

			mat = mat.Rotated(pixel.ZV, angle)
			position := pixel.V(ball.x, ball.y)
			//mat = mat.ScaledXY(position, pixel.V(float64(scaleX), float64(scaleY)))
			mat = mat.Moved(position)
			sprite.Draw(win, mat)
		}
		model.step(1)
		win.Update()
	}
}

// helper function for loading pictures
// Break this up into an additional function for getting width/height!
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getImageConfig(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}
