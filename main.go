package main

import (
	"sync"

	"github.com/Jpisnice/spamfilter/logger"
	bagofwords "github.com/Jpisnice/spamfilter/bagofwords"
	"github.com/Jpisnice/spamfilter/model"
	tokenizer "github.com/Jpisnice/spamfilter/tokenizer"
	parser "github.com/Jpisnice/spamfilter/tokenizer/parser"
	"strings"
)


func main() {
	logger.Logger.Print("Starting tokenization")
	records, count := parser.Parse("enron-spam/Enron_Txt_fn.csv")
	logger.Logger.Print("Parsed records:", count)

	const sections = 8
	if len(records) == 0 {
		logger.Logger.Print("No records to tokenize")
		return
	}

	chunkSize := (len(records) + sections - 1) / sections
	tokenized := make([]tokenizer.TokenizedRecord, len(records))

	var wg sync.WaitGroup
	for section := 0; section < sections; section++ {
		start := section * chunkSize
		if start >= len(records) {
			break
		}

		end := start + chunkSize
		if end > len(records) {
			end = len(records)
		}

		wg.Add(1)
		go func(section, start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				tokenized[i] = tokenizer.Tokenize(records[i])
			}
			logger.Logger.Printf("Section %d tokenized records [%d:%d)", section+1, start, end)
		}(section, start, end)
	}

	wg.Wait()
	logger.Logger.Print("All tokenization workers completed")
	logger.Logger.Print("Tokenized records:", len(tokenized))

	if len(tokenized) == 0 {
		logger.Logger.Print("No records to tokenize")
		return
	}

	chunkSize2 := (len(tokenized) + sections - 1) / sections
	bagRecords := make([]bagofwords.BagOfWords, len(tokenized))

	var wg2 sync.WaitGroup
	for section := 0; section < sections; section++ {
		start := section * chunkSize2
		if start >= len(tokenized) {
			break
		}

		end := start + chunkSize2
		if end > len(tokenized) {
			end = len(tokenized)
		}

		wg2.Add(1)
		go func(section, start, end int) {
			defer wg2.Done()
			for i := start; i < end; i++ {
				bagRecords[i].Init(tokenized[i])
			}
			logger.Logger.Printf("Section %d bagofwords records [%d:%d)", section+1, start, end)
		}(section, start, end)
	}

	wg2.Wait()

	logger.Logger.Print("All bagofwords workers completed")
	logger.Logger.Print("Bagofwords:", len(bagRecords))

	if len(bagRecords) > 0 {
		bagRecords[0].GetFrequency("hello")
	}

	spamBagOfWords := make([]bagofwords.BagOfWords, 0, len(bagRecords))
	hamBagOfWords := make([]bagofwords.BagOfWords, 0, len(bagRecords))
	for _, bag := range bagRecords {
		if bag.GetLabel() == 1 {
			spamBagOfWords = append(spamBagOfWords, bag)
		}
		if bag.GetLabel() == 0 {
			hamBagOfWords = append(hamBagOfWords, bag)
		}
	}
	totalDocs := len(spamBagOfWords) + len(hamBagOfWords)
	logger.Logger.Print("Spam BagOfWords:", len(spamBagOfWords))
	logger.Logger.Print("Ham BagOfWords:", len(hamBagOfWords))
	logger.Logger.Print("Total BagOfWords:", totalDocs)
	// Aggregate total word counts per class
	spamWordCount := make(map[string]int)
	hamWordCount := make(map[string]int)

	for _, bag := range spamBagOfWords {
		for word, count := range bag.GetFrequencies() {
			spamWordCount[word] += count
		}
	}

	for _, bag := range hamBagOfWords {
		for word, count := range bag.GetFrequencies() {
			hamWordCount[word] += count
		}
	}

	vocabulary := make(map[string]struct{})
	totalSpamTokens := 0
	totalHamTokens := 0

	for word, count := range spamWordCount {
		vocabulary[word] = struct{}{}
		totalSpamTokens += count
	}
	for word, count := range hamWordCount {
		vocabulary[word] = struct{}{}
		totalHamTokens += count
	}

	alpha := 1.0
	vocabSize := float64(len(vocabulary))
	spamDenominator := float64(totalSpamTokens) + alpha*vocabSize
	hamDenominator := float64(totalHamTokens) + alpha*vocabSize

	// Laplace-smoothed likelihoods for every word in the vocabulary.
	pWordGivenSpam := make(map[string]float64, len(vocabulary))
	pWordGivenHam := make(map[string]float64, len(vocabulary))
	for word := range vocabulary {
		pWordGivenSpam[word] = (float64(spamWordCount[word]) + alpha) / spamDenominator
		pWordGivenHam[word] = (float64(hamWordCount[word]) + alpha) / hamDenominator
	}

	logger.Logger.Print("Vocabulary size |V|:", len(vocabulary))
	logger.Logger.Print("Total spam tokens:", totalSpamTokens)
	logger.Logger.Print("Total ham tokens:", totalHamTokens)
	logger.Logger.Print("Computed P(w|spam) entries:", len(pWordGivenSpam))
	logger.Logger.Print("Computed P(w|ham) entries:", len(pWordGivenHam))

	// Class priors P(spam) and P(ham).
	pSpam := float64(len(spamBagOfWords)) / float64(totalDocs)
	pHam := float64(len(hamBagOfWords)) / float64(totalDocs)
	logger.Logger.Print("P(spam):", pSpam)
	logger.Logger.Print("P(ham):", pHam)

	// Build Naive Bayes model.
	nb := &model.NaiveBayesModel{
		PriorSpam:      pSpam,
		PriorHam:       pHam,
		PWordGivenSpam: pWordGivenSpam,
		PWordGivenHam:  pWordGivenHam,
	}

	// Example: classify a single existing email (e.g., first tokenized record) as a sanity check.
		exampleFreq := make(map[string]int)
		for _, w := range strings.Fields(strings.ToUpper("hello world im here to advertise our new product that enhances productivity")) {
			exampleFreq[w]++
		}
		label, scoreSpam, scoreHam := nb.Predict(exampleFreq)
		logger.Logger.Print("Example prediction label (1=spam,0=ham):", label)
		logger.Logger.Print("Example scoreSpam:", scoreSpam)
		logger.Logger.Print("Example scoreHam:", scoreHam)
	
}