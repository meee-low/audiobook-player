package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type bubbleteaState int
type tickMsg time.Time

const (
	PLAYER_STATE bubbleteaState = iota
)

type mainModel struct {
	state       bubbleteaState
	playerModel playerModel
}

type playerModel struct {
	trackTitle    string
	trackLength   int // in seconds
	trackPosition int // in seconds
	streamer      beep.StreamSeekCloser
	format        beep.Format
	progressBar   progress.Model
}

func newModel() mainModel {
	playerModel := playerModel{
		trackPosition: 0,
		progressBar:   progress.New(progress.WithDefaultGradient(), progress.WithWidth(50)),
	}

	return mainModel{
		state: PLAYER_STATE, playerModel: playerModel,
	}
}

func (m *mainModel) Init() tea.Cmd {
	// Get the file from the command line arguments
	if len(os.Args) < 2 {
		tea.Quit()
		log.Fatalf("No file provided. Usage: go run main.go <filename>\n")
	}

	fileName := os.Args[1]

	audiofile, err := os.Open(fileName)
	if err != nil {
		tea.Quit()
		log.Fatalf("Could not open file: %s\n", fileName)
	}

	m.playerModel.trackTitle = filepath.Base(fileName)

	// Initialize the streamer and the speaker

	streamer, format, err := mp3.Decode(audiofile)
	if err != nil {
		tea.Quit()
		log.Fatalf("Error decoding MP3 file: %v", err)
	}

	m.playerModel.trackLength = streamer.Len() / format.SampleRate.N(time.Second)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		tea.Quit()
		audiofile.Close()
		streamer.Close()
	})))

	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *mainModel) Close() {
	m.playerModel.streamer.Close()
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case PLAYER_STATE:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				return m, nil
			}
		case tickMsg, progress.FrameMsg:
			var cmd tea.Cmd
			m.playerModel, cmd = m.playerModel.Update(msg)
			return m, cmd
		case tea.QuitMsg:
			tea.Quit()
			return m, nil
		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

func (m *mainModel) View() string {
	switch m.state {
	case PLAYER_STATE:
		return m.playerModel.View()
	default:
		return "Unknown state"
	}
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

func main() {
	mainModel := newModel()
	program := tea.NewProgram(&mainModel)
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running bubbletea program %v", err)
	}
}
