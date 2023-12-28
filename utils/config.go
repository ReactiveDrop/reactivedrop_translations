package main

import "regexp"

const sourceLanguage = "english"
const languagePrefix = "[" + sourceLanguage + "]"

// order of this array matches the Steam documentation table order:
// https://partner.steamgames.com/doc/store/localization/languages
var derivedLanguages = [...]string{
	"arabic",
	"bulgarian",
	"schinese",
	"tchinese",
	"czech",
	"danish",
	"dutch",
	"finnish",
	"french",
	"german",
	"greek",
	"hungarian",
	"indonesian",
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
	"latam",
	"swedish",
	"thai",
	"turkish",
	"ukrainian",
	"vietnamese",
}

// languages that have no translated strings and no translators
var emptyLanguages = map[string]bool{
	"arabic":     true,
	"bulgarian":  true,
	"greek":      true,
	"indonesian": true,
}

var reportedLanguages = map[string]bool{
	"arabic":     false,
	"bulgarian":  false,
	"schinese":   true,
	"tchinese":   true,
	"czech":      true,
	"danish":     false,
	"dutch":      false,
	"finnish":    false,
	"french":     true,
	"german":     true,
	"greek":      false,
	"hungarian":  false,
	"indonesian": false,
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
	"latam":      true,
	"swedish":    true,
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
	"../community/inventory_service/items",
	"../community/stats_website/statsweb",
	"../community/steam_input/steam_input",
	"../community/workshop/workshop_description",
	"../misc/signage/signage",
}

var txtAddonLanguageFiles = [...]string{
	"../addons/*/resource/closecaption",
	"../addons/*/resource/reactivedrop",
}

var importantLanguageFiles = [...][2]string{
	{"resource/basemodui_", ".txt"},
	{"resource/closecaption_", ".txt"},
	{"resource/gameui_", ".txt"},
	{"resource/reactivedrop_", ".txt"},
	{"community/inventory_service/items_", ".vdf"},
	{"community/stats_website/statsweb_", ".vdf"},
}

var eventMaxLength = map[string]int{
	"title":    80,
	"subtitle": 120,
	"summary":  180,
	"body":     32000,
}

var stringMaxLength = [...]struct {
	File      string
	Key       *regexp.Regexp
	MaxLength int
}{
	{
		File:      "../community/workshop/workshop_description",
		Key:       regexp.MustCompile(`_title\z`),
		MaxLength: 50,
	},
	{
		File:      "../community/workshop/workshop_description",
		Key:       regexp.MustCompile(`\AWorkshop_desc\z`),
		MaxLength: 200,
	},
	{
		File:      "../community/workshop/workshop_description",
		Key:       regexp.MustCompile(`_desc\z`),
		MaxLength: 1000,
	},
}

var checkButNoSync = [...]struct {
	category string
	patterns []string
}{
	{
		category: "Steam Store and Community",
		patterns: []string{
			"../community/eula/eula_*.txt",
			"../community/points_shop/app_items_*.json",
			"../community/workshop/workshop_tags_*.json",
			"../store_page/content_warning_*.txt",
			"../store_page/storepage_*.json",
		},
	},
	{
		category: "Upcoming Release Notes",
		patterns: []string{
			"../release_notes/*.xml", // not checking archive
		},
	},
	{
		category: "Credits",
		patterns: []string{
			"../credits/*.txt",
			"../addons/*/resource/*.txt",
		},
	},
	{
		category: "Mail and News",
		patterns: []string{
			"../resource/mail/*.txt",
			"../addons/*/resource/mail/*.txt",
			"../resource/news/*.txt",
			"../addons/*/resource/news/*.txt",
		},
	},
}
