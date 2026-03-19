package model

import "math"

// NaiveBayesModel holds the trained parameters needed for prediction.
type NaiveBayesModel struct {
	PriorSpam float64
	PriorHam  float64

	PWordGivenSpam map[string]float64
	PWordGivenHam  map[string]float64
}

// Predict takes a bag-of-words count map x[w] for a new email
// and returns the predicted label (1 = spam, 0 = ham) and the raw scores.
func (m *NaiveBayesModel) Predict(x map[string]int) (label int, scoreSpam, scoreHam float64) {
	if m == nil {
		return 0, 0, 0
	}

	// Use logs to avoid underflow.
	scoreSpam = math.Log(m.PriorSpam)
	scoreHam = math.Log(m.PriorHam)

	for w, c := range x {
		if c == 0 {
			continue
		}
		pws, okS := m.PWordGivenSpam[w]
		if !okS {
			// If word not seen in training vocabulary, skip it;
			// probabilities for unknown words were not defined.
			continue
		}
		pwh, okH := m.PWordGivenHam[w]
		if !okH {
			continue
		}

		logPws := math.Log(pws)
		logPwh := math.Log(pwh)

		scoreSpam += float64(c) * logPws
		scoreHam += float64(c) * logPwh
	}

	if scoreSpam > scoreHam {
		return 1, scoreSpam, scoreHam
	}
	return 0, scoreSpam, scoreHam
}
