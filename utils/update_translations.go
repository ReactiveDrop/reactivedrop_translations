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

var (
	flagMarkdown      = flag.Bool("markdown", false, "generate a markdown translation progress report")
	flagInputManifest = flag.Bool("input-manifest", false, "compile the steam input manifest")
	flagOnlyUpdate    = flag.Bool("only-update", false, "only update source strings; do not reset differing translations")
)

func main() {
	flag.Parse()

	if *flagMarkdown {
		generateReport()

		return
	}

	if *flagInputManifest {
		compileSteamInputManifest()

		return
	}

	for _, prefix := range txtLanguageFiles {
		syncTranslations(prefix, ".txt", false)
	}
	for _, prefix := range vdfLanguageFiles {
		syncTranslations(prefix, ".vdf", false)
	}

	updateDerivedFiles(sourceLanguage)
	for _, lang := range derivedLanguages {
		if emptyLanguages[lang] {
			continue
		}

		updateDerivedFiles(lang)
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
		fmt.Printf("%s\n", prefix)
	}

	base, err := loadTranslatedStrings(prefix+"_"+sourceLanguage+suffix, sourceLanguage, true)
	if err != nil {
		panic(err)
	}

	for key, i := range base.lookup {
		if everSeenString[key] {
			panic("Already encountered string in different translation file: " + base.strings[i].key)
		}

		everSeenString[key] = true
	}

	for _, lang = range derivedLanguages {
		if emptyLanguages[lang] {
			continue
		}

		if !quiet && reportedLanguages[lang] {
			fmt.Printf("  %10s:", lang)
		}

		upToDate, total := updateLanguageFile(base, prefix, lang, suffix)
		percent := float64(upToDate) / float64(total) * 100

		if !quiet && reportedLanguages[lang] {
			if upToDate >= total {
				fmt.Println(" ✔")
			} else {
				fmt.Printf("%8.0f%% (%d untranslated strings)\n", percent, total-upToDate)
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

var everSeenString = make(map[string]bool)

type translatedString struct {
	key        string
	translated string
	tcomment   string
	source     string
	scomment   string
	tseen      bool
	sseen      bool
	indent     bool
}

type translatedStrings struct {
	strings []translatedString
	lookup  map[string]int
}

func nextRealToken(r *bufio.Reader) (s string, t vdf.Token, indent bool, comments string, err error) {
	for {
		s, t, err = vdf.ReadToken(r)
		if err != nil {
			return
		}

		switch t {
		case vdf.TokenSpace:
			indent = indent || strings.ContainsRune(s, '\t')
		case vdf.TokenComment:
			comments += s
		default:
			return
		}
	}
}

func loadTranslatedStrings(filename, lang string, wantBOM bool) (*translatedStrings, error) {
	ts := &translatedStrings{lookup: make(map[string]int)}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	if wantBOM {
		err = readBOM(r)
		if err != nil {
			return nil, err
		}
	}

	s, t, _, _, err := nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"lang\"" {
		return nil, fmt.Errorf("expected file to start with \"lang\"")
	}

	_, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenOpenBrace {
		return nil, fmt.Errorf("expected opening brace")
	}

	s, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"Language\"" {
		return nil, fmt.Errorf("expected \"Language\"")
	}

	s, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\""+lang+"\"" {
		return nil, fmt.Errorf("expected language name to match filename")
	}

	s, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenQuoted || s != "\"Tokens\"" {
		return nil, fmt.Errorf("expected \"Tokens\"")
	}

	_, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenOpenBrace {
		return nil, fmt.Errorf("expected opening brace")
	}

	for {
		key, t, indent, comment, err := nextRealToken(r)
		if err != nil {
			return nil, err
		}

		if t == vdf.TokenCloseBrace {
			break
		}

		if t != vdf.TokenQuoted {
			return nil, fmt.Errorf("unexpected token type: %q %v", key, t)
		}

		value, t, _, _, err := nextRealToken(r)
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

		origKey := key
		key = strings.ToLower(key)
		if strings.HasPrefix(key, languagePrefix) {
			key = key[len(languagePrefix):]

			i, ok := ts.lookup[key]
			if !ok {
				i = len(ts.strings)
				ts.lookup[key] = i
				ts.strings = append(ts.strings, translatedString{
					key: origKey[len(languagePrefix):],
				})
			}

			x := &ts.strings[i]
			if x.sseen {
				return nil, fmt.Errorf("string appears multiple times: %q", origKey)
			}

			x.source = value
			x.scomment = comment
			x.sseen = true
			x.indent = x.indent || indent
		} else {
			i, ok := ts.lookup[key]
			if !ok {
				i = len(ts.strings)
				ts.lookup[key] = i
				ts.strings = append(ts.strings, translatedString{
					key: origKey,
				})
			}

			x := &ts.strings[i]
			if x.tseen {
				return nil, fmt.Errorf("string appears multiple times: %q", origKey)
			}

			x.translated = value
			x.tcomment = comment
			x.tseen = true
			x.indent = x.indent || indent
		}
	}

	_, t, _, _, err = nextRealToken(r)
	if err != nil {
		return nil, err
	}

	if t != vdf.TokenCloseBrace {
		return nil, fmt.Errorf("expected closing brace")
	}

	_, _, _, _, err = nextRealToken(r)
	if err != io.EOF {
		return nil, fmt.Errorf("expected EOF")
	}

	return ts, nil
}

func check(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

var indentNewLine = strings.NewReplacer("\n", "\n\t\t")

func updateLanguageFile(source *translatedStrings, prefix, lang, suffix string) (upToDate, total int) {
	filename := prefix + "_" + lang + suffix
	dest, err := loadTranslatedStrings(filename, lang, true)
	if errors.Is(err, os.ErrNotExist) {
		dest = &translatedStrings{
			lookup:  make(map[string]int, len(source.strings)),
			strings: make([]translatedString, 0, len(source.strings)),
		}
		err = nil
	}
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	check(buf.WriteString("\ufeff\"lang\"\r\n{\r\n\"Language\"\t\t\""))
	check(vdf.Escape.WriteString(&buf, lang))
	check(buf.WriteString("\"\r\n\"Tokens\"\r\n{\r\n"))

	for _, c := range source.strings {
		lowerKey := strings.ToLower(c.key)

		var x translatedString
		i, ok := dest.lookup[lowerKey]
		if ok {
			x = dest.strings[i]
		}
		x.scomment = c.tcomment

		if *flagOnlyUpdate {
			if x.translated == x.source {
				x.translated = c.translated
			}

			x.source = c.translated
		} else {
			if x.source == "" && c.translated != "" {
				if x.translated == "" {
					x.translated = c.translated
				}

				x.source = c.translated
				x.indent = true
			}

			if x.source != c.translated && x.translated == c.translated {
				x.source, x.translated = x.translated, x.source
			}

			if x.source != c.translated {
				x.translated = c.translated
				x.source = c.translated
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

		if x.indent {
			check(indentNewLine.WriteString(&buf, x.tcomment))
		} else {
			check(buf.WriteString(x.tcomment))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, c.key))
		check(buf.WriteString("\"\t\t\""))
		check(vdf.Escape.WriteString(&buf, x.translated))
		check(buf.WriteString("\"\r\n"))

		if x.indent {
			check(buf.WriteString("\t\t"))
		}

		if x.indent {
			check(indentNewLine.WriteString(&buf, x.scomment))
		} else {
			check(buf.WriteString(x.scomment))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, languagePrefix))
		check(vdf.Escape.WriteString(&buf, c.key))
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

func updateDerivedFiles(lang string) {
	updateAchievements(lang)
	updateRichPresence(lang)
}

func updateAchievements(lang string) {
	kv, err := loadTranslatedStrings("../resource/reactivedrop_"+lang+".txt", lang, true)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../achievements/563560_loc_" + lang + ".vdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	check(f.WriteString("\"lang\"\r\n{\r\n\t\"Language\"\t\""))
	check(vdf.Escape.WriteString(f, lang))
	check(f.WriteString("\"\r\n\t\"Tokens\"\r\n\t{\r\n"))

	findString := func(key string, useSource bool) string {
		i, ok := kv.lookup[strings.ToLower(key)]
		if !ok {
			return ""
		}

		if useSource {
			return kv.strings[i].source
		}

		return kv.strings[i].translated
	}

	for _, a := range achievements {
		aName := findString(a.apiName+"_NAME", false)
		aDesc := findString(a.apiName+"_DESC", false)
		origName, origDesc := "", ""
		if lang != sourceLanguage {
			origName = findString(a.apiName+"_NAME", true)
			origDesc = findString(a.apiName+"_DESC", true)
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

func updateRichPresence(lang string) {
	kv, err := loadTranslatedStrings("../resource/reactivedrop_"+lang+".txt", lang, true)
	if err != nil {
		panic(err)
	}

	oldStrings, err := loadTranslatedStrings("../rich_presence/563560_loc_"+lang+".vdf", lang, false)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../rich_presence/563560_loc_" + lang + ".vdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	check(f.WriteString("\"lang\"\r\n{\r\n\t\"Language\"\t\""))
	check(vdf.Escape.WriteString(f, lang))
	check(f.WriteString("\"\r\n\t\"Tokens\"\r\n\t{\r\n"))

	for _, s := range oldStrings.strings {
		if !strings.HasPrefix(s.key, "#official_") {
			check(f.WriteString("\t\t\""))
			check(vdf.Escape.WriteString(f, s.key))
			check(f.WriteString("\"\t\""))
			check(vdf.Escape.WriteString(f, s.translated))
			check(f.WriteString("\"\r\n"))
		}
	}

	findString := func(key string) string {
		i, ok := kv.lookup[strings.ToLower(key)]
		if !ok {
			return ""
		}

		return kv.strings[i].translated
	}

	check(f.WriteString("\r\n"))

	for _, s := range officialChallenges {
		name := findString("rd_challenge_name_" + s)

		check(fmt.Fprintf(f, "\t\t\"#official_challenge_%s\"\t\"", s))
		check(vdf.Escape.WriteString(f, name))
		check(f.WriteString("\"\r\n"))
	}

	check(f.WriteString("\r\n"))

	for _, s := range officialMissions {
		name := findString("rd_mission_title_" + s)

		check(fmt.Fprintf(f, "\t\t\"#official_mission_%s\"\t\"", s))
		check(vdf.Escape.WriteString(f, name))
		check(f.WriteString("\"\r\n"))
	}

	check(f.WriteString("\t}\r\n}\r\n"))
}
