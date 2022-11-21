package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.lubar.me/ben/valve/vdf"
)

const sourceLanguage = "english"
const languagePrefix = "[" + sourceLanguage + "]"

// order of this array matches the Steam documentation table order:
// https://partner.steamgames.com/doc/store/localization
var derivedLanguages = [...]string{
	//"arabic",
	//"bulgarian",
	"schinese",
	"tchinese",
	"czech",
	"danish",
	"dutch",
	"finnish",
	"french",
	"german",
	//"greek",
	"hungarian",
	"italian",
	"japanese",
	"koreana",
	"norwegian",
	"polish",
	"portuguese",
	"brazilian",
	"romanian",
	"russian",
	"spanish",
	//"latam",
	"swedish",
	"thai",
	"turkish",
	"ukrainian",
	"vietnamese",
}

var reportedLanguages = map[string]bool{
	"arabic":     false,
	"bulgarian":  false,
	"schinese":   true,
	"tchinese":   true,
	"czech":      false,
	"danish":     false,
	"dutch":      false,
	"finnish":    false,
	"french":     true,
	"german":     true,
	"greek":      false,
	"hungarian":  false,
	"italian":    true,
	"japanese":   true,
	"koreana":    true,
	"norwegian":  false,
	"polish":     true,
	"portuguese": true,
	"brazilian":  true,
	"romanian":   false,
	"russian":    true,
	"spanish":    true,
	"latam":      false,
	"swedish":    false,
	"thai":       false,
	"turkish":    false,
	"ukrainian":  true,
	"vietnamese": true,
}

var txtLanguageFiles = [...]string{
	"../platform/servers/serverbrowser",
	"../platform/vgui",
	"../resource/basemodui",
	"../resource/chat",
	"../resource/closecaption",
	"../resource/gameui",
	"../resource/reactivedrop",
	"../resource/valve",
}

var vdfLanguageFiles = [...]string{
	"../community/inventory_service/inventory_service_tags",
	"../community/statsweb",
}

var txtAddonLanguageFiles = [...]string{
	"../addons/*/resource/closecaption",
	"../addons/*/resource/reactivedrop",
}

var (
	flagMarkdown   = flag.Bool("markdown", false, "generate a markdown table as output")
	flagOnlyUpdate = flag.Bool("only-update", false, "only update source strings; do not sync")
)

func main() {
	flag.Parse()

	if *flagMarkdown {
		numReported := 0
		fmt.Print("| File ")
		for _, lang := range derivedLanguages {
			if reportedLanguages[lang] {
				fmt.Printf("| %s ", lang)
				numReported++
			}
		}
		fmt.Printf("|\n|:- |%s\n", strings.Repeat(":-:|", numReported))
	}

	for _, prefix := range txtLanguageFiles {
		syncTranslations(prefix, ".txt", false)
	}
	for _, prefix := range vdfLanguageFiles {
		syncTranslations(prefix, ".vdf", false)
	}

	updateAchievements(sourceLanguage)
	for _, lang := range derivedLanguages {
		updateAchievements(lang)
	}

	for _, prefix := range txtAddonLanguageFiles {
		addonFiles, err := filepath.Glob(prefix + "_" + sourceLanguage + ".txt")
		if err != nil {
			panic(err)
		}
		for _, file := range addonFiles {
			syncTranslations(strings.TrimSuffix(file, "_"+sourceLanguage+".txt"), ".txt", false)
		}
	}
}

func syncTranslations(prefix, suffix string, quiet bool) {
	lang := sourceLanguage
	if quiet {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s_%s%s:\n", prefix, lang, suffix)
				panic(r)
			}
		}()
	}

	if !quiet {
		if *flagMarkdown {
			fmt.Printf("| `%s` |", prefix)
		} else {
			fmt.Printf("%s\n", prefix)
		}
	}

	base, err := loadVDF(prefix + "_" + sourceLanguage + suffix)
	if err != nil {
		panic(err)
	}

	checkForDuplicateSourceStrings(base)

	for _, lang = range derivedLanguages {
		if !quiet && reportedLanguages[lang] && !*flagMarkdown {
			fmt.Printf("  %10s:", lang)
		}

		upToDate, total := updateLanguageFile(base, prefix, lang, suffix)
		percent := float64(upToDate) / float64(total) * 100

		if !quiet && reportedLanguages[lang] {
			if *flagMarkdown {
				if upToDate >= total {
					fmt.Print(" ✔ |")
				} else {
					fmt.Printf(" %.0f%% (%d untranslated strings) |", percent, total-upToDate)
				}
			} else {
				if upToDate >= total {
					fmt.Println(" ✔")
				} else {
					fmt.Printf("%8.0f%% (%d untranslated strings)\n", percent, total-upToDate)
				}
			}
		}
	}

	if !quiet {
		fmt.Println()
	}
}

func readBOM(r io.Reader) error {
	var bom [3]byte

	n, err := r.Read(bom[:])
	if err != nil {
		return err
	}

	if n != 3 || bom[0] != 0xEF || bom[1] != 0xBB || bom[2] != 0xBF {
		return fmt.Errorf("vdf file does not begin with a BOM (byte order mark) or is not UTF-8")
	}

	return nil
}

func loadVDF(name string) (*vdf.KeyValues, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = readBOM(f)
	if err != nil {
		return nil, err
	}

	var kv vdf.KeyValues

	_, err = kv.ReadFrom(f)

	return &kv, err
}

var everSeenString = make(map[string]bool)

func checkForDuplicateSourceStrings(source *vdf.KeyValues) {
	thisFile := make(map[string]bool)
	for c := source.FindKey("Tokens").FirstValue(); c != nil; c = c.NextValue() {
		lowerKey := strings.ToLower(c.Key)
		if thisFile[lowerKey] {
			panic("String appears multiple times in this file: " + c.Key)
		}
		thisFile[lowerKey] = true

		if everSeenString[lowerKey] {
			panic("Already encountered string in different translation file: " + c.Key)
		}
		everSeenString[lowerKey] = true
	}
}

type translatedString struct {
	translated string
	source     string
	indent     bool
}

func nextRealToken(r *bufio.Reader) (s string, t vdf.Token, indent bool, err error) {
	for {
		s, t, err = vdf.ReadToken(r)
		if err != nil {
			return
		}

		if t == vdf.TokenSpace {
			indent = indent || strings.ContainsRune(s, '\t')
		}

		if t != vdf.TokenSpace && t != vdf.TokenComment {
			return
		}
	}
}

func loadTranslatedStrings(filename, lang string) (map[string]translatedString, error) {
	m := make(map[string]translatedString)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	err = readBOM(r)
	if err != nil {
		return nil, err
	}

	s, t, _, err := nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"lang\"" {
		return nil, fmt.Errorf("expected file to start with \"lang\"")
	}

	_, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenOpenBrace {
		return nil, fmt.Errorf("expected opening brace")
	}

	s, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"Language\"" {
		return nil, fmt.Errorf("expected \"Language\"")
	}

	s, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\""+lang+"\"" {
		return nil, fmt.Errorf("expected language name to match filename")
	}

	s, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"Tokens\"" {
		return nil, fmt.Errorf("expected \"Tokens\"")
	}

	_, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenOpenBrace {
		return nil, fmt.Errorf("expected opening brace")
	}

	for {
		key, t, indent, err := nextRealToken(r)
		if err != nil {
			return nil, err
		}

		if t == vdf.TokenCloseBrace {
			break
		}

		if t != vdf.TokenQuoted {
			return nil, fmt.Errorf("unexpected token type: %q %v", key, t)
		}

		value, t, _, err := nextRealToken(r)
		if err != nil {
			return nil, err
		}

		if t != vdf.TokenQuoted {
			return nil, fmt.Errorf("unexpected token type: %q %v", value, t)
		}

		key = key[1 : len(key)-1]
		key = vdf.Unescape.Replace(key)

		value = value[1 : len(value)-1]
		value = vdf.Unescape.Replace(value)

		key = strings.ToLower(key)
		if strings.HasPrefix(key, languagePrefix) {
			key = key[len(languagePrefix):]
			x := m[key]
			x.source = value
			x.indent = x.indent || indent
			m[key] = x
		} else {
			x := m[key]
			x.translated = value
			x.indent = x.indent || indent
			m[key] = x
		}
	}

	_, t, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenCloseBrace {
		return nil, fmt.Errorf("expected closing brace")
	}

	_, _, _, err = nextRealToken(r)
	if err != io.EOF {
		return nil, fmt.Errorf("expected EOF")
	}

	return m, nil
}

func check(n int, err error) {
	if err != nil {
		panic(err)
	}
}

func updateLanguageFile(source *vdf.KeyValues, prefix, lang, suffix string) (upToDate, total int) {
	filename := prefix + "_" + lang + suffix
	dest, err := loadTranslatedStrings(filename, lang)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	check(buf.WriteString("\ufeff\"lang\"\r\n{\r\n\"Language\"\t\t\""))
	check(vdf.Escape.WriteString(&buf, lang))
	check(buf.WriteString("\"\r\n\"Tokens\"\r\n{\r\n"))

	seen := make(map[string]bool)

	for c := source.FindKey("Tokens").FirstValue(); c != nil; c = c.NextValue() {
		if c.Cond != "" {
			panic("unexpected VDF conditional: " + c.String())
		}

		lowerKey := strings.ToLower(c.Key)

		if seen[lowerKey] {
			panic("duplicate translation key: " + c.Key)
		}
		seen[lowerKey] = true

		x := dest[lowerKey]

		if *flagOnlyUpdate {
			if x.translated == x.source {
				x.translated = c.Value
			}

			x.source = c.Value
		} else {
			if x.source == "" && c.Value != "" {
				if x.translated == "" {
					x.translated = c.Value
				}

				x.source = c.Value
				x.indent = true
			}

			if x.source != c.Value && x.translated == c.Value {
				x.source, x.translated = x.translated, x.source
			}

			if x.source != c.Value {
				x.translated = c.Value
				x.source = c.Value
				x.indent = true
			}
		}

		total++
		if !x.indent {
			upToDate++
		}

		if x.indent {
			check(buf.WriteString("\t\t"))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, c.Key))
		check(buf.WriteString("\"\t\t\""))
		check(vdf.Escape.WriteString(&buf, x.translated))
		check(buf.WriteString("\"\r\n"))

		if x.indent {
			check(buf.WriteString("\t\t"))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, languagePrefix))
		check(vdf.Escape.WriteString(&buf, c.Key))
		check(buf.WriteString("\"\t\t\""))
		check(vdf.Escape.WriteString(&buf, x.source))
		check(buf.WriteString("\"\r\n"))
	}

	check(buf.WriteString("}\r\n}\r\n"))

	err = os.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	return
}

func updateAchievements(lang string) {
	kv, err := loadVDF("../resource/reactivedrop_" + lang + ".txt")
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../achievements/563560_loc_" + lang + ".vdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	tokens := kv.FindKey("Tokens")

	check(f.WriteString("\"lang\"\r\n{\r\n\t\"Language\"\t\""))
	check(vdf.Escape.WriteString(f, lang))
	check(f.WriteString("\"\r\n\t\"Tokens\"\r\n\t{\r\n"))

	for _, a := range achievements {
		aName := tokens.FindKey(a.apiName + "_NAME").Value
		aDesc := tokens.FindKey(a.apiName + "_DESC").Value
		origName, origDesc := "", ""
		if lang != sourceLanguage {
			origName = tokens.FindKey(languagePrefix + a.apiName + "_NAME").Value
			origDesc = tokens.FindKey(languagePrefix + a.apiName + "_DESC").Value
		}

		check(fmt.Fprintf(f, "\t\t\"NEW_ACHIEVEMENT_%d_%d_NAME\"\t\"", a.id, a.bit))
		if origName == aName && origDesc == aDesc {
			check(f.WriteString("\" //\""))
		}
		check(vdf.Escape.WriteString(f, aName))

		check(fmt.Fprintf(f, "\"\r\n\t\t\"NEW_ACHIEVEMENT_%d_%d_DESC\"\t\"", a.id, a.bit))
		if origName == aName && origDesc == aDesc {
			check(f.WriteString("\" //\""))
		}
		check(vdf.Escape.WriteString(f, aDesc))
		check(f.WriteString("\"\r\n"))
	}

	check(f.WriteString("\t}\r\n}\r\n"))
}
