package main

import (
	"math/rand"
	"strconv"
	"time"

	rl "github.com/lachee/raylib-goplus/raylib"
)

var ( // MARK: var

	collectedpowerupcount                                           int
	collectedpowerups                                               = make([]blokmap, 10)
	maxpowerups                                                     = 3
	powerupcount                                                    = 1
	powerupmap                                                      = make([]blokmap, 100)
	coinstotal                                                      int
	bulletcount                                                     int
	bulletmap                                                       = make([]blokmap, 500)
	playervelocity                                                  = 1
	explodefade                                                     = float32(1.0)
	explodeon                                                       bool
	explodetimer                                                    = 90
	exploderecs                                                     = make([]rl.Rectangle, 10)
	backmap                                                         = make([]blokback, 100)
	backtype, backtimer, backdirection                              int
	backon                                                          bool
	player                                                          blokplayer
	layermap                                                        = make([]blokmap, grida)
	gridw                                                           = 1000
	gridh                                                           = 1000
	grida                                                           = gridw * gridh
	draww, drawh                                                    = 240, 240
	drawa                                                           = draww * drawh
	drawnext                                                        int
	monw                                                            = 1280
	monh                                                            = 720
	fps                                                             = 30
	framecount                                                      int
	nextblockon, pauseon, numberson, debugon, gridon, centerlineson bool
	onoff2, onoff3, onoff6, onoff10, onoff15, onoff30               bool
	imgs                                                            rl.Texture2D
	camera, camera2x, camera4x, camera8x                            rl.Camera2D
	// long variables
	bulletmodifier1, bulletmodifier2, bulletmodifier3, bulletmodifier4, bulletmodifier5, bulletmodifier6, bulletmodifier7, bulletmodifier8, bulletmodifier9, bulletmodifier10                     int
	poweruptypecount1, poweruptypecount2, poweruptypecount3, poweruptypecount4, poweruptypecount5, poweruptypecount6, poweruptypecount7, poweruptypecount8, poweruptypecount9, poweruptypecount10 int
	// tiles
	tilepowerup = rl.NewRectangle(0, 0, 32, 32)
	tilecrate   = rl.NewRectangle(32, 0, 32, 32)
	// power ups
	coin                = rl.NewRectangle(0, 144, 16, 16)
	coin2               = rl.NewRectangle(0, 109, 32, 32)
	imghart             = rl.NewRectangle(601, 6, 36, 36)
	imghorizbullets     = rl.NewRectangle(562, 9, 36, 36)
	imgasteroid         = rl.NewRectangle(520, 7, 36, 36)
	imgcirclebullets    = rl.NewRectangle(472, 11, 36, 36)
	imgcrossbullets     = rl.NewRectangle(432, 11, 36, 36)
	imgfanbullets       = rl.NewRectangle(386, 11, 36, 36)
	imgelectricbullets  = rl.NewRectangle(349, 11, 36, 36)
	img8bullets         = rl.NewRectangle(308, 11, 36, 36)
	imgexplodingbullets = rl.NewRectangle(264, 10, 36, 36)
	imgrandombullets    = rl.NewRectangle(222, 9, 36, 36)
)

// MARK: structs
type blokback struct {
	color        rl.Color
	fade         float32
	x, y, length int
	fadeon       bool
}
type blokmap struct {
	name                                                             string
	color, color2                                                    rl.Color
	powerup, imgon, coin, activ, special, leftright, bullet, visible bool
	img                                                              rl.Rectangle
	rotation                                                         float32
	velocity, direction, position, x, y, objecttype, numberinmap     int
}
type blokplayer struct {
	color               rl.Color
	activ               bool
	direction, position int
	rotation            float32
}

func raylib() { // MARK: raylib
	rl.InitWindow(monw, monh, "space")
	rl.SetExitKey(rl.KeyEnd)          // key to end the game and close window
	imgs = rl.LoadTexture("imgs.png") // load images
	rl.SetTargetFPS(fps)
	//rl.HideCursor()
	//rl.ToggleFullscreen()
	for !rl.WindowShouldClose() {
		framecount++
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawnocameraback()
		rl.BeginMode2D(camera)
		drawlayers()
		rl.EndMode2D()
		drawnocamera()
		drawmenubars()
		update()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func drawlayers() { // MARK: drawlayers

	// explosions
	if explodeon {
		explodefade -= 0.05

		for a := 0; a < len(exploderecs); a++ {
			rl.DrawRectangleRec(exploderecs[a], rl.Fade(brightbrown(), explodefade))
			switch a {
			case 0:
				exploderecs[a].X -= rFloat32(10, 21)
			case 1:
				exploderecs[a].X -= rFloat32(5, 10)
				exploderecs[a].Y -= rFloat32(10, 21)
			case 2:
				exploderecs[a].X -= rFloat32(1, 4)
				exploderecs[a].Y -= rFloat32(10, 21)
			case 3:
				exploderecs[a].Y -= rFloat32(10, 21)
			case 4:
				exploderecs[a].X += rFloat32(1, 4)
				exploderecs[a].Y -= rFloat32(10, 21)
			case 5:
				exploderecs[a].X += rFloat32(5, 10)
				exploderecs[a].Y -= rFloat32(10, 21)
			case 6:
				exploderecs[a].X += rFloat32(10, 21)
			case 7:
				exploderecs[a].X += rFloat32(5, 10)
				exploderecs[a].Y += rFloat32(10, 21)
			case 8:
				exploderecs[a].X += rFloat32(1, 4)
				exploderecs[a].Y += rFloat32(10, 21)
			case 9:
				exploderecs[a].X += rFloat32(10, 21)
				exploderecs[a].Y += rFloat32(10, 21)
			}
		}
	}
	// layer 1
	x := 0
	y := 0
	count := 0
	drawblock := drawnext

	for a := 0; a < drawa; a++ {
		layermap[drawblock].x = x
		layermap[drawblock].y = y
		// draw stationary blocks
		if layermap[drawblock].activ && layermap[drawblock].visible && !layermap[drawblock].special {
			rl.DrawRectangle(x, y, 16, 16, rl.Black)
			rl.DrawRectangle(x, y, 15, 15, layermap[drawblock].color)
		} else if layermap[drawblock].activ && layermap[drawblock].special { // draw rotating blocks
			destrec := rl.NewRectangle(float32(x+16), float32(y+16), 32, 32)
			v2 := rl.NewVector2(16, 16)
			rl.DrawTexturePro(imgs, layermap[drawblock].img, destrec, v2, layermap[drawblock].rotation, layermap[drawblock].color)
			layermap[drawblock].rotation += 2
		}

		if layermap[drawblock].bullet { // draw bullets
			rl.DrawCircle(x+8, y+8, 6, layermap[drawblock].color)
		}
		if layermap[drawblock].powerup { // draw powerups
			//	rl.DrawRectangle(x, y, 15, 15, layermap[drawblock].color)
			if drawblock == player.position {
				collectpowerup(drawblock)
			}
			if layermap[drawblock].imgon {
				v2 := rl.NewVector2(float32(x), float32(y))
				rl.DrawTextureRec(imgs, layermap[drawblock].img, v2, rl.White)
			}
		}

		if layermap[drawblock].coin { // draw coins
			if drawblock == player.position {
				collectcoin(drawblock)
			}
			if layermap[drawblock].imgon {
				v2 := rl.NewVector2(float32(x), float32(y))
				rl.DrawTextureRec(imgs, coin2, v2, rl.White)
			}
		}
		x += 16
		drawblock++
		count++

		if count == draww {
			count = 0
			drawblock += gridw
			drawblock -= draww
			x = 0
			y += 16
		}
	}

	// layer 2
	x = 0
	y = 0
	count = 0
	drawblock = drawnext
	for a := 0; a < drawa; a++ {

		if drawblock == player.position {
			v2 := rl.NewVector2(float32(x+8), float32(y+8))
			rl.DrawPoly(v2, 3, 16, player.rotation, player.color)

		}

		x += 16
		drawblock++
		count++

		if count == draww {
			count = 0
			drawblock += gridw
			drawblock -= draww
			x = 0
			y += 16
		}
	}

}
func drawmenubars() { // MARK: drawmenubars ██████████████████████████████████████████

	rl.DrawRectangle(0, 0, 60, monh, rl.Fade(rl.Black, 0.7))
	rl.DrawRectangle(monw-60, 0, 60, monh, rl.Fade(rl.Black, 0.7))
	// right menu
	v2 := rl.NewVector2(float32(monw/2)-22, 10)
	rl.BeginMode2D(camera2x)
	rl.DrawTextureRec(imgs, coin, v2, rl.White)
	rl.EndMode2D()
	cointext := strconv.Itoa(coinstotal)
	rl.DrawText(cointext, monw-34, 60, 20, rl.White)
	// left menu

	y := 20
	for a := 0; a < maxpowerups; a++ {
		rl.DrawRectangle(12, y, 36, 36, rl.Fade(rl.Black, 0.8))
		y += 44
	}
	y = 20
	for a := 0; a < len(collectedpowerups); a++ {
		v2 := rl.NewVector2(12, float32(y))
		switch collectedpowerups[a].objecttype {
		case 1:
			rl.DrawTextureRec(imgs, imgasteroid, v2, rl.White)
		case 2:
			rl.DrawTextureRec(imgs, imgcirclebullets, v2, rl.White)
		case 3:
			rl.DrawTextureRec(imgs, imgcrossbullets, v2, rl.White)
		case 4:
			rl.DrawTextureRec(imgs, imgelectricbullets, v2, rl.White)
		case 5:
			rl.DrawTextureRec(imgs, imgfanbullets, v2, rl.White)
		case 6:
			rl.DrawTextureRec(imgs, imghart, v2, rl.White)
		case 7:
			rl.DrawTextureRec(imgs, img8bullets, v2, rl.White)
		case 8:
			rl.DrawTextureRec(imgs, imghorizbullets, v2, rl.White)
		case 9:
			rl.DrawTextureRec(imgs, imgrandombullets, v2, rl.White)
		case 10:
			rl.DrawTextureRec(imgs, imgexplodingbullets, v2, rl.White)
		}
		y += 44
	}
	updatepoweruptypes()
}
func updatepoweruptypes() { // MARK: updatepoweruptypes

	poweruptypecount1 = 0
	poweruptypecount2 = 0
	poweruptypecount3 = 0
	poweruptypecount4 = 0
	poweruptypecount5 = 0
	poweruptypecount6 = 0
	poweruptypecount7 = 0
	poweruptypecount8 = 0
	poweruptypecount9 = 0
	poweruptypecount10 = 0

	for a := 0; a < len(collectedpowerups); a++ {
		switch collectedpowerups[a].objecttype {
		case 1:
			poweruptypecount1++
		case 2:
			poweruptypecount2++
		case 3:
			poweruptypecount3++
		case 4:
			poweruptypecount4++
		case 5:
			poweruptypecount5++
		case 6:
			poweruptypecount6++
		case 7:
			poweruptypecount7++
		case 8:
			poweruptypecount8++
		case 9:
			poweruptypecount9++
		case 10:
			poweruptypecount10++
		}
	}

	bulletmodifier1 = poweruptypecount1
	bulletmodifier2 = poweruptypecount2
	bulletmodifier3 = poweruptypecount3
	bulletmodifier4 = poweruptypecount4
	bulletmodifier5 = poweruptypecount5
	bulletmodifier6 = poweruptypecount6
	bulletmodifier7 = poweruptypecount7
	bulletmodifier8 = poweruptypecount8
	bulletmodifier9 = poweruptypecount9
	bulletmodifier10 = poweruptypecount10

}
func drawnocameraback() { // MARK: drawnocameraback

	if backon {
		switch backtype {
		case 1: // colorful stars
			for a := 0; a < len(backmap); a++ {
				rl.DrawRectangle(backmap[a].x, backmap[a].y, backmap[a].length, backmap[a].length, rl.Fade(backmap[a].color, backmap[a].fade))
				if backmap[a].fadeon {
					backmap[a].fade += 0.02
					if backmap[a].fade >= 1.0 {
						backmap[a].fadeon = false
					}
				} else {
					backmap[a].fade -= 0.02
					if backmap[a].fade <= 0.0 {
						backmap[a].fadeon = true
					}
				}
				if flipcoin() {
					if flipcoin() {
						backmap[a].x += rInt(-2, 3)
					} else {
						backmap[a].y += rInt(-2, 3)
					}
				}
			}
		case 2: // horizontal scrolling colorful stars
			for a := 0; a < len(backmap); a++ {
				rl.DrawRectangle(backmap[a].x, backmap[a].y, backmap[a].length, backmap[a].length, rl.Fade(backmap[a].color, backmap[a].fade))
				backmap[a].x += rInt(5, 11)
				if backmap[a].x > monw {
					backmap[a].x = 10
				}
				if backmap[a].fadeon {
					backmap[a].fade += 0.02
					if backmap[a].fade >= 1.0 {
						backmap[a].fadeon = false
					}
				} else {
					backmap[a].fade -= 0.02
					if backmap[a].fade <= 0.0 {
						backmap[a].fadeon = true
					}
				}
			}
		case 3: // horizontal scrolling colorful stars
			for a := 0; a < len(backmap); a++ {
				rl.DrawRectangle(backmap[a].x, backmap[a].y, backmap[a].length, backmap[a].length, rl.Fade(backmap[a].color, backmap[a].fade))

				switch backdirection {
				case 1:
					backmap[a].x -= rInt(5, 11)
					backmap[a].y += rInt(5, 11)
				case 2:
					backmap[a].y += rInt(5, 11)
				case 3:
					backmap[a].x += rInt(5, 11)
					backmap[a].y += rInt(5, 11)
				case 4:
					backmap[a].x -= rInt(5, 11)
				case 6:
					backmap[a].x += rInt(5, 11)
				case 7:
					backmap[a].x -= rInt(5, 11)
					backmap[a].y -= rInt(5, 11)
				case 8:
					backmap[a].y -= rInt(5, 11)
				case 9:
					backmap[a].x += rInt(5, 11)
					backmap[a].y -= rInt(5, 11)
				}

				if backmap[a].x > monw {
					backmap[a].x = 10
				} else if backmap[a].x < 10 {
					backmap[a].x = monw - 10
				}
				if backmap[a].y > monh {
					backmap[a].y = 10
				} else if backmap[a].y < 10 {
					backmap[a].y = monh - 10
				}

				if backmap[a].fadeon {
					backmap[a].fade += 0.02
					if backmap[a].fade >= 1.0 {
						backmap[a].fadeon = false
					}
				} else {
					backmap[a].fade -= 0.02
					if backmap[a].fade <= 0.0 {
						backmap[a].fadeon = true
					}
				}
			}
			backtimer++
			if backtimer == 120 {
				backtimer = 0
				for {
					backdirection = rInt(1, 10)
					if backdirection != 5 {
						break
					}
				}
			}
		}

	}

}
func drawnocamera() { // MARK: drawnocamera

	// centerlines
	if centerlineson {
		rl.DrawLine(monw/2, 0, monw/2, monh, rl.Green)
		rl.DrawLine(0, monh/2, monw, monh/2, rl.Green)
	}

}

func update() { // MARK: update

	if !pauseon {

		moveplayer()
		updatebullets()
		if onoff3 {
			movepowerups()
		}
		timers()

		if camera.Zoom == 1.0 {
			drawnext = player.position
			drawnext -= 40
			drawnext -= 22 * gridw
		} else if camera.Zoom == 2.0 {
			drawnext = player.position
			drawnext -= 20
			drawnext -= 11 * gridw

		}

		updateplayer()
	}
	input()

	if debugon {
		debug()
	}
}
func updatebullets() { // MARK: updatebullets

	for a := 0; a < len(bulletmap); a++ {
		if bulletmap[a].bullet {
			switch bulletmap[a].direction {
			case 1:
				if !layermap[(bulletmap[a].position+gridw)-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					layermap[(bulletmap[a].position+gridw)-1].bullet = false
					explodebullet(layermap[(bulletmap[a].position+gridw)-1].x, layermap[(bulletmap[a].position+gridw)-1].y)
					if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {
						clearblocks((bulletmap[a].position+gridw)-1, "crate")
						explode((bulletmap[a].position+gridw)-1, "crate")
					}
					if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position + gridw) + 1)
					}
				}
				if !layermap[(bulletmap[a].position+gridw)-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					layermap[(bulletmap[a].position+gridw)-1].bullet = false
					explodebullet(layermap[(bulletmap[a].position+gridw)-1].x, layermap[(bulletmap[a].position+gridw)-1].y)
					if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {
						clearblocks((bulletmap[a].position+gridw)-1, "crate")
						explode((bulletmap[a].position+gridw)-1, "crate")

					}
					if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position + gridw) + 1)
					}
				}
			case 2:
				if !layermap[bulletmap[a].position+gridw].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[bulletmap[a].position+gridw].x, layermap[bulletmap[a].position+gridw].y)
					if layermap[bulletmap[a].position+gridw].name == "crate" {
						clearblocks(bulletmap[a].position+gridw, "crate")
						explode(bulletmap[a].position+gridw, "crate")
					}
					if layermap[bulletmap[a].position+gridw].name == "powerup" {
						newpowerup(bulletmap[a].position + gridw)
					}
				}
				if !layermap[bulletmap[a].position+gridw].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[bulletmap[a].position+gridw].x, layermap[bulletmap[a].position-gridw].y)
					if layermap[bulletmap[a].position+gridw].name == "crate" {
						clearblocks(bulletmap[a].position+gridw, "crate")
						explode(bulletmap[a].position+gridw, "crate")
					}
					if layermap[bulletmap[a].position+gridw].name == "powerup" {
						newpowerup(bulletmap[a].position + gridw)
					}
				}
			case 3:
				if !layermap[(bulletmap[a].position+gridw)+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					layermap[(bulletmap[a].position+gridw)+1].bullet = false
					explodebullet(layermap[(bulletmap[a].position+gridw)+1].x, layermap[(bulletmap[a].position+gridw)+1].y)
					if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {
						explode((bulletmap[a].position+gridw)+1, "crate")
						clearblocks((bulletmap[a].position+gridw)+1, "crate")

					}
					if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position + gridw) + 1)
					}
				}
				if !layermap[(bulletmap[a].position+gridw)+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position += gridw
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[(bulletmap[a].position+gridw)+1].x, layermap[bulletmap[a].position-gridw].y)
					if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {
						explode((bulletmap[a].position+gridw)+1, "crate")
						clearblocks((bulletmap[a].position+gridw)+1, "crate")

					}
					if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position + gridw) + 1)
					}
				}
			case 4:
				if !layermap[bulletmap[a].position-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[bulletmap[a].position-1].x, layermap[bulletmap[a].position-1].y)
					if layermap[bulletmap[a].position-1].name == "crate" {
						explode(bulletmap[a].position-1, "crate")
						clearblocks(bulletmap[a].position-1, "crate")

					}
					if layermap[bulletmap[a].position-1].name == "powerup" {
						newpowerup(bulletmap[a].position - 1)
					}
				}
				if !layermap[bulletmap[a].position-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[bulletmap[a].position-1].x, layermap[bulletmap[a].position-1].y)
					if layermap[bulletmap[a].position-1].name == "crate" {
						explode(bulletmap[a].position-1, "crate")
						clearblocks(bulletmap[a].position-1, "crate")

					}
					if layermap[bulletmap[a].position-1].name == "powerup" {
						newpowerup(bulletmap[a].position - 1)
					}
				}
			case 6:
				if !layermap[bulletmap[a].position+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[bulletmap[a].position+1].x, layermap[bulletmap[a].position+1].y)
					if layermap[bulletmap[a].position+1].name == "crate" {
						explode(bulletmap[a].position+1, "crate")
						clearblocks(bulletmap[a].position+1, "crate")

					}
					if layermap[bulletmap[a].position+1].name == "powerup" {
						newpowerup(bulletmap[a].position + 1)
					}
				}
				if !layermap[bulletmap[a].position+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[bulletmap[a].position+1].x, layermap[bulletmap[a].position+1].y)
					if layermap[bulletmap[a].position+1].name == "crate" {
						explode(bulletmap[a].position+1, "crate")
						clearblocks(bulletmap[a].position+1, "crate")

					}
					if layermap[bulletmap[a].position+1].name == "powerup" {
						newpowerup(bulletmap[a].position + 1)
					}
				}
			case 7:
				if !layermap[(bulletmap[a].position-gridw)-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[(bulletmap[a].position-gridw)-1].x, layermap[(bulletmap[a].position-gridw)-1].y)
					if layermap[(bulletmap[a].position-gridw)-1].name == "crate" {
						explode((bulletmap[a].position-gridw)-1, "crate")
						clearblocks((bulletmap[a].position-gridw)-1, "crate")

					}
					if layermap[(bulletmap[a].position-gridw)-1].name == "powerup" {
						newpowerup((bulletmap[a].position - gridw) - 1)
					}
				}
				if !layermap[(bulletmap[a].position-gridw)-1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
					bulletmap[a].position--
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[(bulletmap[a].position-gridw)-1].x, layermap[(bulletmap[a].position-gridw)-1].y)
					if layermap[(bulletmap[a].position-gridw)-1].name == "crate" {
						explode((bulletmap[a].position-gridw)-1, "crate")
						clearblocks((bulletmap[a].position-gridw)-1, "crate")

					}
					if layermap[(bulletmap[a].position-gridw)-1].name == "powerup" {
						newpowerup((bulletmap[a].position - gridw) - 1)
					}
				}
			case 8:
				if !layermap[bulletmap[a].position-gridw].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[bulletmap[a].position-gridw].x, layermap[bulletmap[a].position-gridw].y)
					if layermap[bulletmap[a].position-gridw].name == "crate" {
						explode(bulletmap[a].position-gridw, "crate")
						clearblocks(bulletmap[a].position-gridw, "crate")

					}
					if layermap[bulletmap[a].position-gridw].name == "powerup" {
						newpowerup(bulletmap[a].position - gridw)
					}
				}
				if !layermap[bulletmap[a].position-gridw].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[bulletmap[a].position-gridw].x, layermap[bulletmap[a].position-gridw].y)
					if layermap[bulletmap[a].position-gridw].name == "crate" {
						explode(bulletmap[a].position-gridw, "crate")
						clearblocks(bulletmap[a].position-gridw, "crate")

					}
					if layermap[bulletmap[a].position-gridw].name == "powerup" {
						newpowerup(bulletmap[a].position - gridw)
					}
				}
			case 9:
				if !layermap[(bulletmap[a].position-gridw)+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					layermap[bulletmap[a].position].bullet = false
					explodebullet(layermap[(bulletmap[a].position-gridw)+1].x, layermap[(bulletmap[a].position-gridw)+1].y)
					if layermap[(bulletmap[a].position-gridw)+1].name == "crate" {
						explode((bulletmap[a].position-gridw)+1, "crate")
						clearblocks((bulletmap[a].position-gridw)+1, "crate")

					}
					if layermap[(bulletmap[a].position-gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position - gridw) + 1)
					}
				}
				if !layermap[(bulletmap[a].position-gridw)+1].activ {
					layermap[bulletmap[a].position].bullet = false
					bulletmap[a].position -= gridw
					bulletmap[a].position++
				} else {
					bulletmap[a].bullet = false
					explodebullet(layermap[(bulletmap[a].position-gridw)+1].x, layermap[(bulletmap[a].position-gridw)+1].y)
					if layermap[(bulletmap[a].position-gridw)+1].name == "crate" {
						explode((bulletmap[a].position-gridw)+1, "crate")
						clearblocks((bulletmap[a].position-gridw)+1, "crate")

					}
					if layermap[(bulletmap[a].position-gridw)+1].name == "powerup" {
						newpowerup((bulletmap[a].position - gridw) + 1)
					}
				}
			}
		}
	}
	if playervelocity <= 3 && playervelocity > 1 {
		for a := 0; a < len(bulletmap); a++ {
			if bulletmap[a].bullet {
				switch bulletmap[a].direction {
				case 1:
					if !layermap[(bulletmap[a].position+gridw)-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[(bulletmap[a].position+gridw)-1].bullet = false
						explodebullet(layermap[(bulletmap[a].position+gridw)-1].x, layermap[(bulletmap[a].position+gridw)-1].y)
						if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {
							clearblocks((bulletmap[a].position+gridw)-1, "crate")
						}
						if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position + gridw) + 1)
						}
					}
				case 2:
					if !layermap[bulletmap[a].position+gridw].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position+gridw].x, layermap[bulletmap[a].position+gridw].y)
						if layermap[bulletmap[a].position+gridw].name == "crate" {

							clearblocks(bulletmap[a].position+gridw, "crate")

						}
						if layermap[bulletmap[a].position+gridw].name == "powerup" {
							newpowerup(bulletmap[a].position + gridw)
						}

					}
				case 3:
					if !layermap[(bulletmap[a].position+gridw)+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[(bulletmap[a].position+gridw)+1].bullet = false
						explodebullet(layermap[(bulletmap[a].position+gridw)+1].x, layermap[(bulletmap[a].position+gridw)+1].y)
						if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {

							clearblocks((bulletmap[a].position+gridw)+1, "crate")

						}
						if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position + gridw) + 1)
						}
					}
				case 4:
					if !layermap[bulletmap[a].position-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position-1].x, layermap[bulletmap[a].position-1].y)
						if layermap[bulletmap[a].position-1].name == "crate" {

							clearblocks(bulletmap[a].position-1, "crate")

						}
						if layermap[bulletmap[a].position-1].name == "powerup" {
							newpowerup(bulletmap[a].position - 1)
						}
					}
				case 6:
					if !layermap[bulletmap[a].position+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position+1].x, layermap[bulletmap[a].position+1].y)
						if layermap[bulletmap[a].position+1].name == "crate" {

							clearblocks(bulletmap[a].position+1, "crate")

						}
						if layermap[bulletmap[a].position+1].name == "powerup" {
							newpowerup(bulletmap[a].position + 1)
						}
					}
				case 7:
					if !layermap[(bulletmap[a].position-gridw)-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[(bulletmap[a].position-gridw)-1].x, layermap[(bulletmap[a].position-gridw)-1].y)
						if layermap[(bulletmap[a].position-gridw)-1].name == "crate" {

							clearblocks((bulletmap[a].position-gridw)-1, "crate")

						}
						if layermap[(bulletmap[a].position-gridw)-1].name == "powerup" {
							newpowerup((bulletmap[a].position - gridw) - 1)
						}
					}
				case 8:
					if !layermap[bulletmap[a].position-gridw].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position-gridw].x, layermap[bulletmap[a].position-gridw].y)
						if layermap[bulletmap[a].position-gridw].name == "crate" {

							clearblocks(bulletmap[a].position-gridw, "crate")

						}
						if layermap[bulletmap[a].position-gridw].name == "powerup" {
							newpowerup(bulletmap[a].position - gridw)
						}
					}
				case 9:
					if !layermap[(bulletmap[a].position-gridw)+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[(bulletmap[a].position-gridw)+1].x, layermap[(bulletmap[a].position-gridw)+1].y)
						if layermap[(bulletmap[a].position-gridw)+1].name == "crate" {

							clearblocks((bulletmap[a].position-gridw)+1, "crate")

						}
						if layermap[(bulletmap[a].position-gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position - gridw) + 1)
						}
					}
				}
			}
		}
	}
	if playervelocity == 3 {
		for a := 0; a < len(bulletmap); a++ {
			if bulletmap[a].bullet {
				switch bulletmap[a].direction {
				case 1:
					if !layermap[(bulletmap[a].position+gridw)-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[(bulletmap[a].position+gridw)-1].bullet = false
						explodebullet(layermap[(bulletmap[a].position+gridw)-1].x, layermap[(bulletmap[a].position+gridw)-1].y)
						if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {

							clearblocks((bulletmap[a].position+gridw)-1, "crate")

						}
						if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position + gridw) + 1)
						}
					}
				case 2:
					if !layermap[bulletmap[a].position+gridw].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position+gridw].x, layermap[bulletmap[a].position+gridw].y)
						if layermap[bulletmap[a].position+gridw].name == "crate" {

							clearblocks(bulletmap[a].position+gridw, "crate")

						}
						if layermap[bulletmap[a].position+gridw].name == "powerup" {
							newpowerup(bulletmap[a].position + gridw)
						}

					}
				case 3:
					if !layermap[(bulletmap[a].position+gridw)+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position += gridw
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[(bulletmap[a].position+gridw)+1].bullet = false
						explodebullet(layermap[(bulletmap[a].position+gridw)+1].x, layermap[(bulletmap[a].position+gridw)+1].y)
						if layermap[(bulletmap[a].position+gridw)+1].name == "crate" {

							clearblocks((bulletmap[a].position+gridw)+1, "crate")

						}
						if layermap[(bulletmap[a].position+gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position + gridw) + 1)
						}
					}
				case 4:
					if !layermap[bulletmap[a].position-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position-1].x, layermap[bulletmap[a].position-1].y)
						if layermap[bulletmap[a].position-1].name == "crate" {

							clearblocks(bulletmap[a].position-1, "crate")

						}
						if layermap[bulletmap[a].position-1].name == "powerup" {
							newpowerup(bulletmap[a].position - 1)
						}
					}
				case 6:
					if !layermap[bulletmap[a].position+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position+1].x, layermap[bulletmap[a].position+1].y)
						if layermap[bulletmap[a].position+1].name == "crate" {

							clearblocks(bulletmap[a].position+1, "crate")

						}
						if layermap[bulletmap[a].position+1].name == "powerup" {
							newpowerup(bulletmap[a].position + 1)
						}
					}
				case 7:
					if !layermap[(bulletmap[a].position-gridw)-1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
						bulletmap[a].position--
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[(bulletmap[a].position-gridw)-1].x, layermap[(bulletmap[a].position-gridw)-1].y)
						if layermap[(bulletmap[a].position-gridw)-1].name == "crate" {

							clearblocks((bulletmap[a].position-gridw)-1, "crate")

						}
						if layermap[(bulletmap[a].position-gridw)-1].name == "powerup" {
							newpowerup((bulletmap[a].position - gridw) - 1)
						}
					}
				case 8:
					if !layermap[bulletmap[a].position-gridw].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[bulletmap[a].position-gridw].x, layermap[bulletmap[a].position-gridw].y)
						if layermap[bulletmap[a].position-gridw].name == "crate" {

							clearblocks(bulletmap[a].position-gridw, "crate")

						}
						if layermap[bulletmap[a].position-gridw].name == "powerup" {
							newpowerup(bulletmap[a].position - gridw)
						}
					}
				case 9:
					if !layermap[(bulletmap[a].position-gridw)+1].activ {
						layermap[bulletmap[a].position].bullet = false
						bulletmap[a].position -= gridw
						bulletmap[a].position++
					} else {
						bulletmap[a].bullet = false
						layermap[bulletmap[a].position].bullet = false
						explodebullet(layermap[(bulletmap[a].position-gridw)+1].x, layermap[(bulletmap[a].position-gridw)+1].y)
						if layermap[(bulletmap[a].position-gridw)+1].name == "crate" {

							clearblocks((bulletmap[a].position-gridw)+1, "crate")

						}
						if layermap[(bulletmap[a].position-gridw)+1].name == "powerup" {
							newpowerup((bulletmap[a].position - gridw) + 1)
						}
					}
				}

			}
		}
	}

	for a := 0; a < len(bulletmap); a++ {
		if bulletmap[a].bullet {
			layermap[bulletmap[a].position] = bulletmap[a]
		}
	}

}
func updateplayer() { // MARK: updateplayer

	switch player.direction {
	case 4:
		player.rotation = 90.0
	case 1:
		player.rotation = 45.0
	case 2:
		player.rotation = 0.0
	case 3:
		player.rotation = 315.0
	case 6:
		player.rotation = 270.0
	case 7:
		player.rotation = 135.0
	case 8:
		player.rotation = 180.0
	case 9:
		player.rotation = 225.0
	}

}
func movepowerups() { // MARK: movepowerups

	for a := 0; a < len(powerupmap); a++ {

		if powerupmap[a].powerup {
			switch powerupmap[a].direction {
			case 1:
				if !layermap[powerupmap[a].position+((gridw*2)-1)].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position += gridw - 1
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 2:
				if !layermap[powerupmap[a].position+gridw*2].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position += gridw
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 3:
				if !layermap[powerupmap[a].position+((gridw*2)+1)].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position += gridw + 1
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 4:
				if !layermap[powerupmap[a].position-1].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position--
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 6:
				if !layermap[powerupmap[a].position+2].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position++
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 7:
				if !layermap[powerupmap[a].position-(gridw-1)].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position -= (gridw - 1)
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 8:
				if !layermap[powerupmap[a].position-gridw].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position -= gridw
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			case 9:
				if !layermap[powerupmap[a].position-(gridw+1)].activ {
					layermap[powerupmap[a].position].powerup = false
					layermap[powerupmap[a].position+1].powerup = false
					layermap[(powerupmap[a].position+1)+gridw].powerup = false
					layermap[powerupmap[a].position+gridw].powerup = false
					powerupmap[a].position -= (gridw + 1)
					if rolldice() == 6 {
						powerupmap[a] = changedirection(powerupmap[a])
					}
				} else {
					powerupmap[a] = changedirection(powerupmap[a])
				}
			}

		}

	}

	for a := 0; a < len(powerupmap); a++ {

		if powerupmap[a].powerup {
			layermap[powerupmap[a].position] = powerupmap[a]
			layermap[powerupmap[a].position+1].powerup = true
			layermap[powerupmap[a].position+1].objecttype = powerupmap[a].objecttype
			layermap[(powerupmap[a].position+1)+gridw].powerup = true
			layermap[(powerupmap[a].position+1)+gridw].objecttype = powerupmap[a].objecttype
			layermap[powerupmap[a].position+gridw].powerup = true
			layermap[powerupmap[a].position+gridw].objecttype = powerupmap[a].objecttype
		}

	}

}
func changedirection(block blokmap) blokmap { // MARK: changedirection

	for {
		block.direction = rInt(1, 10)
		if block.direction != 5 {
			break
		}
	}

	return block

}
func clearblocks(block int, name string) { // MARK: clearblocks

	block--
	block -= gridw
	count := 0
	switch name {
	case "crate":
		for a := 0; a < 9; a++ {
			if layermap[block].name == "crate" {
				layermap[block].activ = false
			}
			block++
			count++
			if count == 3 {
				count = 0
				block += gridw
				block -= 3
			}
		}
	case "coin":
		for a := 0; a < 9; a++ {
			if layermap[block].coin {
				layermap[block].coin = false
				layermap[block].imgon = false
			}
			block++
			count++
			if count == 3 {
				count = 0
				block += gridw
				block -= 3
			}
		}
	case "powerup":
		for a := 0; a < 9; a++ {
			if layermap[block].powerup {
				layermap[block].powerup = false
				layermap[block].objecttype = 0
				layermap[block].imgon = false
				if layermap[block].numberinmap != 0 {
					powerupmap[layermap[block].numberinmap] = blokmap{}
					layermap[block].numberinmap = 0
				}
			}
			block++
			count++
			if count == 3 {
				count = 0
				block += gridw
				block -= 3
			}
		}
	}

}
func explode(block int, name string) { // MARK: explode

	x := layermap[block].x + 8
	y := layermap[block].y + 8

	explodetimer = 90
	explodeon = true
	explodefade = 1.0

	switch name {
	case "crate":
		for a := 0; a < len(exploderecs); a++ {
			length := rInt(8, 15)
			exploderecs[a] = rl.NewRectangle(float32(x+rInt(-4, 5)), float32(y+rInt(-4, 5)), float32(length), float32(length))
		}
		layermap[block].imgon = true
		layermap[block].img = coin2
		layermap[block].coin = true
		layermap[block+1].coin = true
		layermap[block+gridw].coin = true
		layermap[(block+gridw)+1].coin = true
	}

}

func createmap() { // MARK: createmap

	startblock := gridw * 100
	startblock += 100

	player.position = startblock + 10
	player.position += gridw * 10
	drawnext = startblock

	length := 100
	area := length * length
	count := 0
	for a := 0; a < area; a++ {

		layermap[startblock].activ = false
		startblock++
		count++

		if count == length {
			count = 0
			startblock -= length
			startblock += gridw
		}
	}

}
func degrade() { // MARK: degrade

	for a := gridw * 100; a < grida-(gridw*100); a++ {
		if !layermap[a].activ && layermap[a-gridw].activ {
			number := rInt(0, 3)
			for b := 0; b < number; b++ {
				layermap[a-(gridw*b)].activ = false
			}
		}

	}
	for a := grida - (gridw * 100); a > gridw*100; a-- {
		if !layermap[a].activ && layermap[a+1].activ {
			number := rInt(0, 4)
			for b := 0; b < number; b++ {
				layermap[a+b].activ = false
			}
		}

		if layermap[a].activ && !layermap[a-gridw].activ {
			number := rInt(0, 3)
			for b := 0; b < number; b++ {
				layermap[a+(gridw*b)].activ = false
			}
		}
		if !layermap[a].activ && layermap[a-1].activ {
			number := rInt(0, 3)
			for b := 0; b < number; b++ {
				layermap[a-b].activ = false
			}
		}
	}

}
func collectpowerup(block int) { // MARK: collectpowerup ██████████████████████████████████████████

	collectedpowerups[collectedpowerupcount].objecttype = layermap[block].objecttype
	collectedpowerupcount++
	if collectedpowerupcount == maxpowerups {
		collectedpowerupcount = 0
	}
	clearblocks(block, "powerup")

}
func collectcoin(block int) { // MARK: collectcoin
	coinstotal++
	clearblocks(block, "coin")
}
func newpowerup(block int) { // MARK: newpowerup

	x := layermap[block].x
	y := layermap[block].y
	radius := float32(50)

	newpowerupblock := block

	block -= 2
	block -= gridw * 2
	count := 0
	for a := 0; a < 25; a++ {
		if layermap[block].name == "powerup" {
			layermap[block].color = brightgrey()
		}
		block++
		count++
		if count == 5 {
			block += gridw
			block -= 5
			count = 0
		}
	}
	for {
		if flipcoin() {
			newpowerupblock += rInt(3, 5)
		} else {
			newpowerupblock -= rInt(3, 5)
		}
		if flipcoin() {
			newpowerupblock += rInt(3, 5) * gridw
		} else {
			newpowerupblock -= rInt(3, 5) * gridw
		}

		if !layermap[newpowerupblock].activ && !layermap[newpowerupblock+1].activ && !layermap[newpowerupblock+gridw].activ && !layermap[(newpowerupblock+gridw)+1].activ {
			newpowerup := blokmap{}
			newpowerup.powerup = true
			newpowerup.visible = true
			newpowerup.position = block
			newpowerup.objecttype = rInt(1, 11)
			newpowerup.imgon = true
			switch newpowerup.objecttype {
			case 1:
				newpowerup.img = imgasteroid
			case 2:
				newpowerup.img = imgcirclebullets
			case 3:
				newpowerup.img = imgcrossbullets
			case 4:
				newpowerup.img = imgelectricbullets
			case 5:
				newpowerup.img = imgfanbullets
			case 6:
				newpowerup.img = imghart
			case 7:
				newpowerup.img = img8bullets
			case 8:
				newpowerup.img = imghorizbullets
			case 9:
				newpowerup.img = imgrandombullets
			case 10:
				newpowerup.img = imgexplodingbullets
			}
			newpowerup.color = randomorange()
			newpowerup.velocity = 1
			for {
				newpowerup.direction = rInt(1, 10)
				if newpowerup.direction != 5 {
					break
				}
			}
			newpowerup.numberinmap = powerupcount
			powerupmap[powerupcount] = newpowerup
			powerupcount++
			if powerupcount == 100 {
				powerupcount = 1
			}
			break
		}

	}

	for radius < 100 {
		rl.DrawCircleLines(x, y, radius, rl.Fade(randomcolor(), rF32(0.4, 1.0)))
		radius += 5
	}

}
func createpowerups() { // MARK: createpowerups

	number := 10

	for {

		choose := rInt(gridw*100, grida-(gridw*100))

		if !layermap[choose].activ {
			if flipcoin() {
				powerupblock := blokmap{}
				powerupblock.activ = true
				powerupblock.special = true
				powerupblock.img = tilepowerup
				powerupblock.color = brightyellow()
				powerupblock.rotation = 10.0
				powerupblock.name = "powerup"
				layermap[choose] = powerupblock
				choose++
				powerupblock = blokmap{}
				powerupblock.activ = true
				powerupblock.name = "powerup"
				layermap[choose] = powerupblock
				choose--
				choose += gridw
				layermap[choose] = powerupblock
				choose++
				layermap[choose] = powerupblock
			} else {
				powerupblock := blokmap{}
				powerupblock.activ = true
				powerupblock.special = true
				powerupblock.img = tilecrate
				powerupblock.color = brightbrown()
				powerupblock.rotation = 10.0
				powerupblock.name = "crate"
				layermap[choose] = powerupblock
				choose++
				powerupblock = blokmap{}
				powerupblock.activ = true
				powerupblock.name = "crate"
				layermap[choose] = powerupblock
				choose--
				choose += gridw
				layermap[choose] = powerupblock
				choose++
				layermap[choose] = powerupblock
			}
			number--
		}

		if number == 0 {
			break
		}

	}

}

func setinitialvalues() { // MARK: setinitialvalues

	backon = true
	backdirection = 1
	backtype = rInt(1, 4)
	for a := 0; a < len(backmap); a++ {
		backmap[a].color = randomcolor()
		backmap[a].x = rInt(10, monw-10)
		backmap[a].y = rInt(10, monh-10)
		backmap[a].fade = rF32(0.0, 1.0)
		backmap[a].length = rInt(1, 4)
		backmap[a].fadeon = flipcoin()
	}

	for a := 0; a < len(layermap); a++ {
		layermap[a].activ = true
		layermap[a].visible = true
		layermap[a].color = randombluedark()
		layermap[a].color2 = randombluedark()
	}

	createmap()
	degrade()
	createpowerups()
	player.color = rl.Blue

}
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLogLevel(rl.LogError) // hides info window
	rl.InitWindow(monw, monh, "setscreen")
	setscreen()
	rl.CloseWindow()
	setinitialvalues()
	raylib()
}
func timers() { // MARK: timers

	if onoff3 {
		coin.X += 16
		if coin.X > 65 {
			coin.X = 0
		}
		coin2.X += 32
		if coin2.X > 160 {
			coin2.X = 0
		}
	}

	if explodeon {
		explodetimer--
		if explodetimer == 0 {
			explodeon = false
		}
	}

	if framecount%2 == 0 {
		if onoff2 {
			onoff2 = false
		} else {
			onoff2 = true
		}
	}
	if framecount%3 == 0 {
		if onoff3 {
			onoff3 = false
		} else {
			onoff3 = true
		}
	}
	if framecount%6 == 0 {
		if onoff6 {
			onoff6 = false
		} else {
			onoff6 = true
		}
	}
	if framecount%10 == 0 {
		if onoff10 {
			onoff10 = false
		} else {
			onoff10 = true
		}
	}
	if framecount%15 == 0 {
		if onoff15 {
			onoff15 = false
		} else {
			onoff15 = true
		}
	}
	if framecount%30 == 0 {
		if onoff30 {
			onoff30 = false
		} else {
			onoff30 = true
		}
	}

}
func setscreen() { // MARK: setscreen

	rl.SetWindowSize(monw, monh)

	camera.Zoom = 1.0
	camera.Target.X = 0
	camera.Target.Y = 0

	camera2x.Zoom = 2.0
	camera4x.Zoom = 4.0
	camera8x.Zoom = 8.0

	camera4x.Target.X = 0
	camera4x.Target.Y = 0
}
func explodebullet(x, y int) { // MARK: explodebullet

	rl.DrawCircle(x, y, rFloat32(40, 61), rl.Fade(randomorange(), rF32(0.4, 0.9)))
	number := rInt(4, 8)

	for a := 0; a < number; a++ {
		choose := rolldice()
		if choose == 1 {
			rl.DrawCircle(x+(rInt(60, 120)), y+(rInt(60, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		} else if choose == 2 {
			rl.DrawCircle(x+(rInt(60, 120)), y-(rInt(60, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		} else if choose == 3 {
			rl.DrawCircle(x-(rInt(60, 120)), y+(rInt(60, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		} else if choose == 4 {
			rl.DrawCircle(x-(rInt(60, 120)), y-(rInt(60, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		} else if choose == 5 {
			rl.DrawCircle(x+(rInt(-120, 120)), y+(rInt(-120, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		} else if choose == 6 {
			rl.DrawCircle(x+(rInt(-120, 120)), y+(rInt(-120, 120)), rFloat32(5, 21), rl.Fade(randomorange(), rF32(0.4, 0.9)))
		}

	}

}
func shoot() { // MARK: shoot

	bulletmodifier7 = 0

	if bulletmodifier7 != 0 {
		if !layermap[(player.position+gridw)-1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 1
			bulletmap[bulletcount].position = (player.position + gridw) - 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position + gridw)].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 2
			bulletmap[bulletcount].position = (player.position + gridw)
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position+gridw)+1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 3
			bulletmap[bulletcount].position = (player.position + gridw) + 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position)-1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 4
			bulletmap[bulletcount].position = (player.position) - 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position)+1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 6
			bulletmap[bulletcount].position = (player.position) + 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position-gridw)-1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 7
			bulletmap[bulletcount].position = (player.position - gridw) - 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position - gridw)].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 8
			bulletmap[bulletcount].position = (player.position - gridw)
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}
		if !layermap[(player.position-gridw)+1].activ {
			bulletmap[bulletcount].bullet = true
			bulletmap[bulletcount].velocity = playervelocity + 1
			bulletmap[bulletcount].color = randomorange()
			bulletmap[bulletcount].direction = 9
			bulletmap[bulletcount].position = (player.position - gridw) + 1
			bulletcount++
			if bulletcount == 500 {
				bulletcount = 0
			}
		}

	} else {
		switch player.direction {
		case 1:
			if !layermap[(player.position+gridw)-1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 1
				bulletmap[bulletcount].position = (player.position + gridw) - 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 2:
			if !layermap[(player.position + gridw)].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 2
				bulletmap[bulletcount].position = (player.position + gridw)
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 3:
			if !layermap[(player.position+gridw)+1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 3
				bulletmap[bulletcount].position = (player.position + gridw) + 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 4:
			if !layermap[(player.position)-1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 4
				bulletmap[bulletcount].position = (player.position) - 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 6:
			if !layermap[(player.position)+1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 6
				bulletmap[bulletcount].position = (player.position) + 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 7:
			if !layermap[(player.position-gridw)-1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 7
				bulletmap[bulletcount].position = (player.position - gridw) - 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 8:
			if !layermap[(player.position - gridw)].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 8
				bulletmap[bulletcount].position = (player.position - gridw)
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}
		case 9:
			if !layermap[(player.position-gridw)+1].activ {
				bulletmap[bulletcount].bullet = true
				bulletmap[bulletcount].velocity = playervelocity + 1
				bulletmap[bulletcount].color = randomorange()
				bulletmap[bulletcount].direction = 9
				bulletmap[bulletcount].position = (player.position - gridw) + 1
				bulletcount++
				if bulletcount == 500 {
					bulletcount = 0
				}
			}

		}
	}
}

func moveplayer() {

	switch player.direction {
	case 1:
		if !layermap[(player.position-1)+gridw].activ {
			player.position--
			player.position += gridw
		} else {
			wallcollide(1)
		}
	case 2:
		if !layermap[player.position+gridw].activ {
			player.position += gridw
		} else {
			wallcollide(2)
		}
	case 3:
		if !layermap[(player.position+1)+gridw].activ {
			player.position++
			player.position += gridw
		} else {
			wallcollide(3)
		}
	case 4:
		if !layermap[player.position-1].activ {
			player.position--
		} else {
			wallcollide(4)
		}
	case 6:
		if !layermap[player.position+1].activ {
			player.position++
		} else {
			wallcollide(6)
		}
	case 7:
		if !layermap[(player.position-1)-gridw].activ {
			player.position--
			player.position -= gridw
		} else {
			wallcollide(7)
		}
	case 8:
		if !layermap[player.position-gridw].activ {
			player.position -= gridw
		} else {
			wallcollide(8)
		}
	case 9:
		if !layermap[(player.position+1)-gridw].activ {
			player.position++
			player.position -= gridw
		} else {
			wallcollide(9)
		}
	}
	if playervelocity > 1 {
		switch player.direction {
		case 1:
			if !layermap[(player.position-1)+gridw].activ {
				player.position--
				player.position += gridw
			} else {
				wallcollide(1)
			}
		case 2:
			if !layermap[player.position+gridw].activ {
				player.position += gridw
			} else {
				wallcollide(2)
			}
		case 3:
			if !layermap[(player.position+1)+gridw].activ {
				player.position++
				player.position += gridw
			} else {
				wallcollide(3)
			}
		case 4:
			if !layermap[player.position-1].activ {
				player.position--
			} else {
				wallcollide(4)
			}
		case 6:
			if !layermap[player.position+1].activ {
				player.position++
			} else {
				wallcollide(6)
			}
		case 7:
			if !layermap[(player.position-1)-gridw].activ {
				player.position--
				player.position -= gridw
			} else {
				wallcollide(7)
			}
		case 8:
			if !layermap[player.position-gridw].activ {
				player.position -= gridw
			} else {
				wallcollide(8)
			}
		case 9:
			if !layermap[(player.position+1)-gridw].activ {
				player.position++
				player.position -= gridw
			} else {
				wallcollide(9)
			}
		}
	}
	if playervelocity == 3 {
		switch player.direction {
		case 1:
			if !layermap[(player.position-1)+gridw].activ {
				player.position--
				player.position += gridw
			} else {
				wallcollide(1)
			}
		case 2:
			if !layermap[player.position+gridw].activ {
				player.position += gridw
			} else {
				wallcollide(2)
			}
		case 3:
			if !layermap[(player.position+1)+gridw].activ {
				player.position++
				player.position += gridw
			} else {
				wallcollide(3)
			}
		case 4:
			if !layermap[player.position-1].activ {
				player.position--
			} else {
				wallcollide(4)
			}
		case 6:
			if !layermap[player.position+1].activ {
				player.position++
			} else {
				wallcollide(6)
			}
		case 7:
			if !layermap[(player.position-1)-gridw].activ {
				player.position--
				player.position -= gridw
			} else {
				wallcollide(7)
			}
		case 8:
			if !layermap[player.position-gridw].activ {
				player.position -= gridw
			} else {
				wallcollide(8)
			}
		case 9:
			if !layermap[(player.position+1)-gridw].activ {
				player.position++
				player.position -= gridw
			} else {
				wallcollide(9)
			}
		}
	}
}
func wallcollide(direction int) {

	switch direction {
	case 1:
		player.direction = 9
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 3
			} else {
				player.direction = 7
			}
		}
	case 2:
		player.direction = 8
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 4
			} else {
				player.direction = 6
			}
		}
	case 3:
		player.direction = 7
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 1
			} else {
				player.direction = 9
			}
		}
	case 4:
		player.direction = 6
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 9
			} else {
				player.direction = 3
			}
		}
	case 6:
		player.direction = 4
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 1
			} else {
				player.direction = 7
			}
		}
	case 7:
		player.direction = 3
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 1
			} else {
				player.direction = 9
			}
		}
	case 8:
		player.direction = 2
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 4
			} else {
				player.direction = 6
			}
		}
	case 9:
		player.direction = 1
		if rolldice() > 4 {
			if flipcoin() {
				player.direction = 3
			} else {
				player.direction = 7
			}
		}
	}

}
func input() { // MARK: input

	if rl.IsKeyPressed(rl.KeyLeftAlt) {
		playervelocity++
	}
	if rl.IsKeyPressed(rl.KeyZ) {
		playervelocity = 1
	}
	if rl.IsKeyPressed(rl.KeyLeftControl) {
		shoot()
	}
	if rl.IsKeyDown(rl.KeyKp7) {
		player.direction = 7
	}
	if rl.IsKeyDown(rl.KeyKp9) {
		player.direction = 9
	}
	if rl.IsKeyDown(rl.KeyKp1) {
		player.direction = 1
	}
	if rl.IsKeyDown(rl.KeyKp3) {
		player.direction = 3
	}
	if rl.IsKeyDown(rl.KeyKp6) {
		player.direction = 6
	}
	if rl.IsKeyDown(rl.KeyKp4) {
		player.direction = 4
	}
	if rl.IsKeyDown(rl.KeyKp2) {
		player.direction = 2
	}
	if rl.IsKeyDown(rl.KeyKp8) {
		player.direction = 8
	}

	if rl.IsKeyPressed(rl.KeyKp0) {
		if gridon {
			gridon = false
		} else {
			gridon = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpDivide) {
		if numberson {
			numberson = false
		} else {
			numberson = true
		}
	}
	if rl.IsKeyPressed(rl.KeyKpMultiply) {
		if centerlineson {
			centerlineson = false
		} else {
			centerlineson = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpAdd) {
		if camera.Zoom == 0.5 {
			camera.Zoom = 1.0
		} else if camera.Zoom == 1.0 {
			camera.Zoom = 2.0
		} else if camera.Zoom == 2.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom == 3.0 {
			camera.Zoom = 4.0
		}
	}
	if rl.IsKeyPressed(rl.KeyKpSubtract) {
		if camera.Zoom == 1.0 {
			camera.Zoom = 0.5
		} else if camera.Zoom == 2.0 {
			camera.Zoom = 1.0
		} else if camera.Zoom == 3.0 {
			camera.Zoom = 2.0
		} else if camera.Zoom == 4.0 {
			camera.Zoom = 3.0
		}
	}
	if rl.IsKeyPressed(rl.KeyPause) {
		if pauseon {
			pauseon = false
		} else {
			pauseon = true
		}
	}

}
func debug() { // MARK: debug
	rl.DrawRectangle(monw-300, 0, 500, monw, rl.Fade(rl.Black, 0.8))
	rl.DrawFPS(monw-290, monh-100)

	playerpositiontext := strconv.Itoa(player.position)
	powerupcounttext := strconv.Itoa(powerupcount)

	rl.DrawText(playerpositiontext, monw-290, 10, 10, rl.White)
	rl.DrawText("playerposition", monw-150, 10, 10, rl.White)
	rl.DrawText(powerupcounttext, monw-290, 20, 10, rl.White)
	rl.DrawText("powerupcount", monw-150, 20, 10, rl.White)

}

// MARK: colors https://www.rapidtables.com/web/color/RGB_Color.html
func randomgrey() rl.Color {
	color := rl.NewColor(uint8(rInt(105, 192)), uint8(rInt(105, 192)), uint8(rInt(105, 192)), 255)
	return color
}
func randombluelight() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 180)), uint8(rInt(120, 256)), uint8(rInt(120, 256)), 255)
	return color
}
func randombluedark() rl.Color {
	color := rl.NewColor(0, 0, uint8(rInt(120, 250)), 255)
	return color
}
func randomyellow() rl.Color {
	color := rl.NewColor(255, uint8(rInt(150, 256)), 0, 255)
	return color
}
func randomorange() rl.Color {
	color := rl.NewColor(uint8(rInt(250, 256)), uint8(rInt(60, 210)), 0, 255)
	return color
}
func randomred() rl.Color {
	color := rl.NewColor(uint8(rInt(128, 256)), uint8(rInt(0, 129)), uint8(rInt(0, 129)), 255)
	return color
}
func randomgreen() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 170)), uint8(rInt(100, 256)), uint8(rInt(0, 50)), 255)
	return color
}
func randomcolor() rl.Color {
	color := rl.NewColor(uint8(rInt(0, 256)), uint8(rInt(0, 256)), uint8(rInt(0, 256)), 255)
	return color
}
func brightyellow() rl.Color {
	color := rl.NewColor(uint8(255), uint8(255), uint8(0), 255)
	return color
}
func brightbrown() rl.Color {
	color := rl.NewColor(uint8(218), uint8(165), uint8(32), 255)
	return color
}
func brightgrey() rl.Color {
	color := rl.NewColor(uint8(212), uint8(212), uint8(213), 255)
	return color
}

// random numbers
func rF32(min, max float32) float32 {
	return (rand.Float32() * (max - min)) + min
}
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
