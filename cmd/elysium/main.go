package main

import(
	"fmt"
	"time"
	"math/rand"

	"github.com/hanleo-6/elysium/internal/render"
	"github.com/hanleo-6/elysium/internal/physics"
)

func main() {
	c := render.New(50, 25)

	// c.DrawCircle(c.DotWidth()/3, c.DotHeight()/2, 3.0, render.ColorAdult)
	// c.DrawLine(0, 0, 5, 2, render.ColorDeepWater)
	skelly := physics.NewSpine(physics.Vec2{10, 10}, 10, 3.0)

	fmt.Print("\033[2J\033[?25l") // clear terminal screen + hide cursor
	defer fmt.Print("\033[?25h") // restore cursor when done

	dt := 0.033
	// setup randomly switching target
	target := physics.Vec2{X: 40, Y: 20}
	timer := 0.0

	for {
		c.Clear()
		c.DrawSpine(skelly, render.ColorAdult)
		c.SetDotWithColor(int(target.X), int(target.Y), render.ColorLightBand)

		accel := skelly.Seek(target, dt)
		skelly.Step(dt, accel)

		timer += dt
		if timer > 6.0 {
			timer = 0.0
			target = physics.Vec2{
				X: rand.Float64() * float64(c.DotWidth()),
				Y: rand.Float64() * float64(c.DotHeight()),
			}
		}

		fmt.Print("\033[H") // move cursor to top left again
		fmt.Print(c.Render())
		time.Sleep(33 * time.Millisecond)
	}
}

