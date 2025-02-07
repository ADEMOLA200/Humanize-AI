package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
	"undetectable-ai/DSk/repositories"

	"github.com/neurosnap/sentences"
)

var (
	englishTraining = sentences.NewStorage()
	fillerSentences = loadFillers()
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// func RewriteText(text string) string {
//     paraphrased, err := repositories.ParaphraseText(text)
//     if err == nil && paraphrased != "" {
//         return paraphrased
//     }
//     return text
// }

// RewriteText uses the paraphrase, then applies light rewriting modifications.
func RewriteText(text string) string {
	// First, call the paraphrase endpoint.
	paraphrased, err := repositories.ParaphraseText(text)
	if err == nil && paraphrased != "" {
		text = paraphrased
	}

	sents := splitSentences(text)
	var transformed []string
	for _, sent := range sents {
		if rand.Float64() < 0.6 {
			sent = varySentenceStructure(sent)
		}

		if rand.Float64() < 0.3 {
			sent = replaceSynonyms(sent)
		}
		
		if rand.Float64() < 0.3 {
			sent = addNaturalNoise(sent)
		}
		transformed = append(transformed, sent)
	}

	if rand.Float64() < 0.3 {
		rand.Shuffle(len(transformed), func(i, j int) {
			transformed[i], transformed[j] = transformed[j], transformed[i]
		})
	}

	if rand.Float64() < 0.3 {
		transformed = append(transformed, getContextualFiller(transformed))
	}

	return strings.Join(transformed, " ")
}

func splitSentences(text string) []string {
	tokenizer := sentences.NewSentenceTokenizer(englishTraining)
	sents := tokenizer.Tokenize(text)
	var result []string
	for _, s := range sents {
		result = append(result, s.Text)
	}
	return result
}

func varySentenceStructure(sentence string) string {
	// Lower the probability for converting sentence voice.
	if rand.Float64() < 0.2 {
		return convertVoice(sentence)
	}
	return sentence
}

func convertVoice(sentence string) string {
	words := strings.Fields(sentence)
	// A basic active/passive conversion example.
	if len(words) > 2 && strings.Contains(sentence, " by ") {
		return strings.Replace(sentence, " by ", " ", 1)
	}
	if len(words) > 2 {
		return words[1] + " " + words[0] + " " + strings.Join(words[2:], " ")
	}
	return sentence
}

func getWordAffixes(word string) (string, string) {
	var prefix, suffix string
	for _, ch := range word {
		if !unicode.IsLetter(ch) {
			prefix += string(ch)
		} else {
			break
		}
	}
	for i := len(word) - 1; i >= 0; i-- {
		if !unicode.IsLetter(rune(word[i])) {
			suffix = string(word[i]) + suffix
		} else {
			break
		}
	}
	return prefix, suffix
}

func replaceSynonyms(sentence string) string {
	words := strings.Fields(sentence)
	for i, word := range words {
		cleanWord := strings.Trim(word, ".,!?;:\"'")
		// Skip if too short or common.
		if len(cleanWord) < 5 || isCommonWord(cleanWord) {
			continue
		}
		// Lower replacement probability.
		if rand.Float64() < 0.1 {
			synonym, err := repositories.GetSynonym(cleanWord)
			if err == nil && synonym != "" {
				prefix, suffix := getWordAffixes(word)
				if unicode.IsUpper(rune(word[0])) {
					synonym = strings.Title(synonym)
				}
				words[i] = prefix + synonym + suffix
			}
		}
	}
	return strings.Join(words, " ")
}

func addNaturalNoise(sentence string) string {
	replacements := map[string]string{
		"because":   "'cause",
		"however":   "though",
		"therefore": "so",
		" students": " learners",
		" utilize":  " use",
	}
	pauseFormats := []string{
		" -- ", ", you know, ", " ... ", ", well, ",
	}
	for k, v := range replacements {
		if rand.Float64() < 0.1 && strings.Contains(sentence, k) {
			sentence = strings.Replace(sentence, k, v, 1)
		}
	}
	if rand.Float64() < 0.1 && len(sentence) > 40 {
		pause := pauseFormats[rand.Intn(len(pauseFormats))]
		insertPos := rand.Intn(len(sentence)-20) + 10
		sentence = sentence[:insertPos] + pause + sentence[insertPos:]
	}
	return sentence
}

func getContextualFiller(sentences []string) string {
	if len(sentences) == 0 {
		return ""
	}
	lastSentence := sentences[len(sentences)-1]
	keywords := extractKeywords(lastSentence)
	templates := []string{
		"Now, this makes you wonder about %s...",
		"It's clear that %s plays a big role here...",
		"Doesn't this make you think about %s?",
		"This really highlights the importance of %s...",
	}
	if len(keywords) > 0 {
		return fmt.Sprintf(templates[rand.Intn(len(templates))],
			strings.Join(keywords[:min(2, len(keywords))], " and "))
	}
	return fillerSentences[rand.Intn(len(fillerSentences))]
}

func extractKeywords(sentence string) []string {
	var keywords []string
	words := strings.Fields(sentence)
	for _, word := range words {
		if len(word) > 5 && !isCommonWord(word) {
			keywords = append(keywords, strings.Trim(word, ".,!?"))
		}
	}
	return keywords
}

func isCommonWord(word string) bool {
	common := map[string]bool{
		"the": true, "and": true, "that": true, "this": true, "with": true,
	}
	return common[strings.ToLower(word)]
}

func loadFillers() []string {
	return []string{
		"From this perspective, it becomes evident that",
		"Considering these factors holistically",
		"Natural language processing reveals",
		"Contemporary analysis suggests",
		"Modern interpretations emphasize",
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
