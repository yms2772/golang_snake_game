package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/mattn/go-runewidth"

	"github.com/nsf/termbox-go"
)

type Move int

type Snake struct {
	X int
	Y int
}

var (
	gameOver = make(chan bool)
	field    [60][30]int
	snake    []*Snake
	move     = Right
	prevMove = move
)

const (
	Road int = iota
	Item
)

const (
	Up Move = iota
	Down
	Right
	Left
)

func SetBlock(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
	_ = termbox.Flush()
}

func PrintField() {
	for y := 0; y < 30; y++ {
		for x := 0; x < 60; x++ {
			switch field[x][y] {
			case Road:
				SetBlock(x, y, termbox.ColorDefault, termbox.ColorDefault, ".")
			case Item:
				SetBlock(x, y, termbox.ColorGreen, termbox.ColorDefault, "O")
			}
		}
	}

	for _, item := range snake {
		SetBlock(item.X, item.Y, termbox.ColorRed, termbox.ColorDefault, "O")
	}
}

func GetHeadLocation() (x, y int) {
	return snake[len(snake)-1].X, snake[len(snake)-1].Y
}

func GetTailLocation() (x, y int) {
	return snake[0].X, snake[0].Y
}

func OnKeyPress() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			prevMove = move

			switch ev.Key {
			case termbox.KeyArrowUp:
				move = Up
			case termbox.KeyArrowDown:
				move = Down
			case termbox.KeyArrowLeft:
				move = Left
			case termbox.KeyArrowRight:
				move = Right
			case termbox.KeyEsc:
				gameOver <- true
			}
		}
	}
}

func AppendSnake(x, y int) {
	snake = append(snake, &Snake{
		X: x,
		Y: y,
	})
}

func RemoveSnake(x, y int) {
	for i, item := range snake {
		if x == item.X && y == item.Y {
			snake = append(snake[:i], snake[i+1:]...)
		}
	}

	field[x][y] = Road
}

func IsSnake(x, y int) bool {
	for _, item := range snake {
		if x == item.X && y == item.Y {
			return true
		}
	}

	return false
}

func GenerateItem() {
	x := rand.Intn(60)
	y := rand.Intn(30)

	field[x][y] = Item

	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	moveTicker := time.NewTicker(200 * time.Millisecond)

	_ = termbox.Init()
	defer termbox.Close()

	for y := 0; y < 30; y++ {
		for x := 0; x < 60; x++ {
			if x == 29 && y == 15 {
				AppendSnake(x, y)

				continue
			}

			if x == 30 && y == 15 {
				AppendSnake(x, y)

				continue
			}

			field[x][y] = Road
		}
	}

	GenerateItem()
	PrintField()

	go OnKeyPress()

	go func() {
		for {
			switch {
			case <-gameOver:
				termbox.Close()
				os.Exit(1)
			}
		}
	}()

	for {
		select {
		case <-moveTicker.C:
			headX, headY := GetHeadLocation()
			tailX, tailY := GetTailLocation()

			RemoveSnake(tailX, tailY)

			switch move {
			case Up:
				if headY-1 < 0 {
					gameOver <- true
				}

				if IsSnake(headX, headY-1) {
					move = prevMove
					AppendSnake(tailX, tailY)

					continue
				}

				if field[headX][headY-1] == Item {
					GenerateItem()
					AppendSnake(tailX, tailY)
				}

				AppendSnake(headX, headY-1)
			case Down:
				if headY+1 >= 30 {
					gameOver <- true
				}

				if IsSnake(headX, headY+1) {
					move = prevMove
					AppendSnake(tailX, tailY)

					continue
				}

				if field[headX][headY+1] == Item {
					GenerateItem()
					AppendSnake(tailX, tailY)
				}

				AppendSnake(headX, headY+1)
			case Left:
				if headX-1 < 0 {
					gameOver <- true
				}

				if IsSnake(headX-1, headY) {
					move = prevMove
					AppendSnake(tailX, tailY)

					continue
				}

				if field[headX-1][headY] == Item {
					GenerateItem()
					AppendSnake(tailX, tailY)
				}

				AppendSnake(headX-1, headY)
			case Right:
				if headX+1 >= 60 {
					gameOver <- true
				}

				if IsSnake(headX+1, headY) {
					move = prevMove
					AppendSnake(tailX, tailY)

					continue
				}

				if field[headX+1][headY] == Item {
					GenerateItem()
					AppendSnake(tailX, tailY)
				}

				AppendSnake(headX+1, headY)
			}

			PrintField()
		}
	}
}
