package entity

import "github.com/charmbracelet/bubbles/viewport"

type Command struct {
	Content  string
	Cursor   int
	Selected map[int]struct{}
	Ready    bool
	Viewport viewport.Model
}

func NewCmd() *Command {
	m := Command{}

	return &m
}
