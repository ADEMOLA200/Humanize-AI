package repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kljensen/snowball"
)

type SparqlResponse struct {
	Results struct {
		Bindings []struct {
			Synonym struct {
				Value string `json:"value"`
			} `json:"synonym"`
		} `json:"bindings"`
	} `json:"results"`
}

func GetSynonym(word string) (string, error) {
	// SPARQL query to fetch synonyms
	query := fmt.Sprintf(`SELECT ?synonym WHERE {
		?s <http://wordnet-rdf.princeton.edu/ontology#word> "%s" .
		?s <http://wordnet-rdf.princeton.edu/ontology#synonym> ?synonym .
	}`, word)

	// Encode the query
	encodedQuery := url.QueryEscape(query)
	sparqlURL := fmt.Sprintf("http://ldf.fi/wordnet/sparql?query=%s", encodedQuery)

	// Make the request
	resp, err := http.Get(sparqlURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var result SparqlResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Extract synonyms
	if len(result.Results.Bindings) == 0 {
		return "", fmt.Errorf("no synonyms found for the word")
	}

	wordStem, _ := snowball.Stem(word, "english", true)

	validSynonyms := []string{}
	for _, binding := range result.Results.Bindings {
		synonym := binding.Synonym.Value
		if isValidSynonym(synonym, wordStem) {
			validSynonyms = append(validSynonyms, synonym)
		}
	}

	if len(validSynonyms) == 0 {
		return "", fmt.Errorf("no valid synonyms found for the context")
	}

	return validSynonyms[0], nil
}

func isValidSynonym(synonym, wordStem string) bool {
	synonymStem, _ := snowball.Stem(synonym, "english", true)
	return wordStem == synonymStem
}
