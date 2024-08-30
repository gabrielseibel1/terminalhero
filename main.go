package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabrielseibel1/fungo/apply"
)

const (
	boardHeight = 11
	boardWidth  = 50
)

const rollMsg = "roll"

type Model struct {
	board [boardHeight][boardWidth]rune
	sheet [][]rune
	index int
}

func NewModel() *Model {
	m := &Model{}
	return m
}

func (m *Model) Prepare(sheet []string) error {
	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {
			if err := m.Set(i, j, ' '); err != nil {
				return err
			}
		}
	}
	for i := 1; i < boardHeight; i += 2 {
		for j := 0; j < boardWidth; j++ {
			if err := m.Set(i, j, '-'); err != nil {
				return err
			}
		}
	}
	m.sheet = make([][]rune, len(sheet))
	for i := range m.sheet {
		for _, r := range sheet[i] {
			m.sheet[i] = append(m.sheet[i], r)
		}
	}
	return nil
}

func (m *Model) Set(i, j int, r rune) error {
	if !((0 <= i && i < boardHeight) && (0 <= j && j < boardWidth)) {
		return fmt.Errorf("not a valid cell (%d, %d) <- %v", i, j, r)
	}
	m.board[i][j] = r
	return nil
}

func (m *Model) Roll() error {
	if m.sheet == nil || m.index >= len(m.sheet[0]) {
		return nil // finished
	}
	// shift window right
	for i := range m.board {
		for j := 0; j < boardWidth; j++ {
			if j < boardWidth-1 {
				if err := m.Set(i, j, m.board[i][j+1]); err != nil {
					return err
				}
			}
		}
	}
	// roll-in new notes
	for i := range m.board {
		if err := m.Set(i, boardWidth-1, m.sheet[i][m.index]); err != nil {
			return err
		}
	}
	m.index++
	return nil
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String(), tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	case string:
		switch string(msg) {
		case rollMsg:
			if err := m.Roll(); err != nil {
				panic(err)
			}
			return m, nil
		}
	}
	return m, nil
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	lines := apply.ToSlice(m.board[:], func(runesArray [boardWidth]rune) string {
		runesSlice := runesArray[:]
		return string(runesSlice)
	})
	board := strings.Join(lines, "\n")
	footer := fmt.Sprintf("Index = %d", m.index)
	return board + "\n" + footer
}

func main() {
	m := NewModel()
	err := m.Prepare([]string{
		"1              ",
		"-2-------------",
		"  3            ",
		"---4-----------",
		"    5          ",
		"-----6---------",
		"      7       F",
		"-------8-----E-",
		"        9   D  ",
		"---------A-C---",
		"          B    ",
	})
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	go func(ticker *time.Ticker) {
		for range ticker.C {
			p.Send(rollMsg)
		}
	}(time.NewTicker(time.Millisecond * 100))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
