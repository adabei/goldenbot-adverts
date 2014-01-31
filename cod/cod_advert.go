// Package advert provides a plugin to display messages in a set interval
// to all players.
package cod

import (
	"bufio"
	"fmt"
	"github.com/adabei/goldenbot/rcon"
	"log"
	"os"
	"time"
)

type Advert struct {
	input    string
	Interval int
	requests chan rcon.RCONQuery
}

func NewAdvert(input string, interval int, requests chan rcon.RCONQuery) *Advert {
	a := new(Advert)
	a.input = input
	a.requests = requests
	a.Interval = interval
	return a
}

func (a *Advert) Setup() error {
	return nil
}

func (a *Advert) Start() {
	ads, err := read(a.input)
	if err != nil {
		log.Fatal("Failed to load adverts from file ", a.input, ": ", err)
	}

	for {
		for _, ad := range ads {
			if ad != "" {
				// TODO missing say prefix
				log.Println("adverts: sending", ad, "to RCON")
				a.requests <- rcon.RCONQuery{Command: fmt.Sprint("say \"", ad, "\""), Response: nil}
			} else {
				time.Sleep(time.Duration(a.Interval) * time.Millisecond)
			}
		}
	}
}

func read(from string) ([]string, error) {
	ads := make([]string, 0)
	fi, err := os.Open(from)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	scanner := bufio.NewScanner(fi)

	for scanner.Scan() {
		if val := scanner.Text(); len(val) == 0 {
			ads = append(ads, "")
		} else {
			ads = append(ads, val)
		}
	}

	return append(ads, ""), nil
}
