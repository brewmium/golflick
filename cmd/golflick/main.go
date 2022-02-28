package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1080
	screenHeight = 640
)

//go:embed ball1.png
var ball1_png []byte

//go:embed ball2.png
var ball2_png []byte

//go:embed grass.jpg
var grass_jpg []byte

var (
	ball1Image *ebiten.Image
	ball2Image *ebiten.Image
	grassImage *ebiten.Image
)

type ball struct {
	body   *box2d.B2Body
	radius float64
	image  *ebiten.Image
	scale  float64
}

type Game struct {
	world box2d.B2World
	balls []ball

	y  int
	dy int
}

func (g *Game) Update() error {
	g.world.Step(1.0/60.0, 8, 3)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(grassImage, op)

	for _, b := range g.balls {
		pos := b.body.GetPosition()
		op.GeoM.Reset()
		op.GeoM.Translate(-b.radius, -b.radius)
		op.GeoM.Rotate(-b.body.GetAngle())
		op.GeoM.Scale(b.scale, b.scale)
		op.GeoM.Translate(pos.X+540.0, 320.0-pos.Y)
		screen.DrawImage(b.image, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func initBox(g *Game) {
	gravity := box2d.MakeB2Vec2(0.0, -100.0)

	g.world = box2d.MakeB2World(gravity)

	// Ground body
	{
		bd := box2d.MakeB2BodyDef()
		ground := g.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(-540.0, -320.0), box2d.MakeB2Vec2(540.0, -320.0))
		ground.CreateFixture(&shape, 0.0)
	}

	// Left wall body
	{
		bd := box2d.MakeB2BodyDef()
		wall := g.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(-540.0, -320.0), box2d.MakeB2Vec2(-540.0, 320.0))
		wall.CreateFixture(&shape, 0.0)
	}

	// Right wall body
	{
		bd := box2d.MakeB2BodyDef()
		wall := g.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(540.0, -320.0), box2d.MakeB2Vec2(540.0, 320.0))
		wall.CreateFixture(&shape, 0.0)
	}

	// Top wall body
	{
		bd := box2d.MakeB2BodyDef()
		wall := g.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(-540.0, 320.0), box2d.MakeB2Vec2(540.0, 320.0))
		wall.CreateFixture(&shape, 0.0)
	}

	addBall(g, -300.0, -100.0, 75.0, ball1Image, 1.0)
	addBall(g, -250.0, 200.0, 50.0, ball2Image, 1.0)
	addBall(g, -320.0, 250.0, 75.0, ball1Image, 0.4)
	addBall(g, -200.0, 100.0, 50.0, ball2Image, 1.5)
}

func addBall(g *Game, x, y, r float64, im *ebiten.Image, scale float64) {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(x, y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	bd.AngularDamping = 0.75

	body := g.world.CreateBody(&bd)

	shape := box2d.MakeB2CircleShape()
	shape.M_radius = r * scale

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Friction = 1.0
	fd.Density = 20.0
	fd.Restitution = 0.6
	body.CreateFixtureFromDef(&fd)

	b := ball{
		body:   body,
		radius: r,
		image:  im,
		scale:  scale,
	}
	g.balls = append(g.balls, b)
}

func initImage(im []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(im))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func main() {
	ball1Image = initImage(ball1_png)
	ball2Image = initImage(ball2_png)
	grassImage = initImage(grass_jpg)

	g := &Game{}
	initBox(g)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Waaa!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
