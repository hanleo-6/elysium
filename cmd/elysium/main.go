package main

import(
	"fmt"
	"time"
	"math/rand"
	"os"

	"github.com/hanleo-6/elysium/internal/sim"
	"github.com/hanleo-6/elysium/internal/physics"
	"github.com/hanleo-6/elysium/internal/render"

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
	slugs  []*sim.ElysiaChlorotica
	algae  []*sim.Alga
	width  int
	height int
	ready  bool // true if canvas is setup
}

type config struct {
	numSlugs int
	initialAlgae int
}

func getConfig() config {
	return config{
		numSlugs: 2.0,
		initialAlgae: 20,
	}
}

func initialModel() model {
	return model{
		ready: false,
	}
}

func (m *model) initSim() {
	m.initSlugs()
	m.initAlgae()
}

func (m *model) initSlugs() {
	c := getConfig()
	slugs := make([]*sim.ElysiaChlorotica, 0, c.numSlugs)
	for i := 0; i < c.numSlugs; i++ {
		x := int(rand.Float64() * float64(m.canvas.DotWidth()))
		y := int(rand.Float64() * float64(m.canvas.DotHeight()))
		newSlug := sim.NewECSlug(x, y, 2.0)
		slugs = append(slugs, newSlug)
	}
	m.slugs = slugs
}

func (m *model) initAlgae() {
	c := getConfig()
	algae := make([]*sim.Alga, 0, c.initialAlgae)
	for i := 0; i < c.initialAlgae; i++ {
		x := int(rand.Float64() * float64(m.canvas.DotWidth()))
		y := int(rand.Float64() * float64(m.canvas.DotHeight()))
		newAlga := sim.NewAlga(x, y)
		algae = append(algae, newAlga)
	}
	m.algae = algae
}

func (m *model) RegenerateAlgae() {
	for _, alga := range m.algae {
		if alga.Eaten {
			x := int(rand.Float64() * float64(m.canvas.DotWidth()))
			y := int(rand.Float64() * float64(m.canvas.DotHeight()))
			alga.Pos = physics.Vec2{X: float64(x), Y: float64(y)}
			alga.Eaten = false
		}
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

		if !m.ready {
			m.initSim()
			m.ready = true
		}

	case tickMsg:
		// carry out a simulation tick
		dt := 0.033

		for _, slug := range m.slugs {
			slug.Tick(m.algae, dt)
		}

		if len(m.algae) < 10 {
			m.RegenerateAlgae()
		}

		return m, doTick()
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.canvas == nil {
		return tea.NewView("Initialising canvas...")
	}
	m.canvas.Clear()

	for _, alga := range m.algae {
		m.canvas.DrawAlga(alga)
	}
	for _, slug := range m.slugs {
		m.canvas.DrawECSlug(slug)
	}

	return tea.NewView(m.canvas.Render())
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Encountered error running bubbletea on the model: %v", err)
		os.Exit(1)
	}
}

