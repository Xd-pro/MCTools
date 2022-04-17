package main

import (
	"errors"
	"strings"
)

const WHITESPACES = " \t\n\r"

const ERROR_UNCLOSED_QUOTE = "unclosed quote"
const ERROR_EMPTY_INPUT = "empty input"

func HandleQuotes(in string) ([]string, []string, error) {

	if len(in) == 0 {
		return nil, nil, errors.New(ERROR_EMPTY_INPUT)
	}

	done := false
	char := 0
	lexed := []string{}
	var err error = nil
	for !done {
		if strings.ContainsRune(WHITESPACES, rune(in[char])) {
			char++
			if char >= len(in) {
				done = true
				break
			}
		} else if in[char] == '"' {
			text := ""
			char++
			if char >= len(in) {
				done = true
				break
			}
			for in[char] != '"' {
				text += string(in[char])
				char++
				if char >= len(in) {
					done = true
					err = errors.New(ERROR_UNCLOSED_QUOTE)
					break
				}
			}
			char++
			if char >= len(in) {
				done = true
			}
			lexed = append(lexed, text)
		} else {
			text := ""
			for !strings.ContainsRune(WHITESPACES, rune(in[char])) && in[char] != '"' {
				text += string(in[char])
				char++
				if char >= len(in) {
					done = true
					break
				}
			}
			lexed = append(lexed, text)
		}

	}

	out := make([]string, 0)
	flags := make([]string, 0)
	for _, value := range lexed {
		if strings.HasPrefix(value, "#") {
			flags = append(flags, strings.TrimPrefix(value, "#"))
		} else {
			out = append(out, value)
		}
	}

	return out, flags, err

}
