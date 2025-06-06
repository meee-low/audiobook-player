package tui

import (
	"fmt"

	"github.com/gopxl/beep/v2"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type playerModel struct {
	trackTitle    string
	trackLength   int // in seconds
	trackPosition int // in seconds
	streamer      beep.StreamSeekCloser
	format        beep.Format
	progressBar   progress.Model
}

func (m playerModel) Update(msg tea.Msg) (playerModel, tea.Cmd) {
	// FIX: The progress bar looks wonky because it flips in irregular intervals.
	switch msg.(type) {
	case tickMsg:
		// TODO: use m.streamer.Position() and m.streamer.trackLength() instead
		if m.trackPosition < m.trackLength {
			m.trackPosition++
			percentage := float64(m.trackPosition) / float64(m.trackLength)
			cmd := m.progressBar.SetPercent(percentage)
			return m, tea.Batch(cmd, tick())
		} else {
			return m, tea.Quit
		}

	case progress.FrameMsg:
		p, cmd := m.progressBar.Update(msg)
		progressBar, ok := p.(progress.Model)
		if !ok {
			panic("Could not cast this tea.Model to progress.Model, but it really should be a progress.Model, so something very weird happened")
		}
		m.progressBar = progressBar
		return m, cmd

	default:
		return m, tea.Quit
	}

}

func (m *playerModel) View() string {
	return fmt.Sprintf("Now playing: %s\n%s\n %d:%02d / %d:%02d\n\n",
		m.trackTitle,
		m.progressBar.View(),
		m.trackPosition/60, m.trackPosition%60,
		m.trackLength/60, m.trackLength%60,
	)
}
