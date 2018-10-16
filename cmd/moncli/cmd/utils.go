package cmd

import "strconv"

func parseID(i string) (uint64, error) {
	return strconv.ParseUint(i, 10, 64)
}
