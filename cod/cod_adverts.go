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

type Adverts struct {
	input    string
	Interval int
	cfg      Config
	requests chan rcon.RCONQuery
}

type Config struct {
	Prefix   string
	Input    string
	Interval int
}

func NewAdverts(cfg Config, requests chan rcon.RCONQuery) *Adverts {
	a := new(Adverts)
	a.cfg = cfg
	a.requests = requests
	return a
}

func (a *Adverts) Setup() error {
	return nil
}

func (a *Adverts) Start() {
	ads, err := read(a.cfg.Input)
	if err != nil {
		log.Fatal("Failed to load adverts from file ", a.cfg.Input, ": ", err)
	}

	for {
		for _, ad := range ads {
			if ad != "" {
				// TODO missing say prefix
				log.Println("adverts: sending", ad, "to RCON")
				a.requests <- rcon.RCONQuery{Command: fmt.Sprint("say \"", a.cfg.Prefix, ad, "\""), Response: nil}
			} else {
				time.Sleep(time.Duration(a.cfg.Interval) * time.Millisecond)
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
