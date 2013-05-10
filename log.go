package main

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"time"
)

const TimeFormat = "L 01/02/2006 - 15:04:05: "

var (
	pattern_joined_team  = regexp.MustCompile(`^"([^<"]+)<(\d*)><(BOT|STEAM_\d:\d:\d+)><([A-Z][a-z]*)>" joined team "([A-Z][a-z]*)"$`)
	pattern_spawned_as_a = regexp.MustCompile(`^([^<"]+) spawned as a ([A-Z][A-Za-z]*)$`)
	pattern_respawning   = regexp.MustCompile(`^Respawning ([^<"]+)$`)
	pattern_death        = regexp.MustCompile(`^\(DEATH\)"([^<"]+)<(\d*)><(BOT|STEAM_\d:\d:\d+|)><([A-Za-z]*)><([A-Za-z]*)><([A-Z]+)><(-?\d+(?:\+\d+)?)><setpos(?:_exact)? (-?\d+\.\d\d) (-?\d+\.\d\d) (-?\d+\.\d\d); setang (-?\d+\.\d\d) (-?\d+\.\d\d) (-?\d+\.\d\d)><Area (\d+)>" killed "([^<"]+)<(\d*)><(BOT|STEAM_\d:\d:\d+|)><([A-Za-z]*)><([A-Za-z]*)><([A-Z]+)><(-?\d+(?:\+\d+)?)><setpos(?:_exact)? (-?\d+\.\d\d) (-?\d+\.\d\d) (-?\d+\.\d\d); setang (-?\d+\.\d\d) (-?\d+\.\d\d) (-?\d+\.\d\d)><Area (\d+)>" with "([a-z0-9_]+)"( \(headshot\))?$`)
)

func LogReader(r io.Reader) func() (int, map[string]interface{}, error) {
	br := bufio.NewReaderSize(r, 0x10000)
	ln := 0
	return func() (int, map[string]interface{}, error) {
		ln++
		line, err := br.ReadString('\n')
		if err != nil {
			return ln, nil, err
		}
		line = line[:len(line)-1] // remove trailing newline
		parsed := map[string]interface{}{"Text": line, "Line": ln}

		t, err := time.Parse(TimeFormat, line[:len(TimeFormat)])
		if err != nil {
			return ln, nil, err
		}
		parsed["Time"] = t

		line = line[len(TimeFormat):]

		if m := pattern_joined_team.FindStringSubmatch(line); m != nil {
			parsed["Type"] = "joined team"
			parsed["Name"] = m[1]
			parsed["UID"], err = strconv.Atoi(m[2])
			if err != nil {
				return ln, nil, err
			}
			parsed["Addr"] = m[3]
			parsed["OldTeam"] = m[4]
			parsed["NewTeam"] = m[5]
		} else if m := pattern_spawned_as_a.FindStringSubmatch(line); m != nil {
			parsed["Type"] = "spawned as a"
			parsed["Name"] = m[1]
			parsed["Class"] = m[2]
		} else if m := pattern_respawning.FindStringSubmatch(line); m != nil {
			parsed["Type"] = "respawning"
			parsed["Name"] = m[1]
		} else if m := pattern_death.FindStringSubmatch(line); m != nil {
			parsed["Type"] = "killed"
			attacker := make(map[string]interface{})
			parsed["Attacker"] = attacker
			attacker["Name"] = m[1]
			attacker["UID"], err = strconv.Atoi(m[2])
			if err != nil {
				return ln, nil, err
			}
			attacker["Addr"] = m[3]
			attacker["Team"] = m[4]
			attacker["Class"] = m[5]
			attacker["Status"] = m[6]
			attacker["Health"] = m[7]
			attacker_pos := make([]float64, 6)
			attacker["Pos"] = attacker_pos
			attacker_pos[0], err = strconv.ParseFloat(m[8], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker_pos[1], err = strconv.ParseFloat(m[9], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker_pos[2], err = strconv.ParseFloat(m[10], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker_pos[3], err = strconv.ParseFloat(m[11], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker_pos[4], err = strconv.ParseFloat(m[12], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker_pos[5], err = strconv.ParseFloat(m[13], 64)
			if err != nil {
				return ln, nil, err
			}
			attacker["Area"], err = strconv.Atoi(m[14])
			if err != nil {
				return ln, nil, err
			}

			victim := make(map[string]interface{})
			parsed["Victim"] = victim
			victim["Name"] = m[15]
			victim["UID"], err = strconv.Atoi(m[16])
			if err != nil {
				return ln, nil, err
			}
			victim["Addr"] = m[17]
			victim["Team"] = m[18]
			victim["Class"] = m[19]
			victim["Status"] = m[20]
			victim["Health"] = m[21]
			victim_pos := make([]float64, 6)
			victim["Pos"] = attacker_pos
			victim_pos[0], err = strconv.ParseFloat(m[22], 64)
			if err != nil {
				return ln, nil, err
			}
			victim_pos[1], err = strconv.ParseFloat(m[23], 64)
			if err != nil {
				return ln, nil, err
			}
			victim_pos[2], err = strconv.ParseFloat(m[24], 64)
			if err != nil {
				return ln, nil, err
			}
			victim_pos[3], err = strconv.ParseFloat(m[25], 64)
			if err != nil {
				return ln, nil, err
			}
			victim_pos[4], err = strconv.ParseFloat(m[26], 64)
			if err != nil {
				return ln, nil, err
			}
			victim_pos[5], err = strconv.ParseFloat(m[27], 64)
			if err != nil {
				return ln, nil, err
			}
			victim["Area"], err = strconv.Atoi(m[28])
			if err != nil {
				return ln, nil, err
			}

			parsed["Weapon"] = m[29]
			parsed["Headshot"] = m[30] != ""
		} else {
			parsed["Unparsed"] = true
		}

		return ln, parsed, nil
	}
}
