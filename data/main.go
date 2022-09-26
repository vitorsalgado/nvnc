package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vitorsalgado/nvnc/internal/httpx"
	"github.com/vitorsalgado/nvnc/internal/tse"
	"golang.org/x/text/encoding/charmap"
)

const (
	_fCandState      = 12
	_fCandType       = 15
	_fCandNumber     = 17
	_fCandName       = 18
	_fCandBallotName = 19
	_fCandSituation  = 24
	_fCandParty      = 29
	_fCandPartyName  = 30
)

type Cand struct {
	Name         string `json:"name"`
	BallotNumber string `json:"ballot_number"`
	BallotName   string `json:"ballot_name"`
	Type         string `json:"type"`
	Situation    string `json:"situation"`
	Party        string `json:"party"`
	PartyName    string `json:"party_name"`
	State        string `json:"state"`
}

func main() {
	_ = godotenv.Load()

	// cwd, _ := os.Getwd()
	httpclient := httpx.Client(&httpx.Conf{})

	// configs
	candurl := os.Getenv("DATA_URL")
	filename := os.Getenv("DATA_FILENAME")
	dest := os.Getenv("DATA_DEST")

	// fetching candidates csv
	b, err := tse.FetchCandidates(httpclient, candurl, filename)
	if err != nil {
		log.Fatal(fmt.Errorf("error fetching candidates csv from TSE. reason=%w", err))
	}

	c := 0
	r := charmap.ISO8859_1.NewDecoder().Reader(bytes.NewReader(b))
	s := bufio.NewScanner(r)
	f, err := os.Create(dest)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating destination file=%s. reason=%w", dest, err))
	}

	defer f.Close()

	f.WriteString("[")

	for s.Scan() {
		c++

		if c == 1 {
			continue
		}

		if c > 2 {
			f.WriteString(",")
		}

		entries := strings.Split(s.Text(), ";")
		cand := Cand{
			Name:         strings.Replace(entries[_fCandName-1], "\"", "", -1),
			BallotNumber: strings.Replace(entries[_fCandNumber-1], "\"", "", -1),
			BallotName:   strings.Replace(entries[_fCandBallotName-1], "\"", "", -1),
			Type:         strings.Replace(entries[_fCandType-1], "\"", "", -1),
			Situation:    strings.Replace(entries[_fCandSituation-1], "\"", "", -1),
			Party:        strings.Replace(entries[_fCandParty-1], "\"", "", -1),
			PartyName:    strings.Replace(entries[_fCandPartyName-1], "\"", "", -1),
			State:        strings.Replace(entries[_fCandState-1], "\"", "", -1),
		}

		jzon, _ := json.Marshal(cand)

		f.WriteString(string(jzon))
	}

	f.WriteString("]")
}
