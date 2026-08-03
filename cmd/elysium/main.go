package main

import(
	"fmt"
	"time"
	"math/rand"
	"os"

	"github.com/hanleo-6/elysium/internal/render"
	"github.com/hanleo-6/elysium/internal/physics"
	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(33 * time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	canvas *render.Canvas
	skelly *physics.Skeleton
	target physics.Vec2
	timer  float64
	width  int
	height int
}

func initialModel() model {
	return model {
		skelly: physics.NewSpine(physics.Vec2{X: 10, Y: 10}, 10, 3.0),
		target: physics.Vec2{X: 20, Y: 20},
		timer: 3.0,
	}
}

func (m model) Init() tea.Cmd {
	return doTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// resize canvas if terminal is resized
		m.width = msg.Width
		m.height = msg.Height
		m.canvas = render.NewCanvas(m.width, m.height)

	case tickMsg:
		// carry out a simulation tick
		dt := 0.033

		m.timer += dt
		if m.timer > 6.0 {
			m.timer = 0.0
			m.target = physics.Vec2{
				X: rand.Float64() * float64(m.canvas.DotWidth()),
				Y: rand.Float64() * float64(m.canvas.DotHeight()),
			}
		}
		accel := m.skelly.Seek(m.target, dt)
		m.skelly.Step(dt, accel)

		return m, doTick()
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.canvas == nil {
		return tea.NewView("Initialising canvas...")
	}
	m.canvas.Clear()

	m.canvas.DrawSpine(m.skelly, render.ColorAdult)
	m.canvas.SetDotWithColor(int(m.target.X), int(m.target.Y), render.ColorLightBand)

	return tea.NewView(m.canvas.Render())
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Encountered error running bubbletea on the model: %v", err)
		os.Exit(1)
	}
}

