package tokenizer

import (
	parser "github.com/Jpisnice/spamfilter/tokenizer/parser"
	"strings"
)
type TokenizedRecord struct {
	Words []string
	Label int
}
func Tokenize(record parser.Record) TokenizedRecord {
	words := strings.Fields(record.Message)
	return TokenizedRecord{Words: words, Label: record.Label}
}
