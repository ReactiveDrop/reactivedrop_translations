# Alien Swarm: Reactive Drop Translations

This repository contains all the files required to translate [Alien Swarm: Reactive Drop](https://store.steampowered.com/app/563560/) user interface into any supported language. Reactive Drop(RD) supports all [languages which are supported by Steam](https://partner.steamgames.com/doc/store/localization#supported_languages). You can also translate to an unsupported language but this will only work as a VPK add-on file through workshop.

### Pre-requisites

- Any unicode plain text editor. E.g. [Notepad++](https://notepad-plus-plus.org/) or [Visual Studio Code](https://code.visualstudio.com/)
- Fork the repository, add your changes in your own fork, and then submit a pull request to merge in your changes.
- If you have problems using forks or GIT you can submit your translation through the Issues tab by providing an archive.

### Directory structure

- achievements: Steam Achievement descriptions which will be shown on Steams website.
- addons: Translations for custom campaigns and challenges.
- community: Steam Community (workshop, inventory, points shop) translations.
- credits: Credits texts for official campaigns.
- platform: Translations for generic source-engine UI ingame eg. the server browser.
- release_notes: Update notes shown on Steams website.
- resource: Game localization files. These are UTF-8 encoded, whereas the files in `\Alien Swarm Reactive Drop\reactivedrop\resource` are UTF-16LE encoded.
- rich_presence: Online status messages in Steams Friends List.
- store_page: [Steam's store page](https://store.steampowered.com/app/563560/)
- utils: Helper tools for translators.

### How to download the repository

If you are not familiar with using GIT you can [download](https://github.com/ReactiveDrop/reactivedrop_translations/archive/refs/heads/master.zip) the repository as a ZIP archive and submit it using [Issues](https://github.com/ReactiveDrop/reactivedrop_translations/issues) tab, or edit the files in GitHub's web editor.

### How to translate?

### File Encoding

Please make sure that your text editor preserves file's encoding. The new files you create need to have same encoding as their English counterpart(e.g. labsmail1_czech.txt needs to be UTF-8 encoded just like labsmail1_english.txt). If you use Notepad++ you can see the file's encoding in the menu Encoding.

### Achievements

These files are automatically created based on the `resource/reactivedrop_*.txt` files.

### Mail and News

Create a copy of each mail and news file and rename it to your language sufix(e.g. labsmail1_russian.txt). Translate the contents of each file. See labsmail1_russian.txt as an example.

### BaseModUi, CloseCaption, GameUi, ReactiveDrop, and Community VDF files

In these files(e.g. basemodui_czech.txt) the untranslated strings are indented by two tabs to be easily visible. Translate each indented line and remove the two tabs in front of it.

### Workshop

Create a copy of the English file and rename it to your language suffix (e.g. workshop_tags_schinese.json). Translate the contents of the file. In JSON files, only translate text on the right side of the `:`.

### Item Schema

Copy the lines with `_english` for your language. Try to keep each block of translations in alphabetical order by language name. If there are words between `%percent_signs%`, leave them as-is (they are replaced with numbers by the game client). If there is a colon in the text between `%percent_signs:like this%`, translate only the part after the colon (it is used as a replacement if the number is missing from the item data). Don't change any lines that aren't labelled with a language.

### Changing English strings

If you change an English string in a way that does not require editing other languages (such as fixing a typo that doesn't change the meaning), you need to also change the `[english]` copy of the string in each of the other language files that have it translated. Adding new English strings or changing strings in a way that does require re-translation is handled by the sync script/executable in the utils folder.

### How to test your translation before submitting it

- Go into Steam - Library - Alien Swarm: Reactive Drop, right click and choose Properties. In the Language tab select the language you are translating into.
- Copy your files into respective folders in `C:\Program Files (x86)\Steam\steamapps\common\Alien Swarm Reactive Drop\reactivedrop\resource`

### Got Questions?

Feel free to ask them on this repository's issue tracker or [the Reactive Drop forum](https://steamcommunity.com/app/563560/discussions/1/).
