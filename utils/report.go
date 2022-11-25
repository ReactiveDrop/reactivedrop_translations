package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"git.lubar.me/ben/valve/translation"
	"golang.org/x/text/language/display"
)

func eachFile(patterns []string, addSuffix string, cb func(string, io.Reader)) {
	for _, pattern := range patterns {
		names, err := filepath.Glob(pattern + addSuffix)
		if err != nil {
			panic(err)
		}

		for _, name := range names {
			f, err := os.Open(name)
			if err != nil {
				panic(err)
			}

			cleanName := filepath.ToSlash(name)[3:] // strip the "../"

			cb(cleanName, f)

			err = f.Close()
			if err != nil {
				panic(err)
			}
		}
	}
}

func countIndentedLines(r io.Reader) int {
	// don't care about performance - this runs offline
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	state, count := 0, 0
	for _, ch := range b {
		switch ch {
		case '\t':
			if state == 0 {
				state = 1
			}
		case '/':
			if state == 1 {
				state = 2
			}
		case '\n':
			state = 0
		default:
			if state == 1 {
				count++
			}

			state = 2
		}
	}

	return count
}

type fileForReport struct {
	prefix   string
	suffix   string
	indented map[string]int
}

func sortFilesForReport(m map[[2]string]map[string]int) []fileForReport {
	s, i := make([]fileForReport, len(m)), 0
	for k, v := range m {
		s[i].prefix = k[0]
		s[i].suffix = k[1]
		s[i].indented = v
		i++
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].prefix < s[j].prefix
	})

	return s
}

func generateReport() {
	languageFiles := make(map[[2]string]map[string]int)
	processLanguageFile := func(s string, r io.Reader) {
		indented := countIndentedLines(r)
		if indented&1 == 1 {
			panic("unexpected odd number of indented lines in " + s)
		}

		i, j := strings.LastIndexByte(s, '_'), strings.LastIndexByte(s, '.')
		prefix, lang, suffix := s[:i+1], s[i+1:j], s[j:]

		if lang == sourceLanguage && indented != 0 {
			panic("unexpected indented lines in source file " + s)
		}

		m, ok := languageFiles[[2]string{prefix, suffix}]
		if !ok {
			m = make(map[string]int)
			languageFiles[[2]string{prefix, suffix}] = m
		}

		m[lang] = indented / 2
	}
	eachFile(txtLanguageFiles[:], "_*.txt", processLanguageFile)
	eachFile(vdfLanguageFiles[:], "_*.vdf", processLanguageFile)
	eachFile(txtAddonLanguageFiles[:], "_*.txt", processLanguageFile)
	sortedLanguageFiles := sortFilesForReport(languageFiles)

	var checked [][]fileForReport

	for _, check := range checkButNoSync {
		checkCategory := make(map[[2]string]map[string]int)
		eachFile(check.patterns, "", func(s string, r io.Reader) {
			indented := countIndentedLines(r)

			i, j := strings.LastIndexByte(s, '_'), strings.LastIndexByte(s, '.')
			prefix, lang, suffix := s[:i+1], s[i+1:j], s[j:]

			if lang == sourceLanguage && indented != 0 {
				panic("unexpected indented lines in source file " + s)
			}

			if _, ok := languageFiles[[2]string{prefix, suffix}]; ok {
				return // don't count language files twice
			}

			m, ok := checkCategory[[2]string{prefix, suffix}]
			if !ok {
				m = make(map[string]int)
				checkCategory[[2]string{prefix, suffix}] = m
			}

			m[lang] = indented
		})

		checked = append(checked, sortFilesForReport(checkCategory))
	}

	var inventoryItems []map[string]interface{}
	eachFile(checkInventorySchema[:], "", func(s string, r io.Reader) {
		dec := json.NewDecoder(r)
		dec.UseNumber()

		var data struct {
			AppID int                      `json:"appid"`
			Items []map[string]interface{} `json:"items"`
		}

		err := dec.Decode(&data)
		if err != nil {
			panic(err)
		}

		name := s
		for _, item := range data.Items {
			item["_file"] = name
		}

		inventoryItems = append(inventoryItems, data.Items...)
	})

	fmt.Print("| Language | Strings |")
	for _, check := range checkButNoSync {
		fmt.Printf(" %s |", check.short)
	}
	fmt.Print(" Inventory |\n| --- | --- | --- |")
	for range checkButNoSync {
		fmt.Print(" --- |")
	}

	for _, lang := range derivedLanguages {
		if emptyLanguages[lang] {
			continue
		}

		slug := "#non-curated-languages"
		if reportedLanguages[lang] {
			slug = strings.ReplaceAll(strings.ToLower(lang+" "+display.Self.Name(translation.FromSteamLanguage[lang])), " ", "-")
		}

		fmt.Printf("\n| [%s](#%s \"%s\") |", translation.FromSteamLanguage[lang].String(), slug, display.Self.Name(translation.FromSteamLanguage[lang]))

		// for strings, count each incomplete string; don't count entirely missing
		// files as we auto-sync these, but report them below in case we mess up
		incomplete := 0
		for _, file := range sortedLanguageFiles {
			incomplete += file.indented[lang]
		}

		if incomplete == 0 {
			fmt.Print(" ✓ |")
		} else {
			fmt.Printf(" %d |", incomplete)
		}

		// for other files, count each file as one unit
		for _, files := range checked {
			incomplete = 0

			for _, file := range files {
				indented, ok := file.indented[lang]
				if !ok || indented != 0 {
					incomplete++
				}
			}

			if incomplete == 0 {
				fmt.Print(" ✓ |")
			} else {
				fmt.Printf(" %d |", incomplete)
			}
		}

		// count each missing key in the inventory schema
		incomplete = 0
		for _, item := range inventoryItems {
			for _, prefix := range inventoryKeyPrefixes {
				_, ok1 := item[prefix+sourceLanguage]
				_, ok2 := item[prefix+lang]
				if ok1 && !ok2 {
					incomplete++
				}
			}
		}

		if incomplete == 0 {
			fmt.Print(" ✓ |")
		} else {
			fmt.Printf(" %d |", incomplete)
		}
	}

	fmt.Print("\n# Per-File Breakdown\n\n")

	anyNonReported := false
	for _, lang := range derivedLanguages {
		if !reportedLanguages[lang] {
			anyNonReported = true
			continue
		}

		fmt.Printf("## %s (%s)\n\n", lang, display.Self.Name(translation.FromSteamLanguage[lang]))

		any := false

		anyStrings := false
		for _, file := range sortedLanguageFiles {
			indented, ok := file.indented[lang]
			if ok && indented == 0 {
				continue
			}

			if !anyStrings {
				any = true
				anyStrings = true

				fmt.Printf("<details><summary>Strings</summary>\n\n")
			}

			name := file.prefix + lang + file.suffix
			sourceName := file.prefix + sourceLanguage + file.suffix

			if !ok {
				fmt.Printf("- [%s](%s) is missing.\n", filepath.Base(name), sourceName)
			} else {
				fmt.Printf("- [%s](%s) has %d untranslated strings.\n", filepath.Base(name), name, indented)
			}
		}

		if anyStrings {
			fmt.Printf("\n</details>\n\n")
		}

		for i, files := range checked {
			anyOther := false

			for _, file := range files {
				indented, ok := file.indented[lang]
				if ok && indented == 0 {
					continue
				}

				if !anyOther {
					any = true
					anyOther = true

					fmt.Printf("<details><summary>%s</summary>\n\n", checkButNoSync[i].category)
				}

				name := file.prefix + lang + file.suffix
				sourceName := file.prefix + sourceLanguage + file.suffix

				if !ok {
					fmt.Printf("- [%s](%s) is missing.\n", filepath.Base(name), sourceName)
				} else {
					fmt.Printf("- [%s](%s) has %d indented lines.\n", filepath.Base(name), name, indented)
				}
			}

			if anyOther {
				fmt.Printf("\n</details>\n\n")
			}
		}

		anyItems := false
		for _, item := range inventoryItems {
			anyThisItem := false

			for _, prefix := range inventoryKeyPrefixes {
				_, ok1 := item[prefix+sourceLanguage]
				_, ok2 := item[prefix+lang]
				if ok1 && !ok2 {
					if !anyItems {
						any = true
						anyItems = true
						fmt.Print("<details><summary>Inventory Schema</summary>\n\n")
					}

					if !anyThisItem {
						anyThisItem = true
						file := item["_file"].(string)
						fmt.Printf("- Item %s %q in [`%s`](%s) is missing `%s`", item["itemdefid"], item["name_"+sourceLanguage], filepath.Base(file), file, prefix+lang)
					} else {
						fmt.Printf(", `%s`", prefix+lang)
					}
				}
			}

			if anyThisItem {
				fmt.Println()
			}
		}

		if anyItems {
			fmt.Print("\n</details>\n\n")
		}

		if !any {
			fmt.Print("✓ Up to date!\n")
		}
	}

	if anyNonReported {
		fmt.Print("# Non-Curated Languages\n\nThese languages have not been substantially updated since the start of the Reactive Drop translation project. If you are fluent in one of these (or the ones above) and would like to contribute, don't hesitate to ask for directions :)\n\n")

		for _, lang := range derivedLanguages {
			if !reportedLanguages[lang] {
				fmt.Printf("- %s (%s)\n", lang, display.Self.Name(translation.FromSteamLanguage[lang]))
			}
		}
	} else {
		fmt.Print("# Non-Curated Languages\n\nAlien Swarm: Reactive Drop currently supports [every language with full Steam platform support](https://partner.steamgames.com/doc/store/localization/languages)! If you are fluent in a language above and would like to contribute, don't hesitate to ask for directions :)\n")
	}
}
