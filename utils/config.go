package main

import "regexp"

const sourceLanguage = "english"
const languagePrefix = "[" + sourceLanguage + "]"

// order of this array matches the Steam documentation table order:
// https://partner.steamgames.com/doc/store/localization
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
	"arabic":    true,
	"bulgarian": true,
	"greek":     true,
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

var checkInventorySchema = [...]string{
	"../community/inventory_service/item-schema-*.json",
}

var inventoryKeyPrefixes = [...]string{
	"name_",
	"briefing_name_",
	"description_",
	"ingame_description_",
	"before_description_",
	"after_description_",
	"accessory_description_",
	"display_type_",
	"style_0_name_",
	"style_1_name_",
	"style_2_name_",
	"style_3_name_",
	"style_4_name_",
	"style_5_name_",
	"style_6_name_",
	"style_7_name_",
	"style_8_name_",
	"style_9_name_",
	"style_10_name_",
	"style_11_name_",
	"style_12_name_",
	"style_13_name_",
	"style_14_name_",
	"style_15_name_",
	"style_16_name_",
	"style_17_name_",
	"style_18_name_",
	"style_19_name_",
	"style_20_name_",
	"style_21_name_",
	"style_22_name_",
	"style_23_name_",
	"style_24_name_",
	"style_25_name_",
	"style_26_name_",
	"style_27_name_",
	"style_28_name_",
	"style_29_name_",
	"style_30_name_",
	"style_31_name_",
	"style_32_name_",
	"style_33_name_",
	"style_34_name_",
	"style_35_name_",
	"style_36_name_",
	"style_37_name_",
	"style_38_name_",
	"style_39_name_",
	"style_40_name_",
	"style_41_name_",
	"style_42_name_",
	"style_43_name_",
	"style_44_name_",
	"style_45_name_",
	"style_46_name_",
	"style_47_name_",
	"style_48_name_",
	"style_49_name_",
	"style_50_name_",
	"style_51_name_",
	"style_52_name_",
	"style_53_name_",
	"style_54_name_",
	"style_55_name_",
	"style_56_name_",
	"style_57_name_",
	"style_58_name_",
	"style_59_name_",
	"style_60_name_",
	"style_61_name_",
	"style_62_name_",
	"style_63_name_",
	// 64 styles is enough for anyone --Bill Gates, probably
	"notification_name_0_",
	"notification_name_1_",
}
