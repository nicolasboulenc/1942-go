package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BONUS_TYPE_MORE     = 1
	BONUS_TYPE_RATE     = 2
	BONUS_TYPE_SPEED    = 3
	BONUS_MVT_OSCILLATE = 1
	BONUS_MVT_STRAIGHT  = 2
	DEBUG_CONF_INTERVAL = 0.5 // check every x seconds for config file change
)

type Player struct {
	pos            rl.Vector2
	dir            rl.Vector2
	speed          float32 // pixel per second
	fire_rate      int     // shots per second
	fire_prev_time float64
	is_firing      bool
}

type Config struct {
	player_velocity     float32
	player_fire_rate    int
	projectile_velocity float32
	bonus_velocity      float32
}

type Bonus struct {
	type_ int
	pos   rl.Vector2
	speed float32
}

type Projectile struct {
	enabled bool
	pos     rl.Vector2
	dir     rl.Vector2
}

type GameState struct {
	in_menu   bool
	is_paused bool
	prev_time float64
}

type Debug struct {
	conf_prev_time       time.Time // the time the file was changed
	conf_check_prev_time float64   // last time we checked
}

var player Player = Player{pos: rl.Vector2{X: 100, Y: 100}, dir: rl.Vector2{X: 0, Y: 0}, speed: config.player_velocity, fire_rate: 16}
var game_state GameState = GameState{in_menu: true, is_paused: true}
var projectiles = make([]Projectile, 10)
var bonus Bonus
var debug Debug
var config Config

func main() {
	screenWidth := int32(400)
	screenHeight := int32(400)

	debug.conf_prev_time = config_check()
	config_load()

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(screenWidth, screenHeight, "1942-go")
	rl.SetWindowPosition(10, 10)

	buffer := rl.LoadRenderTexture(screenWidth, screenHeight)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading

	rl.SetTargetFPS(60)

	game_state.prev_time = rl.GetTime()
	player.fire_prev_time = game_state.prev_time

	bonus = Bonus{type_: BONUS_TYPE_MORE, pos: rl.Vector2{}, speed: config.bonus_velocity}

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

		if debug.conf_check_prev_time+DEBUG_CONF_INTERVAL > ntime {
			current_time := config_check()
			if debug.conf_prev_time != current_time {
				config_load()
			}
			debug.conf_check_prev_time = ntime
		}

		for i := 0; i < len(projectiles); i++ {

			projectiles[i].pos.Y -= config.projectile_velocity * float32(dtime)

			proj_rect := rl.Rectangle{X: projectiles[i].pos.X, Y: projectiles[i].pos.Y, Width: 1, Height: 10}
			view_rect := rl.Rectangle{X: 0, Y: 0, Width: float32(screenWidth), Height: float32(screenHeight)}
			if !rl.CheckCollisionRecs(proj_rect, view_rect) {
				projectiles[i].enabled = false
			}
		}

		// draw to texture
		rl.BeginTextureMode(buffer)
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, int32(player.pos.X)-texture.Width/2, int32(player.pos.Y)-texture.Height/2, rl.White)

		for i := 0; i < len(projectiles); i++ {
			if projectiles[i].enabled {
				rl.DrawRectangle(int32(projectiles[i].pos.X), int32(projectiles[i].pos.Y), 2, 6, rl.Black)
			}
		}

		rl.EndTextureMode()

		// draw to screen
		rl.BeginDrawing()
		rl.DrawTextureRec(buffer.Texture, rl.Rectangle{X: 0, Y: 0, Width: float32(buffer.Texture.Width), Height: float32(-buffer.Texture.Height)}, rl.Vector2{X: 0, Y: 0}, rl.White)
		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

func process_inputs() {

	player.dir = rl.Vector2Zero()
	player.is_firing = false

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

func config_check() time.Time {
	fileInfo, err := os.Stat("1942.conf")
	if err != nil {
		log.Fatal(err)
	}

	return fileInfo.ModTime()
}

func config_load() {

	file, err := os.Open("1942.conf")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanWords)

	tt := reflect.ValueOf(config)
	field_name := ""
	for fileScanner.Scan() {
		if field_name == "" {
			field_name = fileScanner.Text()
		} else {
			value_text := fileScanner.Text()
			if field_name == "player_velocity" {
				config.player_velocity, _ = strconv.ParseFloat(value_text, 32)
			} else if field_name == "player_fire_rate" {
				config.player_velocity, _ = strconv.ParseFloat(value_text, 32)
			} else if field_name == "projectile_velocity" {
				config.player_velocity, _ = strconv.ParseFloat(value_text, 32)
			} else if field_name == "bonus_velocity" {
				config.player_velocity, _ = strconv.ParseFloat(value_text, 32)
			}
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
