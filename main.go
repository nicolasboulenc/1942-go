package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// movement types
	MOVEMENT_TYPE_OSCILLATE = 1
	MOVEMENT_TYPE_STRAIGHT  = 2
	// weapon fire types
	FIRE_TYPE_ALTERNATE = 1
	FIRE_TYPE_DUAL      = 2
	FIRE_TYPE_QUAD      = 3
	// debug stuff
	DEBUG_CONF_INTERVAL = 0.5 // check every x seconds for config file change
)

type Player struct {
	pos                rl.Vector2
	dir                rl.Vector2
	velocity_modifier  float32 // 1.2x
	fire_rate_modifier float64 // 2x
	is_firing          bool
	weapons            []Weapon
}

type Config struct {
	player_velocity     float32
	player_fire_rate    int
	projectile_velocity float32
	bonus_velocity      float32
}

type Projectile struct {
	enabled  bool
	pos      rl.Vector2
	dir      rl.Vector2
	velocity float32
}

type GameState struct {
	screen_width  int32
	screen_height int32
	in_menu       bool
	is_paused     bool
	prev_time     float64
}

type Debug struct {
	conf_prev_time       time.Time // the time the file was changed
	conf_check_prev_time float64   // last time we checked
}

type Weapon struct {
	projectile_velocity float32 // 1.2x
	projectile_damage   float32
	state               int
	fire_type           int     // single, burst, auto, etc
	fire_rate           float64 // 2x
	fire_prev_time      float64
	is_firing           bool
}

var player Player
var game_state GameState
var projectiles []Projectile
var debug Debug
var config Config

func main() {

	game_init()

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(game_state.screen_width, game_state.screen_height, "1942-go")
	rl.SetWindowPosition(10, 10)

	buffer := rl.LoadRenderTexture(game_state.screen_width, game_state.screen_height)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		ntime := rl.GetTime()
		dtime := ntime - game_state.prev_time
		game_state.prev_time = ntime

		// inputs
		process_inputs()
		player.pos.X += player.dir.X * config.player_velocity * player.velocity_modifier * float32(dtime)
		player.pos.Y += player.dir.Y * config.player_velocity * player.velocity_modifier * float32(dtime)

		// logic
		if player.is_firing {
			for i := 0; i < len(player.weapons); i++ {
				weapon := &player.weapons[i]
				if ntime-weapon.fire_prev_time >= 1/weapon.fire_rate {
					projectile := player_weapon_fire(weapon, player)
					projectiles = append(projectiles, projectile)
				}
			}
		}

		for i := 0; i < len(projectiles); i++ {

			projectiles[i].pos.Y -= config.projectile_velocity * float32(dtime)

			proj_rect := rl.Rectangle{X: projectiles[i].pos.X, Y: projectiles[i].pos.Y, Width: 1, Height: 10}
			view_rect := rl.Rectangle{X: 0, Y: 0, Width: float32(game_state.screen_width), Height: float32(game_state.screen_height)}
			if !rl.CheckCollisionRecs(proj_rect, view_rect) {
				projectiles[i].enabled = false
			}
		}

		if debug.conf_check_prev_time+DEBUG_CONF_INTERVAL > ntime {
			current_time := config_check()
			if debug.conf_prev_time != current_time {
				config_load()
			}
			debug.conf_check_prev_time = ntime
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

func game_init() {

	player = Player{pos: rl.Vector2{X: 100, Y: 100}, dir: rl.Vector2{X: 0, Y: 0}, velocity_modifier: 1, fire_rate_modifier: 1, is_firing: false, weapons: make([]Weapon, 0, 3)}
	game_state = GameState{in_menu: true, is_paused: true, screen_width: 640, screen_height: 480}
	projectiles = make([]Projectile, 0, 10)

	weapon := Weapon{projectile_velocity: 400, projectile_damage: 2, state: 0, fire_type: FIRE_TYPE_ALTERNATE, fire_rate: 16, fire_prev_time: 0, is_firing: false}
	player_weapon_add(weapon)

	debug.conf_prev_time = config_check()
	config_load()

	game_state.prev_time = rl.GetTime()
}

func player_weapon_add(weapon Weapon) {

	// check the weapon doesn already exist
	for i := 0; i < len(player.weapons); i++ {
		if player.weapons[i].fire_type == weapon.fire_type {
			// remove that weapon
			player.weapons = append(player.weapons[:i], player.weapons[i+1:]...)
			break
		}
	}

	player.weapons = append(player.weapons, weapon)
}

func player_weapon_fire(weapon *Weapon, player Player) Projectile {

	var origin rl.Vector2
	origin.Y = player.pos.Y

	switch weapon.fire_type {
	case FIRE_TYPE_ALTERNATE:
		if weapon.state == 0 {
			weapon.state = 1
			origin.X = player.pos.X - 10
		} else {
			weapon.state = 0
			origin.X = player.pos.X + 10
		}
		weapon.fire_prev_time = rl.GetTime()
		return Projectile{enabled: true, pos: origin, velocity: weapon.projectile_velocity}
	}
	return Projectile{enabled: false}
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
	if rl.IsKeyDown(rl.KeySpace) {
		player.is_firing = true
	}
	if rl.IsKeyDown(rl.KeyU) {
		weapon := Weapon{projectile_velocity: 400, projectile_damage: 2, fire_type: FIRE_TYPE_ALTERNATE, fire_rate: 400, fire_prev_time: 0, is_firing: false}
		player.weapons = append(player.weapons, weapon)
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

	field_name := ""
	for fileScanner.Scan() {
		if field_name == "" {
			field_name = fileScanner.Text()
		} else {
			value_text := fileScanner.Text()
			switch field_name {
			case "player_velocity":
				val, _ := strconv.ParseInt(value_text, 10, 32)
				config.player_velocity = float32(val)
			case "player_fire_rate":
				val, _ := strconv.ParseInt(value_text, 10, 32)
				config.player_fire_rate = int(val)
			case "projectile_velocity":
				val, _ := strconv.ParseInt(value_text, 10, 32)
				config.projectile_velocity = float32(val)
			case "bonus_velocity":
				val, _ := strconv.ParseInt(value_text, 10, 32)
				config.bonus_velocity = float32(val)
			}
			field_name = ""
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
