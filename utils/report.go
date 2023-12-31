package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"git.lubar.me/ben/valve/translation"
	"golang.org/x/text/language/display"
)

var commonDirectoryRemover = regexp.MustCompile(`resource/|\Acommunity/[a-z_]+/|\Astore_page/|\Arelease_notes/|\Acredits/`)

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
		case '\r':
			// ignore
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

	fmt.Print("# Overview\n")
	fmt.Print("| Language |")
	for _, file := range importantLanguageFiles {
		fmt.Printf(" %s |", filepath.Base(file[0][:len(file[0])-1]))
	}
	fmt.Print(" Strings | Files |\n| --- | --- | --- |")
	for range importantLanguageFiles {
		fmt.Print(" --- |")
	}

	for _, lang := range derivedLanguages {
		if emptyLanguages[lang] {
			continue
		}

		slug := "non-curated-languages"
		if reportedLanguages[lang] {
			slug = strings.ReplaceAll(strings.ToLower(lang+" "+display.Self.Name(translation.FromSteamLanguage[lang])), " ", "-")
		}

		fmt.Printf("\n| [%s](#%s) |", display.English.Tags().Name(translation.FromSteamLanguage[lang]), slug)

		isImportantFile := make(map[[2]string]bool)
		for _, file := range importantLanguageFiles {
			isImportantFile[file] = true
			if languageFiles[file][lang] == 0 {
				fmt.Print(" ✔️ |")
			} else {
				fmt.Printf(" %d |", languageFiles[file][lang])
			}
		}

		// for strings, count each incomplete string; don't count entirely missing
		// files as we auto-sync these, but report them below in case we mess up
		incomplete := 0
		for _, file := range sortedLanguageFiles {
			if !isImportantFile[[2]string{file.prefix, file.suffix}] {
				incomplete += file.indented[lang]
			}
		}

		if incomplete == 0 {
			fmt.Print(" ✔️ |")
		} else {
			fmt.Printf(" %d |", incomplete)
		}

		// for other files, count each file as one unit
		incomplete = 0
		for _, files := range checked {
			for _, file := range files {
				indented, ok := file.indented[lang]
				if !ok || indented != 0 {
					incomplete++
				}
			}
		}

		if incomplete == 0 {
			fmt.Print(" ✔️ |")
		} else {
			fmt.Printf(" %d |", incomplete)
		}
	}

	fmt.Print("\n### Legend\n- ***Non-capitalized column headers*** are the most important specific files and their number of missing strings. All of these txt-files are located in the resources folder. Except statsweb, which lies in community/stats_website, it's an vdf-file.\n- ***Strings*** is the number of missing strings not included in one of the files which get an individual non-capitalized column.\n- ***Inventory*** is the number of missing language-specific keys in the inventory schema.\n- ***Files*** is the total number of other files that are missing or need updating.")

	fmt.Print("\n\n# Per-File Breakdown\n\n")

	anyNonReported := false
	for _, lang := range derivedLanguages {
		if !reportedLanguages[lang] {
			anyNonReported = true
			continue
		}

		fmt.Printf("<details><summary>\n\n## %s (%s)\n\n</summary>\n", lang, display.Self.Name(translation.FromSteamLanguage[lang]))

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

				fmt.Printf("\n### Strings\n\n")
			}

			name := file.prefix + lang + file.suffix
			sourceName := file.prefix + sourceLanguage + file.suffix

			if !ok {
				fmt.Printf("- [%s](%s) is missing.\n", commonDirectoryRemover.ReplaceAllLiteralString(name, ""), sourceName)
			} else {
				fmt.Printf("- [%s](%s) has %d untranslated strings.\n", commonDirectoryRemover.ReplaceAllLiteralString(name, ""), name, indented)
			}
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

					fmt.Printf("\n### %s\n\n", checkButNoSync[i].category)
				}

				name := file.prefix + lang + file.suffix
				sourceName := file.prefix + sourceLanguage + file.suffix

				if !ok {
					fmt.Printf("- [%s](%s) is missing.\n", commonDirectoryRemover.ReplaceAllLiteralString(name, ""), sourceName)
				} else {
					fmt.Printf("- [%s](%s) has %d indented lines.\n", commonDirectoryRemover.ReplaceAllLiteralString(name, ""), name, indented)
				}
			}
		}

		if !any {
			fmt.Print("\n✓ Up to date!\n")
		}

		fmt.Print("\n</details>\n\n")
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
