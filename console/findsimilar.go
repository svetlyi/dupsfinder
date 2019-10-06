package console

//finds the most similar sample to needle from haystack
func findSimilar(needle string, haystack []string) string {
	var needleBytes []byte = []byte(needle)
	var sampleBytes []byte

	var bytesToLoop *[]byte
	var bytesToCompare *[]byte

	var similarBytesCounter int
	var maxSimilarBytes int
	var sampleIndexWithMaxBytes int

	for sampleIndex, sample := range haystack {
		similarBytesCounter = 0
		sampleBytes = []byte(sample)

		if len(sampleBytes) > len(needleBytes) {
			bytesToLoop = &needleBytes
			bytesToCompare = &sampleBytes
		} else {
			bytesToLoop = &sampleBytes
			bytesToCompare = &needleBytes
		}

		for _, currentByte := range *bytesToLoop {
			for _, currentByteToComapre := range *bytesToCompare {
				if currentByte == currentByteToComapre {
					similarBytesCounter++
				}
			}
		}
		if similarBytesCounter > maxSimilarBytes {
			maxSimilarBytes = similarBytesCounter
			sampleIndexWithMaxBytes = sampleIndex
		}
	}

	return haystack[sampleIndexWithMaxBytes]
}
