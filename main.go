package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PLAYER_SPEED = 200 // pixels per second
)

type Player struct {
	pos            rl.Vector2
	dir            rl.Vector2
	speed          float32 // pixel per second
	fire_rate      int     // shots per second
	fire_prev_time float64
	fire_vel       float32
	is_firing      bool
}

type Projectile struct {
	enabled bool
	pos     rl.Vector2
	dir     rl.Vector2
}

type State struct {
	in_menu   bool
	is_paused bool
	prev_time float64
}

var player Player = Player{pos: rl.Vector2{X: 100, Y: 100}, dir: rl.Vector2{X: 0, Y: 0}, speed: PLAYER_SPEED}
var game_state State = State{in_menu: true, is_paused: true}
var projectiles = make([]Projectile, 0, 10)

func main() {
	screenWidth := int32(400)
	screenHeight := int32(400)

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(screenWidth, screenHeight, "1942-go")
	rl.SetWindowPosition(10, 10)

	buffer := rl.LoadRenderTexture(screenWidth, screenHeight)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading

	rl.SetTargetFPS(60)

	game_state.prev_time = rl.GetTime()
	player.fire_prev_time = game_state.prev_time

	for !rl.WindowShouldClose() {

		ntime := rl.GetTime()
		dtime := ntime - game_state.prev_time
		game_state.prev_time = ntime

		// inputs
		process_inputs()
		player.pos.X += player.dir.X * player.speed * float32(dtime)
		player.pos.Y += player.dir.Y * player.speed * float32(dtime)

		// logic
		if player.is_firing && ntime-player.fire_prev_time >= 1/float64(player.fire_rate) {
			projectile := Projectile{enabled: true, pos: rl.Vector2{X: player.pos.X, Y: player.pos.Y}}
			projectiles = append(projectiles, projectile)
			player.fire_prev_time = ntime
		}

		for i := 0; i < len(projectiles); i++ {
		}

		// draw to texture
		rl.BeginTextureMode(buffer)
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, int32(player.pos.X)-texture.Width/2, int32(player.pos.Y)-texture.Height/2, rl.White)

		for i := 0; i < len(projectiles); i++ {
			rl.DrawRectangle(int32(projectiles[i].pos.X), int32(projectiles[i].pos.Y), 10, 10, rl.Black)
		}

		rl.EndTextureMode()

		// draw to screen
		rl.BeginDrawing()
		rl.DrawTextureRec(buffer.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(buffer.Texture.Width), Height: float32(-buffer.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.DrawText(fmt.Sprintf("this IS a texture!", player.is_firing), 10, 10, 10, rl.Gray)
		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

func process_inputs() {

	player.dir = rl.Vector2Zero()

	// keyboard
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		player.dir.X += 0.8
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		player.dir.X -= 0.8
	}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		player.dir.Y -= 0.8
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		player.dir.Y += 0.8
	}
	if rl.IsKeyDown(rl.KeySpace) {
		player.is_firing = true
	}

	player.dir = rl.Vector2Normalize(player.dir)
}
