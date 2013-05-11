package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogReaderBasic(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob(filepath.Join("testdata", "L*.log"))
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range files {
		base := filepath.Base(name)

		f, err := os.Open(name)
		if err != nil {
			t.Errorf("Opening %s: %v", base, err)
			continue
		}

		read := LogReader(f)

		prevLine := 0

		for {
			line, parsed, err := read()

			if line != prevLine+1 {
				t.Errorf("Line %d of %s came after line %d", line, base, prevLine)
			}
			prevLine = line

			if err != nil {
				if err != io.EOF {
					t.Errorf("Reading %s+%d: %v", base, line, err)
				}

				break
			}

			//t.Logf("%s+%d: %v", base, line, parsed)

			if parsed["Line"] != line {
				t.Errorf("Line %d of %s is missing a Line field. (%+v)", line, name, parsed["Line"])
			}

			if _, ok := parsed["Time"].(time.Time); !ok {
				t.Errorf("Line %d of %s is missing a Time field. (%+v)", line, name, parsed["Time"])
			}

			if _, ok := parsed["Text"].(string); !ok {
				t.Errorf("Line %d of %s is missing a Text field. (%+v)", line, name, parsed["Text"])
			}
		}

		if err = f.Close(); err != nil {
			t.Errorf("Closing %s: %v", base, err)
		}
	}
}

func TestLogReaderParse(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob(filepath.Join("testdata", "L*.log"))
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range files {
		base := filepath.Base(name)

		f, err := os.Open(name)
		if err != nil {
			t.Errorf("Opening %s: %v", base, err)
			continue
		}

		read := LogReader(f)

		for {
			line, parsed, err := read()
			if err != nil {
				if err != io.EOF {
					t.Errorf("Reading %s+%d: %v", base, line, err)
				}

				break
			}

			//t.Logf("%s+%d: %v", base, line, parsed)

			text, ok := parsed["Text"].(string)
			if !ok {
				t.Errorf("Text field missing on one or more lines of %s", base)
				break
			}

			if strings.Contains(text, "(DEATH)") {
				if strings.Contains(text, " used pills on ") {
					// bug in L4D2. This line is most likely corrupt.
					continue
				}

				if _, unparsed := parsed["Unparsed"]; unparsed {
					t.Errorf("%s+%d: (DEATH) message unparsed", base, line)
					continue
				}

				// TODO: more checking
			} else if strings.Contains(text, " spawned as a ") {
				if _, unparsed := parsed["Unparsed"]; unparsed {
					t.Errorf("%s+%d: \"spawned as a\" message unparsed", base, line)
					continue
				}

				// TODO: more checking
			} else if strings.Contains(text, " joined team ") {
				if _, unparsed := parsed["Unparsed"]; unparsed {
					t.Errorf("%s+%d: \"joined team\" message unparsed", base, line)
					continue
				}

				// TODO: more checking
			} else if strings.Contains(text, "Respawning ") {
				if _, unparsed := parsed["Unparsed"]; unparsed {
					t.Errorf("%s+%d: \"Respawning\" message unparsed", base, line)
					continue
				}

				// TODO: more checking
			} else if _, unparsed := parsed["Unparsed"]; !unparsed {
				t.Errorf("%s+%d: Unexpected parsed message: %v", base, line, parsed)
				continue
			}

		}

		if err = f.Close(); err != nil {
			t.Errorf("Closing %s: %v", base, err)
		}
	}
}
