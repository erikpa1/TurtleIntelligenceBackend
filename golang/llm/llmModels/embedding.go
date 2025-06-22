package llmModels

import "math"

type Embedding [][]float32

func (self *Embedding) GetSimilarity(queryVector Embedding) float32 {

	highest := float32(0.0)

	for _, aVecRow := range *self {

		for _, bVecRow := range queryVector {
			similarity := cosineSimilarity(aVecRow, bVecRow)
			if similarity > highest {
				highest = similarity
			}
		}
	}

	return highest

}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
