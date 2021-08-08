package inidata

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const formatCommentKey = "Comment Line %d"

var (
	errorIsComment = errors.New("is INI comment")
	errorIsSection = errors.New("is INI section")
)

type DataMap struct {
	data     map[string]map[string]string
	section  string
	sections []string
}

func (dm DataMap) addKV(k, v string) {
	if dm.data[dm.section] == nil {
		dm.data[dm.section] = make(map[string]string)
	}
	dm.data[dm.section][k] = v
}

// GetKey accepts a key name and optional section name and returns the associated value for it.
// The section defaults to GLOBAL if not included as the 2nd argument
func (dm DataMap) GetKey(args ...string) (v string, ok bool) {
	section := "GLOBAL"
	key := ""
	if len(args) > 0 {
		key = args[0]
		if len(args) > 1 {
			dm.section = args[1]
			section = args[1]
		}
		v, ok = dm.data[dm.section][key]
		if ok == false && section != "GLOBAL" {
			v, ok = dm.data["GLOBAL"][key]
		}
		return v, ok
	}

	return "", false
}

// GetSections returns a slice of all the sections
func (dm DataMap) GetSections() []string {
	count := len(dm.sections)
	if count == 0 || count < len(dm.data) {
		dm.sections = []string{}
		for s, _ := range dm.data {
			dm.sections = append(dm.sections, s)
		}
	}
	return dm.sections
}

// ParseBytes takes a byte slice and parses it into data
func (dm DataMap) ParseBytes(data []byte) error {
	line := ""
	lastC := '\n'
	i := 0
	dataLength := len(data)

	dm.data = make(map[string]map[string]string)
	dm.section = "GLOBAL"

	for i < dataLength {
		r, runeWidth := utf8.DecodeRune(data[i:])

		switch r {
		case '\n':
			if lastC != '\n' && len(line) > 0 {
				k, v, err := lineParse(line)
				if err != nil {
					switch err {
					case errorIsComment:
						k = fmt.Sprintf(formatCommentKey, i+1)

						// dm.data[dm.section][k] = v
						dm.addKV(k, v)
					case errorIsSection:
						dm.section = v
					default:
						return err
					}
				} else {
					// if dm.data[dm.section] == nil {
					// 	dm.data[dm.section] = make(map[string]string)
					// }
					// dm.data[dm.section][k] = v
					dm.addKV(k, v)
				}
			}
			line = ""
		default:
			line += string(r)
			if i+runeWidth == dataLength { // We've reached the end of the data and it doesn't end with a newline
				k, v, err := lineParse(line)
				if err != nil {
					switch err {
					case errorIsComment:
						k = fmt.Sprintf(formatCommentKey, i+1)
					case errorIsSection:
						dm.section = v
					default:
						return err
					}
				} else {
					if dm.data[dm.section] == nil {
						dm.data[dm.section] = make(map[string]string)
					}
					dm.data[dm.section][k] = v
				}
			}
		}
		i += runeWidth
		lastC = r
	}

	// for section, dict := range dm.data {
	// 	// log.Printf("inidata.DataMap.ParseBytes() | section: %q | dict: %q\n", section, dict)
	// 	for k, v := range dict {
	// 		log.Printf("inidata.DataMap.ParseBytes() | section: %q | %q: %q\n", section, k, v)
	// 	}
	// }

	return nil
}

// SetSection sets the active section name
func (dm DataMap) SetSection(s string) {
	dm.section = s
}

// NewDataMap returns a new DataMap
func NewDataMap() DataMap {
	return DataMap{
		data:    make(map[string]map[string]string),
		section: "GLOBAL",
	}
}

func lineParse(line string) (k, v string, err error) {
	line = strings.TrimSpace(line)
	// log.Printf("inidata.lineParse() | line: %q (%d)\n", line, len(line))

	if len(line) > 0 {
		switch line[0] {
		case '#', ';':
			// Comment
			comment := strings.TrimSpace(line[1:])
			// log.Printf("inidata.lineParse() | Comment: %q\n", comment)

			return "Comment", comment, errorIsComment
		case '[':
			// Section
			section := strings.Trim(line, `[]`)
			section = strings.TrimSpace(section)
			// log.Printf("inidata.lineParse() | Section: %q\n", section)

			return "Section", section, errorIsSection
		default:
			// Key/Value
			kv := strings.SplitN(line, "=", 2)
			if len(kv) > 0 {
				k = strings.TrimSpace(kv[0])
				if len(kv) > 1 {
					v = strings.TrimSpace(kv[1])
				}
			}
			v = strings.Trim(v, `"`)
			// log.Printf("inidata.lineParse() | Key: %q | Value: %q\n", k, v)
		}
	}

	return k, v, nil
}
