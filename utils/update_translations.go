package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"git.lubar.me/ben/valve/vdf"
)

var (
	flagMarkdown      = flag.Bool("markdown", false, "generate a markdown translation progress report")
	flagInputManifest = flag.Bool("input-manifest", false, "compile the steam input manifest")
	flagRender        = flag.Bool("render", false, "render derived files (except the steam input manifest)")
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

	if *flagRender {
		updateDerivedFiles(sourceLanguage)
		for _, lang := range derivedLanguages {
			if emptyLanguages[lang] {
				continue
			}

			updateDerivedFiles(lang)
		}

		renderInventorySchema()
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

	checkReleaseNotes()
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
			comments += s
		case vdf.TokenComment:
			comments += s
		default:
			// remove anything before the newline after the previous string
			if i := strings.IndexByte(comments, '\n'); i != -1 {
				comments = comments[i+1:]
			}

			// remove any indentation in the comment string (it will be re-added after newlines if needed; use spaces in comments)
			comments = strings.ReplaceAll(comments, "\t", "")

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

		for _, checkLen := range stringMaxLength {
			if !strings.HasPrefix(filename, checkLen.File) {
				continue
			}

			if checkLen.Key.MatchString(key) {
				if count := utf8.RuneCountInString(value); count > checkLen.MaxLength {
					panic(fmt.Sprintf("%q cannot be longer than %d characters, but it is %d characters", key, checkLen.MaxLength, count))
				}
				break
			}
		}

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
var unindentBlankLine = strings.NewReplacer("\t\t\r\n", "\r\n")

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

	merge, err := loadTranslatedStrings(prefix+"_"+lang+"_merge"+suffix, lang, true)
	if errors.Is(err, os.ErrNotExist) {
		merge, err = nil, nil
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

		if merge != nil && x.translated == c.translated {
			i, ok = merge.lookup[lowerKey]
			if ok && merge.strings[i].source == c.translated && merge.strings[i].translated != x.translated {
				x.translated = merge.strings[i].translated
				x.tcomment = merge.strings[i].tcomment
				x.indent = true
			}
		}

		total++
		if !x.indent {
			upToDate++
		}

		if x.indent {
			c := indentNewLine.Replace("\t\t" + x.scomment)
			check(unindentBlankLine.WriteString(&buf, c))
		} else {
			check(buf.WriteString(x.scomment))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, languagePrefix))
		check(vdf.Escape.WriteString(&buf, c.key))
		check(buf.WriteString("\"\t\t\""))
		check(vdf.Escape.WriteString(&buf, x.source))
		check(buf.WriteString("\"\r\n"))

		if x.indent {
			c := indentNewLine.Replace("\t\t" + x.tcomment)
			check(unindentBlankLine.WriteString(&buf, c))
		} else {
			check(buf.WriteString(x.tcomment))
		}

		check(1, buf.WriteByte('"'))
		check(vdf.Escape.WriteString(&buf, c.key))
		check(buf.WriteString("\"\t\t\""))
		check(vdf.Escape.WriteString(&buf, x.translated))
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

func renderInventorySchema() {
	shared, err := loadTranslatedStrings("../community/inventory_service/items_shared.vdf", "none", true)
	if err != nil {
		panic(err)
	}

	languages := make([]*inventoryStrings, len(derivedLanguages)+1)
	source, err := loadInventoryStrings(sourceLanguage)
	if err != nil {
		panic(err)
	}

	languages[len(derivedLanguages)] = source

	for i, lang := range derivedLanguages {
		languages[i], err = loadInventoryStrings(lang)
		if err != nil {
			panic(err)
		}
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].lang < languages[j].lang
	})

	names, err := filepath.Glob("../community/inventory_service/template/item-schema-*.json")
	if err != nil {
		panic(err)
	}

	for _, name := range names {
		renderInventorySchemaFile(name, filepath.Join("../community/inventory_service/rendered", filepath.Base(name)), shared, source, languages)
	}
}

func renderInventorySchemaFile(srcName, dstName string, shared *translatedStrings, source *inventoryStrings, languages []*inventoryStrings) {
	src, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	var buf bytes.Buffer

	in := json.NewDecoder(src)
	in.UseNumber()

	delims := make([]json.Delim, 0, 8)
	lastDelim := json.Delim(0)
	expectingValue := false
	eachLangPrefix := ""

	for {
		tok, err := in.Token()
		if err == io.EOF {
			var buf2 bytes.Buffer

			err := json.Indent(&buf2, buf.Bytes(), "", "\t")
			if err != nil {
				panic(err)
			}

			buf2.WriteByte('\n') // trailing newline

			err = os.WriteFile(dstName, buf2.Bytes(), 0644)
			if err != nil {
				panic(err)
			}

			break
		}

		if err != nil {
			panic(err)
		}

		switch t := tok.(type) {
		case json.Delim:
			buf.WriteString(t.String())

			if t == '{' || t == '[' {
				delims = append(delims, t)
			} else if t == '}' {
				if delims[len(delims)-1] != '{' {
					panic("unexpected } after " + delims[len(delims)-1].String())
				}

				delims = delims[:len(delims)-1]

				if in.More() {
					buf.WriteByte(',')
				}
			} else if t == ']' {
				if delims[len(delims)-1] != '[' {
					panic("unexpected ] after " + delims[len(delims)-1].String())
				}

				delims = delims[:len(delims)-1]

				if in.More() {
					buf.WriteByte(',')
				}
			} else {
				panic("unexpected json.Delim " + t.String())
			}

			if len(delims) != 0 {
				lastDelim = delims[len(delims)-1]
			}

			if eachLangPrefix != "" {
				panic("unexpected language prefix")
			}

			expectingValue = false
		case json.Number:
			if lastDelim != '[' && !expectingValue {
				panic("unexpected number")
			}

			if eachLangPrefix != "" {
				panic("unexpected language prefix")
			}

			buf.WriteString(t.String())
			if in.More() {
				buf.WriteByte(',')
			}

			expectingValue = false
		case string:
			if lastDelim == '[' || expectingValue {
				langs := languages
				if eachLangPrefix == "" {
					langs = []*inventoryStrings{source}
				}

				sourceString := performInventorySchemaReplacements(t, shared, source, source)
				if eachLangPrefix != "" && sourceString == t {
					panic("no tokens in " + srcName + " field " + eachLangPrefix + "%LANG% (value " + strconv.Quote(t) + ")")
				}

				for _, lang := range langs {
					str := performInventorySchemaReplacements(t, shared, source, lang)

					if eachLangPrefix != "" {
						if lang != source && sourceString == str {
							continue
						}

						b, err := json.Marshal(eachLangPrefix + lang.lang)
						if err != nil {
							panic(err)
						}

						buf.Write(b)
						buf.WriteByte(':')
					}

					b, err := json.Marshal(str)
					if err != nil {
						panic(err)
					}

					buf.Write(b)

					if in.More() {
						buf.WriteByte(',')
					}
				}

				expectingValue = false
				eachLangPrefix = ""
			} else {
				if strings.HasSuffix(t, "_"+sourceLanguage) {
					panic("unexpected language suffix in template " + srcName + " (" + t + ")")
				}
				for _, lang := range derivedLanguages {
					if strings.HasSuffix(t, "_"+lang) {
						panic("unexpected language suffix in template " + srcName + " (" + t + ")")
					}
				}

				if strings.HasSuffix(t, "_%LANG%") {
					eachLangPrefix = t[:len(t)-len("%LANG%")]
				} else {
					b, err := json.Marshal(t)
					if err != nil {
						panic(err)
					}

					buf.Write(b)

					buf.WriteByte(':')
				}
				expectingValue = true
			}
		case bool:
			if lastDelim != '[' && !expectingValue {
				panic("unexpected number")
			}

			if eachLangPrefix != "" {
				panic("unexpected language prefix")
			}

			if t {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}

			if in.More() {
				buf.WriteByte(',')
			}

			expectingValue = false
		default:
			panic(fmt.Sprintf("unexpected %T %v", t, t))
		}
	}
}

var inventorySchemaPlaceholder = regexp.MustCompile(`\{#[a-z0-9_]+\}`)

func performInventorySchemaReplacements(str string, shared *translatedStrings, source, derived *inventoryStrings) string {
	return inventorySchemaPlaceholder.ReplaceAllStringFunc(str, func(placeholder string) string {
		token := placeholder[2 : len(placeholder)-1]

		i, ok := shared.lookup[token]
		if ok {
			return shared.strings[i].translated
		}

		for _, f := range derived.files {
			if f != nil {
				i, ok := f.lookup[token]
				if ok {
					return f.strings[i].translated
				}
			}
		}

		for _, f := range source.files {
			if f != nil {
				i, ok := f.lookup[token]
				if ok {
					return f.strings[i].translated
				}
			}
		}

		panic("missing string in inventory schema: " + placeholder)
	})
}

type inventoryStrings struct {
	lang  string
	files [3]*translatedStrings
}

func loadInventoryStrings(lang string) (*inventoryStrings, error) {
	reactivedrop, err := loadTranslatedStrings("../resource/reactivedrop_"+lang+".txt", lang, true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	tags, err := loadTranslatedStrings("../community/inventory_service/inventory_service_tags_"+lang+".vdf", lang, true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	items, err := loadTranslatedStrings("../community/inventory_service/items_"+lang+".vdf", lang, true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &inventoryStrings{
		lang: lang,
		files: [...]*translatedStrings{
			reactivedrop,
			tags,
			items,
		},
	}, nil
}

type ReleaseNotes struct {
	XMLName xml.Name `xml:"content"`
	Strings []struct {
		ID   string `xml:"id,attr"`
		Text string `xml:",chardata"`
	} `xml:"string"`
}

func checkReleaseNotes() {
	names, err := filepath.Glob("../release_notes/*.xml")
	if err != nil {
		panic(err)
	}

	englishLines := make(map[string]int)
	for _, name := range names {
		if strings.HasSuffix(name, "_english.xml") {
			var content ReleaseNotes

			b, err := os.ReadFile(name)
			if err != nil {
				panic(err)
			}

			err = xml.Unmarshal(b, &content)
			if err != nil {
				fmt.Println(name)
				panic(err)
			}

			for _, s := range content.Strings {
				if s.ID == "body" {
					lines := strings.FieldsFunc(s.Text, func(ch rune) bool { return ch == '\n' || ch == '\r' })
					englishLines[name[:len(name)-len("_english.xml")]] = len(lines)
					break
				}
			}
		}
	}

	for _, name := range names {
		var content ReleaseNotes

		b, err := os.ReadFile(name)
		if err != nil {
			panic(err)
		}

		err = xml.Unmarshal(b, &content)
		if err != nil {
			fmt.Println(name)
			panic(err)
		}

		for _, s := range content.Strings {
			count := utf8.RuneCountInString(s.Text)
			max := eventMaxLength[s.ID]

			if max == 0 {
				panic(fmt.Sprintf("%s: unexpected string id %q", name, s.ID))
			}

			if count > max {
				panic(fmt.Sprintf("%s: %q cannot be longer than %d characters, but it is %d characters", name, s.ID, max, count))
			}

			if s.ID == "body" {
				lines := strings.FieldsFunc(s.Text, func(ch rune) bool { return ch == '\n' || ch == '\r' })
				expected := englishLines[name[:strings.LastIndexByte(name, '_')]]

				if expected == 0 {
					panic(fmt.Sprintf("%s: cannot find English file for this release notes document", name))
				}

				if len(lines) != expected {
					panic(fmt.Sprintf("%s: body of release notes has %d lines but English version has %d lines (this is usually Ben's fault)", name, len(lines), expected))
				}
			}
		}
	}
}
