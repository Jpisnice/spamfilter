package bagofwords

import (
	tokenizer "github.com/Jpisnice/spamfilter/tokenizer"
	"strings"
)


type BagOfWords struct {
	label int
	frequency map[string]int
}
func (b *BagOfWords) Init(TokenizedRecord tokenizer.TokenizedRecord){
	b.frequency = calculateFrequency(TokenizedRecord.Words)
	b.label = TokenizedRecord.Label

}
 func calculateFrequency(words []string) map[string]int {
	frequency := make(map[string]int)
	for _, word := range words {
		frequency[strings.ToUpper(word)]++
	}
	return frequency
 }

 func (b *BagOfWords) GetFrequency(word string) int {
	return b.frequency[strings.ToUpper(word)]
 }
 func (b *BagOfWords) GetLabel() int {
	return b.label
 }

 func (b *BagOfWords) GetFrequencies() map[string]int {
	return b.frequency
 }