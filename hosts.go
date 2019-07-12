package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const expandFmt = `([[:word:]-]+)(\{([\d]+)..([\d]+)\})?([[:word:]\.-]*)`

var expandRe = regexp.MustCompile(expandFmt)

// hosts takes the raw string input from --hosts and resolves
// the actual list of hosts that commando will execute against.
func hosts(input string) []string {
	split := strings.Split(input, ",")
	return resolve(split)
}

func resolve(resolvable []string) []string {
	var resolved []string
	for _, raw := range resolvable {
		resolved = append(resolved, expand(raw)...)
	}
	return resolved
}

func expand(raw string) []string {
	var expanded []string
	raw = strings.TrimSpace(raw)

	matches := expandRe.FindAllStringSubmatch(raw, -1)

	if !strings.Contains(matches[0][0], "..") {
		expanded = append(expanded, matches[0][0])
	} else {
		root := matches[0][1]
		low, err := strconv.Atoi(matches[0][3])
		if err != nil {
			return nil
		}
		high, err := strconv.Atoi(matches[0][4])
		if err != nil {
			return nil
		}
		tail := "" // trailing domain
		if len(matches[0]) == 6 {
			tail = matches[0][5]
		}
		for n := low; n <= high; n++ {
			host := fmt.Sprintf("%s%d%s", root, n, tail)
			expanded = append(expanded, host)
		}
	}

	return expanded
}
