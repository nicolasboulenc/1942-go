package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PLAYER_SPEED = 2
)

type Player struct {
	x   int
	y   int
	dir float32
}

var player Player = Player{x: 400, y: 225, dir: 0.0}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "gamesafoot")
	// NOTE: Textures MUST be loaded after Window initialization (OpenGL context is required)
	texture := rl.LoadTexture("images/UK_Spitfire.png") // Texture loading

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		process_inputs()

		rl.DrawTexture(texture, screenWidth/2-texture.Width/2, screenHeight/2-texture.Height/2, rl.White)
		rl.DrawText("this IS a texture!", 360, 370, 10, rl.Gray)

		rl.EndDrawing()
	}

	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

func process_inputs() {

	rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2)

	// keyboard
	if rl.IsKeyDown(rl.KeyRight) {
		ballPosition.X += 0.8
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		ballPosition.X -= 0.8
	}
	if rl.IsKeyDown(rl.KeyUp) {
		ballPosition.Y -= 0.8
	}
	if rl.IsKeyDown(rl.KeyDown) {
		ballPosition.Y += 0.8
	}

	// gamepad
	var gamepad int32 = 0 // which gamepad to display
	fmt.Println(rl.GetGamepadName(gamepad))
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
		rl.DrawCircle(259, 152, 39, rl.Black)
		rl.DrawCircle(259, 152, 34, rl.LightGray)
		rl.DrawCircle(int32(259+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftX)*20)),
			int32(152-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftY)*20)), 25, rl.Black)

		// Draw axis: right joystick
		rl.DrawCircle(461, 237, 38, rl.Black)
		rl.DrawCircle(461, 237, 33, rl.LightGray)
		rl.DrawCircle(int32(461+(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightX)*20)),
			int32(237-(rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightY)*20)), 25, rl.Black)

		// Draw axis: left-right triggers
		rl.DrawRectangle(170, 30, 15, 70, rl.Gray)
		rl.DrawRectangle(604, 30, 15, 70, rl.Gray)
		rl.DrawRectangle(170, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisLeftTrigger))/2.0)*70), rl.Red)
		rl.DrawRectangle(604, 30, 15, int32(((1.0+rl.GetGamepadAxisMovement(gamepad, rl.GamepadAxisRightTrigger))/2.0)*70), rl.Red)
	}
}
