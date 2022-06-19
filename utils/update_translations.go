package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"git.lubar.me/ben/valve/vdf"
)

const onlyUpdateSourceStrings = false
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

var txtLanguageFiles = [...]string{
	"../resource/basemodui",
	"../resource/chat",
	"../resource/closecaption",
	"../resource/gameui",
	"../resource/reactivedrop",
	"../resource/valve",
}

var vdfLanguageFiles = [...]string{
	"../community/points_shop",
	"../community/statsweb",
}

func main() {
	for _, prefix := range txtLanguageFiles {
		syncTranslations(prefix, ".txt")
	}
	for _, prefix := range vdfLanguageFiles {
		syncTranslations(prefix, ".vdf")
	}

	updateAchievements(sourceLanguage)
	for _, lang := range derivedLanguages {
		updateAchievements(lang)
	}
}

func syncTranslations(prefix, suffix string) {
	fmt.Printf("%s:\n", prefix)

	base, err := loadVDF(prefix + "_" + sourceLanguage + suffix)
	if err != nil {
		panic(err)
	}

	for _, lang := range derivedLanguages {
		fmt.Printf("  %10s: ", lang)

		upToDate, total := updateLanguageFile(base, prefix, lang, suffix)
		percent := float64(upToDate) / float64(total) * 100

		fmt.Printf("% 7.3f%%\n", percent)
	}

	fmt.Println()
}

func readBOM(r io.Reader) error {
	var bom [3]byte

	n, err := r.Read(bom[:])
	if err != nil {
		return err
	}

	if n != 3 || bom[0] != 0xEF || bom[1] != 0xBB || bom[2] != 0xBF {
		return fmt.Errorf("vdf file does not begin with a BOM")
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

type translatedString struct {
	translated string
	source     string
	indent     bool
}

func nextRealToken(r *bufio.Reader) (string, vdf.Token, bool, error) {
	indent := false

	for {
		s, t, err := vdf.ReadToken(r)
		if err != nil {
			return s, t, indent, err
		}

		if t == vdf.TokenSpace {
			indent = indent || strings.ContainsRune(s, '\t')
		}

		if t != vdf.TokenSpace && t != vdf.TokenComment {
			return s, t, indent, err
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

	s, t, _, err = nextRealToken(r)
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

	s, t, _, err = nextRealToken(r)
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

	buf.WriteString("\ufeff\"lang\"\r\n{\r\n\"Language\"\t\t\"")
	vdf.Escape.WriteString(&buf, lang)
	buf.WriteString("\"\r\n\"Tokens\"\r\n{\r\n")

	for c := source.FindKey("Tokens").FirstValue(); c != nil; c = c.NextValue() {
		if c.Cond != "" {
			panic("unexpected VDF conditional: " + c.String())
		}

		x := dest[strings.ToLower(c.Key)]

		if onlyUpdateSourceStrings {
			x.source = c.Value
		} else {
			if x.source == "" && c.Value != "" {
				if x.translated == "" {
					x.translated = c.Value
				}

				x.source = c.Value
				x.indent = true
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
			buf.WriteString("\t\t")
		}

		buf.WriteByte('"')
		vdf.Escape.WriteString(&buf, c.Key)
		buf.WriteString("\"\t\t\"")
		vdf.Escape.WriteString(&buf, x.translated)
		buf.WriteString("\"\r\n")

		if x.indent {
			buf.WriteString("\t\t")
		}

		buf.WriteByte('"')
		vdf.Escape.WriteString(&buf, languagePrefix)
		vdf.Escape.WriteString(&buf, c.Key)
		buf.WriteString("\"\t\t\"")
		vdf.Escape.WriteString(&buf, x.source)
		buf.WriteString("\"\r\n")
	}

	buf.WriteString("}\r\n}\r\n")

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

	f.WriteString("\"lang\"\r\n{\r\n\t\"Language\"\t\"")
	vdf.Escape.WriteString(f, lang)
	f.WriteString("\"\r\n\t\"Tokens\"\r\n\t{\r\n")

	for _, a := range achievements {
		aName := tokens.FindKey(a.apiName + "_NAME").Value
		aDesc := tokens.FindKey(a.apiName + "_DESC").Value
		origName, origDesc := "", ""
		if lang != sourceLanguage {
			origName = tokens.FindKey(languagePrefix + a.apiName + "_NAME").Value
			origDesc = tokens.FindKey(languagePrefix + a.apiName + "_DESC").Value
		}

		fmt.Fprintf(f, "\t\t\"NEW_ACHIEVEMENT_%d_%d_NAME\"\t\"", a.id, a.bit)
		if origName == aName && origDesc == aDesc {
			f.WriteString("\" //\"")
		}
		vdf.Escape.WriteString(f, aName)

		fmt.Fprintf(f, "\"\r\n\t\t\"NEW_ACHIEVEMENT_%d_%d_DESC\"\t\"", a.id, a.bit)
		if origName == aName && origDesc == aDesc {
			f.WriteString("\" //\"")
		}
		vdf.Escape.WriteString(f, aDesc)
		f.WriteString("\"\r\n")
	}

	f.WriteString("\t}\r\n}\r\n")
}
