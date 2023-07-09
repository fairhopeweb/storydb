package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

type tickMsg time.Time

func syncUpdate(msg tea.Msg, m ModelHome) (*ModelHome, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			m.home.ActiveSyncScreen = true
			m.home.Viewport.SetContent(syncView(&m))
			return &m, syncTickCmd()
		case "right":
			m.home.ActiveSyncScreen = false
			m.home.Viewport.SetContent(syncView(&m))
			return &m, nil
		case "enter":
			m.home.StatusSyncScreen = true
		case "q":
			if m.home.StatusSyncScreen {
        m.home.StatusSyncScreen = false
				m.home.ProgressSync = progress.NewModel(progress.WithDefaultGradient())
				m.home.Viewport.SetContent(m.GetDataView())
				return &m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.home.Viewport.Width = msg.Width
		m.home.Viewport.Height = msg.Height
		m.home.ProgressSync.Width = msg.Width - padding*2 - 4
		if m.home.ProgressSync.Width > maxWidth {
			m.home.ProgressSync.Width = maxWidth
		}
		m.home.Viewport.SetContent(syncView(&m))
	case tickMsg:
		if m.home.ProgressSync.Percent() == 1.0 {
			return &m, nil
		}
		Percentage := (float64(3700) / float64(3876)) * 100
		cmd := m.home.ProgressSync.SetPercent(Percentage)
		return &m, tea.Batch(syncTickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.home.ProgressSync.Update(msg)
		m.home.ProgressSync = progressModel.(progress.Model)
		m.home.Viewport.SetContent(syncView(&m))
		return &m, cmd
	}
	return &m, cmd
}

func syncView(m *ModelHome) string {
	var okButton, cancelButton string

	if m.home.ActiveSyncScreen {
		okButton = ActiveButtonStyle.Render("Yes")
		cancelButton = ButtonStyle.Render("No, take me back")
	} else {
		okButton = ButtonStyle.Render("Yes")
		cancelButton = ActiveButtonStyle.Render("No, take me back")
	}

	question := lipgloss.NewStyle().
		Width(m.home.Viewport.Width - 50).
		Align(lipgloss.Center).
		Render("Are you sure you want to sync")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(
		(m.home.Viewport.Width - 50),
		(m.home.Viewport.Height - 50),
		lipgloss.Left, lipgloss.Center,
		DialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(SubtleStyle),
	)
	return BaseStyle.PaddingLeft(20).
		PaddingTop((m.home.Viewport.Height / 2)).
		Render(dialog+"\n\n", syncProgressView(m))
}

func syncProgressView(m *ModelHome) string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.home.ProgressSync.View() + "\n"
}

func syncTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
