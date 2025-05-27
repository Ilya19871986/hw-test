package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	Word  string
	Count int
}

// создать мапу word -> count.
func countWordFrequencies(text string) map[string]int {
	frequencyMap := make(map[string]int)
	words := strings.Fields(text)
	for _, word := range words {
		frequencyMap[word]++
	}
	return frequencyMap
}

// отсортировать мапу word -> count по count && word и вернуть word []string.
func sortWordFrequencies(frequencyMap map[string]int) []string {
	frequency := make([]wordCount, 0, len(frequencyMap))
	for word, count := range frequencyMap {
		frequency = append(frequency, wordCount{Word: word, Count: count})
	}
	sort.Slice(frequency, func(i, j int) bool {
		if frequency[i].Count != frequency[j].Count {
			return frequency[i].Count > frequency[j].Count
		}
		return frequency[i].Word < frequency[j].Word
	})
	result := make([]string, 0, len(frequency))
	for _, word := range frequency {
		result = append(result, word.Word)
	}
	return result
}

func Top10(text string) []string {
	frequencyMap := countWordFrequencies(text)
	sortedWords := sortWordFrequencies(frequencyMap)
	if len(sortedWords) > 10 {
		return sortedWords[:10]
	}
	return sortedWords
}
