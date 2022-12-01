package main

import (
	"os"

	"git.lubar.me/ben/valve/vdf"
)

func loadVDF(filename string) (*vdf.KeyValues, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = readBOM(f)
	if err != nil {
		return nil, err
	}

	kv := &vdf.KeyValues{}
	_, err = kv.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return kv, f.Close()
}

func compileSteamInputManifest() {
	manifest, err := loadVDF("../community/steam_input/action_manifest.vdf")
	if err != nil {
		panic(err)
	}

	loc := manifest.FindKey("localization")

	err = addSteamInputManifestLanguage(loc, sourceLanguage)
	if err != nil {
		panic(err)
	}

	for _, lang := range derivedLanguages {
		if emptyLanguages[lang] {
			continue
		}

		err = addSteamInputManifestLanguage(loc, lang)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create("../community/steam_input/configs/steam_input_manifest.vdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = manifest.WriteTo(f)
	if err != nil {
		panic(err)
	}
}

func addSteamInputManifestLanguage(loc *vdf.KeyValues, lang string) error {
	ts, err := loadTranslatedStrings("../community/steam_input/steam_input_"+lang+".vdf", lang)
	if err != nil {
		return err
	}

	locLang := &vdf.KeyValues{Key: lang}
	any := false

	for _, s := range ts.strings {
		if s.indent && s.translated == s.source {
			continue
		}

		any = true
		locLang.AddSubKey(&vdf.KeyValues{
			Key:      s.key,
			Value:    s.translated,
			HasValue: true,
		})
	}

	if any {
		loc.AddSubKey(locLang)
	}

	return nil
}
