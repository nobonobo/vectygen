package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

// NewSample ...
func NewSample(d map[string]func(*vecty.Event)) *Sample {
	return &Sample{
		dispatcher: d,
	}
}

// Sample ...
type Sample struct {
	vecty.Core
	dispatcher map[string]func(*vecty.Event)
}

// Render ...
func (c *Sample) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(
			vecty.Property("", "true"),
		),
		elem.Input(
			vecty.Markup(
				vecty.ClassMap{
					"boo1": true,
					"boo2": true,
					"boo3": true,
					"boo4": true,
					"boo5": true,
				},
				prop.Disabled(true),
			),
			vecty.Text("Hello"),
			elem.Break(
				vecty.Markup(
					vecty.Class("hoge"),
				),
			),
			vecty.Text("World!"),
		),
		elem.Button(
			vecty.Markup(
				event.Click(c.Click),
			),
			vecty.Text("Click"),
		),
	)
}

// Click ...
func (c *Sample) Click(event *vecty.Event) {
	f, ok := c.dispatcher["Click"]
	if !ok {
		panic("unknown func: \"Click\"")
	}
	f(event)
}
