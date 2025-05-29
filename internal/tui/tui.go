package tui

import (
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

func NewModel() mainModel {
	playerModel := playerModel{
		trackPosition: 0,
		progressBar:   progress.New(progress.WithDefaultGradient(), progress.WithWidth(50)),
	}

	return mainModel{
		state: PLAYER_STATE, playerModel: playerModel,
	}
}

func (m *mainModel) Init() tea.Cmd {
	// TODO: swap this for actual logic loading latest book from the database
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
