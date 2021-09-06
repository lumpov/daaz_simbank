package utils

import (
	"math"
	"strconv"
	"strings"
)

// DecodeUCS2 string
func DecodeUCS2(str string) string {
	chanked := chunkSubstr(str, 4)

	builder := strings.Builder{}
	for i := range chanked {
		a, err := strconv.ParseInt(chanked[i], 16, 0)
		if err == nil {
			builder.WriteString(string(a))
		}
	}

	return builder.String()
}

func chunkSubstr(str string, size int) []string {
	numChunks := int(math.Ceil(float64(len(str) / size)))
	chunks := make([]string, numChunks)

	o := 0
	for i := 0; i < numChunks; i++ {
		chunks = append(chunks, str[o:o+size])

		o += size
	}

	return chunks
}
