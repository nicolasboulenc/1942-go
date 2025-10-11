package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
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
	Window struct {
		Title string `json:"title"`
	}
	Player struct {
		Velocity float32 `json:"velocity"`
	}
	Weapons []WeaponConfig `json:"weapons"`
}

type WeaponConfig struct {
	Name                string  `json:"name"`
	FireType            int     `json:"fire_type"`
	FireRate            float64 `json:"fire_rate"`
	ProjectileVelocity  float32 `json:"projectile_velocity"`
	ProjectileDamage    float32 `json:"projectile_damage"`
	ProjectileDirection Vector2 `json:"projectile_direction"`
}

type Vector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
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
	conf           WeaponConfig
	state          int
	fire_prev_time float64
	is_firing      bool
}

var player Player
var game_state GameState
var projectiles []Projectile
var debug Debug
var config Config
var tileMap *TileMap
var tileSet *TileSet

func main() {

	gameInit()

	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(game_state.screen_width, game_state.screen_height, "1942-go")
	rl.SetWindowPosition(10, 10)

	buffer := rl.LoadRenderTexture(game_state.screen_width, game_state.screen_height)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading
	atlas := rl.LoadTexture(tileSet.Image)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		ntime := rl.GetTime()
		dtime := ntime - game_state.prev_time
		game_state.prev_time = ntime

		// inputs
		gameProcessInputs()
		player.pos.X += player.dir.X * config.Player.Velocity * player.velocity_modifier * float32(dtime)
		player.pos.Y += player.dir.Y * config.Player.Velocity * player.velocity_modifier * float32(dtime)

		// logic
		if player.is_firing {
			for i := 0; i < len(player.weapons); i++ {
				weapon := &player.weapons[i]
				if ntime-weapon.fire_prev_time >= 1/weapon.conf.FireRate {
					playerWeaponFire(weapon, player, &projectiles)
				}
			}
		}

		view_rect := rl.Rectangle{X: 0, Y: 0, Width: float32(game_state.screen_width), Height: float32(game_state.screen_height)}
		for i := 0; i < len(projectiles); i++ {

			projectiles[i].pos.X += projectiles[i].velocity * projectiles[i].dir.X * float32(dtime)
			projectiles[i].pos.Y -= projectiles[i].velocity * projectiles[i].dir.Y * float32(dtime)
			proj_rect := rl.Rectangle{X: projectiles[i].pos.X, Y: projectiles[i].pos.Y, Width: 1, Height: 10}
			if !rl.CheckCollisionRecs(proj_rect, view_rect) {
				projectiles[i].enabled = false
			}
		}

		if debug.conf_check_prev_time+DEBUG_CONF_INTERVAL > ntime {
			latest_time := configFileCheck()
			if debug.conf_prev_time != latest_time {
				configFileLoad()
			}
			debug.conf_check_prev_time = ntime
		}

		// draw to texture
		rl.BeginTextureMode(buffer)
		rl.ClearBackground(rl.RayWhite)

		src := rl.Rectangle{X: 0, Y: 0, Width: float32(tileMap.TileWidth), Height: float32(tileMap.TileHeight)}
		dst := rl.Vector2{X: 0, Y: 0}
		for i := 0; i < len(tileMap.Layers[0].Data); i++ {
			tileId := tileMap.Layers[0].Data[i] - 1
			src.X = float32(tileId % tileSet.Columns * tileSet.TileWidth)
			src.Y = float32(math.Floor(float64(tileId)/float64(tileSet.Columns)) * float64(tileSet.TileHeight))
			dst.X = float32(i % tileMap.Width * tileMap.TileWidth)
			dst.Y = float32(math.Floor(float64(i)/float64(tileMap.Width)) * float64(tileMap.TileHeight))
			rl.DrawTextureRec(atlas, src, dst, rl.White)
		}

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

func gameInit() {

	debug.conf_prev_time = configFileCheck()
	configFileLoad()

	tileMap = LoadTileMap("map.json")
	fmt.Printf("%+v\n", tileMap)

	tileSet = LoadTileSet("tileset.json")
	fmt.Printf("%+v\n", tileSet)

	player = Player{pos: rl.Vector2{X: 100, Y: 100}, dir: rl.Vector2{X: 0, Y: 0}, velocity_modifier: 1, fire_rate_modifier: 1, is_firing: false, weapons: make([]Weapon, 0, 3)}
	game_state = GameState{in_menu: true, is_paused: true, screen_width: 800, screen_height: 600}
	projectiles = make([]Projectile, 0, 10)

	weaponConf := configWeaponGet(&config, "alternate")
	// weaponConf := WeaponConfig{Name: "basic", FireType: FIRE_TYPE_ALTERNATE, FireRate: 8, ProjectileVelocity: 400, ProjectileDamage: 2}
	weapon := Weapon{conf: *weaponConf, state: 0, fire_prev_time: 0, is_firing: false}
	playerWeaponAdd(weapon)

	game_state.prev_time = rl.GetTime()
}

func gameProcessInputs() {

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
		configWeapon := configWeaponGet(&config, "quad")
		playerWeaponSwap(configWeapon)
	}
	player.dir = rl.Vector2Normalize(player.dir)
}

func playerWeaponAdd(weapon Weapon) {

	// check the weapon doesn already exist
	for i := 0; i < len(player.weapons); i++ {
		if player.weapons[i].conf.FireType == weapon.conf.FireType {
			// remove that weapon
			player.weapons = append(player.weapons[:i], player.weapons[i+1:]...)
			break
		}
	}

	player.weapons = append(player.weapons, weapon)
}

func playerWeaponSwap(weaponConfig *WeaponConfig) {

	player.weapons[0].conf = *weaponConfig
}

func playerWeaponFire(weapon *Weapon, player Player, projectiles *[]Projectile) {

	var origin rl.Vector2
	origin.Y = player.pos.Y

	switch weapon.conf.FireType {
	case FIRE_TYPE_ALTERNATE:
		weapon.fire_prev_time = rl.GetTime()
		if weapon.state == 0 {
			weapon.state = 1
			origin.X = player.pos.X - 10
		} else {
			weapon.state = 0
			origin.X = player.pos.X + 10
		}
		projectile := Projectile{enabled: true, pos: origin, dir: rl.Vector2(weapon.conf.ProjectileDirection), velocity: weapon.conf.ProjectileVelocity}
		*projectiles = append(*projectiles, projectile)
	case FIRE_TYPE_DUAL:
		weapon.fire_prev_time = rl.GetTime()

		origin.X = player.pos.X - 10
		projectile := Projectile{enabled: true, pos: origin, dir: rl.Vector2(weapon.conf.ProjectileDirection), velocity: weapon.conf.ProjectileVelocity}
		*projectiles = append(*projectiles, projectile)

		origin.X = player.pos.X + 10
		projectile = Projectile{enabled: true, pos: origin, dir: rl.Vector2(weapon.conf.ProjectileDirection), velocity: weapon.conf.ProjectileVelocity}
		*projectiles = append(*projectiles, projectile)
	case FIRE_TYPE_QUAD:
		weapon.fire_prev_time = rl.GetTime()

		origin.X = player.pos.X - 15
		projectile := Projectile{enabled: true, pos: origin, dir: rl.Vector2(weapon.conf.ProjectileDirection), velocity: weapon.conf.ProjectileVelocity}
		*projectiles = append(*projectiles, projectile)

		origin.X = player.pos.X + 15
		dir2 := rl.Vector2{X: -weapon.conf.ProjectileDirection.X, Y: weapon.conf.ProjectileDirection.Y}
		projectile = Projectile{enabled: true, pos: origin, dir: dir2, velocity: weapon.conf.ProjectileVelocity}
		*projectiles = append(*projectiles, projectile)
	}
}

func configWeaponGet(config *Config, name string) *WeaponConfig {

	for i := 0; i < len(config.Weapons); i++ {
		if config.Weapons[i].Name == name {
			return &config.Weapons[i]
		}
	}
	return nil
}

func configFileCheck() time.Time {
	fileInfo, err := os.Stat("config.json")
	if err != nil {
		log.Fatal(err)
	}

	return fileInfo.ModTime()
}

func configFileLoad() {

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config)
}
