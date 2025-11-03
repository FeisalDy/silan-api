package miscellaneous

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
	RTL  int    `json:"rtl,omitempty"`
}

var languages []Language

func init() {
	loadLanguages()
}

func loadLanguages() {
	// Get the path to list.json
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	jsonPath := filepath.Join(exePath, "pkg/miscellaneous/list.json")

	// Try loading from executable directory first
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		// Fallback to relative path from current working directory
		data, err = os.ReadFile("pkg/miscellaneous/list.json")
		if err != nil {
			panic("Could not load language list: " + err.Error())
		}
	}

	err = json.Unmarshal(data, &languages)
	if err != nil {
		panic("Could not parse language list: " + err.Error())
	}
}

// GetAllLanguages returns all available languages
func GetAllLanguages() []Language {
	return languages
}

// SearchLanguages searches both name and native fields
func SearchLanguages(query string) []Language {
	if query == "" {
		return languages
	}

	query = strings.ToLower(query)
	var results []Language
	seen := make(map[string]bool)

	for _, lang := range languages {
		if strings.Contains(strings.ToLower(lang.Name), query) {
			if !seen[lang.Code] {
				results = append(results, lang)
				seen[lang.Code] = true
			}
		}
	}

	return results
}

// GetLanguageByCode returns a language by its ISO 639-1 code
func GetLanguageByCode(code string) []Language {
	code = strings.ToLower(code)
	for _, lang := range languages {
		if strings.ToLower(lang.Code) == code {
			return []Language{lang}
		}
	}
	return nil
}
