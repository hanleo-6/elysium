package main

import(
	"fmt"
	"github.com/hanleo-6/elysium/internal/render"
)

func main() {
	c := render.New(50, 50)

	c.SetDot(25, 25)
	c.SetDot(10, 10)

	fmt.Print(c.Render())
}

