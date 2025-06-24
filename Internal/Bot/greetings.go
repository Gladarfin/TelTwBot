package bot

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strings"
)

type Greeting struct {
	Language string
	Text     string
}

type Greeter struct {
	greetings []Greeting
	rand      *rand.Rand
}

func NewGreeter(filename string, rnd *rand.Rand) (*Greeter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var greetings []Greeting

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			greetings = append(greetings, Greeting{
				Language: strings.TrimSpace(parts[0]),
				Text:     strings.TrimSpace(parts[1]),
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Greeter{
		greetings: greetings,
		rand:      rnd,
	}, nil
}

func (g *Greeter) GetRandomGreeting() Greeting {
	if len(g.greetings) == 0 {
		return Greeting{}
	}

	return g.greetings[g.rand.Intn(len(g.greetings))]
}

func (g *Greeter) Count() int {
	return len(g.greetings)
}

func (g *Greeter) Validate() error {
	if len(g.greetings) == 0 {
		return errors.New("no greetings loaded")
	}
	return nil
}
