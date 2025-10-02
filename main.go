package main

import (
	"fmt"

	raygui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PLAYER_SPEED = 200
)

type Player struct {
	pos   rl.Vector2
	dir   rl.Vector2
	speed float32
}

type State struct {
	in_menu   bool
	is_paused bool
}

var player Player = Player{pos: rl.Vector2{X: 100, Y: 100}, dir: rl.Vector2{X: 0, Y: 0}, speed: PLAYER_SPEED}
var game_state State = State{in_menu: true, is_paused: true}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(400)

	rl.InitWindow(screenWidth, screenHeight, "1942-go")
	scale_dpi := rl.GetWindowScaleDPI()
	fmt.Println("scale %i %i", scale_dpi.X, scale_dpi.Y)

	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		dtime := rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if game_state.in_menu {
			if raygui.Button(rl.NewRectangle(350, 200, 100, 30), "Start Game") {
				game_state.in_menu = false
			}
		} else {

			process_inputs()
			player.pos.X += player.dir.X * player.speed * dtime
			player.pos.Y += player.dir.Y * player.speed * dtime

			rl.DrawTexture(texture, int32(player.pos.X)-texture.Width/2, int32(player.pos.Y)-texture.Height/2, rl.White)
			rl.DrawText("this IS a texture!", 360, 370, 10, rl.Gray)
		}

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
	player.dir = rl.Vector2Normalize(player.dir)

	// gamepad
	var gamepad int32 = 0 // which gamepad to display
	// fmt.Println(rl.GetGamepadName(gamepad))
	if rl.IsGamepadAvailable(gamepad) {
		rl.DrawText(fmt.Sprintf("GP1: %s", rl.GetGamepadName(gamepad)), 10, 10, 10, rl.Black)

		// Draw buttons: xbox home
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddle) {
			rl.DrawText("Button: Middle", 10, 25, 10, rl.Black)
		}

		// Draw buttons: basic
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleRight) {
			rl.DrawText("Button: START", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonMiddleLeft) {
			rl.DrawText("Button: BACK", 10, 25, 10, rl.Black)
		}

		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceLeft) {
			rl.DrawText("Button: X", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceDown) {
			rl.DrawText("Button: A", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceRight) {
			rl.DrawText("Button: B", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightFaceUp) {
			rl.DrawText("Button: Y", 10, 25, 10, rl.Black)
		}

		// Draw buttons: d-pad
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceUp) {
			rl.DrawText("Button: D-UP", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceDown) {
			rl.DrawText("Button: D-DOWN", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceLeft) {
			rl.DrawText("Button: D-LEFT", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftFaceRight) {
			rl.DrawText("Button: D-RIGHT", 10, 25, 10, rl.Black)
		}

		// Draw buttons: left-right back
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonLeftTrigger1) {
			rl.DrawText("Button: Shoulder-L", 10, 25, 10, rl.Black)
		}
		if rl.IsGamepadButtonDown(gamepad, rl.GamepadButtonRightTrigger1) {
			rl.DrawText("Button: Shoulder-R", 10, 25, 10, rl.Black)
		}

		// Draw axis: left joystick
		// rl.DrawCircle(259, 152, 39, rl.Black)
		// rl.DrawCircle(259, 152, 34, rl.LightGray)
		// rl.DrawCircle(int32(259+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftX)*20)),
		// 	int32(152-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftY)*20)), 25, rl.Black)

		// Draw axis: right joystick
		// rl.DrawCircle(461, 237, 38, rl.Black)
		// rl.DrawCircle(461, 237, 33, rl.LightGray)
		// rl.DrawCircle(int32(461+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightX)*20)),
		// 	int32(237-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightY)*20)), 25, rl.Black)

		// Draw axis: left-right triggers
		// rl.DrawRectangle(170, 30, 15, 70, rl.Gray)
		// rl.DrawRectangle(604, 30, 15, 70, rl.Gray)
		// rl.DrawRectangle(170, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftTrigger))/2.0)*70), rl.Red)
		// rl.DrawRectangle(604, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightTrigger))/2.0)*70), rl.Red)
	}
}
