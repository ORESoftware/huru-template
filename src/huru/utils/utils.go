package utils

import "bytes"

// JoinArgs joins strings
func JoinArgs(strangs ...string) string {
	buffer := bytes.NewBufferString("")
	for _, s := range strangs {
		buffer.WriteString(s)
	}
	return buffer.String()
}
