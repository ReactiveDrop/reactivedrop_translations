# Alien Swarm: Reactive Drop Translations #

This repository contains all the files required to translate [Alien Swarm: Reactive Drop](http://store.steampowered.com/app/563560/) user interface into any supported language. Reactive Drop(RD) supports all languages which are supported by Steam. You can also translate to an unsupported language but this will only work as a VPK add-on file through workshop.

### Pre-requisites ###

* Any plain text editor. E.g. [Notepad++](https://notepad-plus-plus.org/)
* Fork the repository, add your changes in your own fork, and then submit a pull request to merge in your changes.
* If you have problems using forks or GIT you can submit your translation through the Issues tab by providing an archive.

### Directory structure ###
* achievements folder contains files for translating achievements which will be shown on Steam websites
* resource folder contains game localization files. Same as Alien Swarm Reactive Drop\reactivedrop\resource
* store_page is [Steam's store page](http://store.steampowered.com/app/563560/)

### How to download the repository ###
If you are not familiar with using GIT you can [download](https://bitbucket.org/reactivedropteam/reactivedrop_translations/downloads/) the repository as a ZIP archive and submit it using [Issues](https://bitbucket.org/reactivedropteam/reactivedrop_translations/issues) tab.

### CONTRIBUTING ###
Thanks for your interest in the Alien Swarm: Reactive Drop project.  When you make a
contribution to the project (e.g. create an Issue or submit a Pull Request)
(a "Contribution"), Reactive Drop Team wants to be able to use your Contribution to improve
the game.

As a condition of providing a Contribution, you agree that: 

* You irrevocably grant anyone the right to use your work under the following license: Creative Commons CC0 Waiver (release all rights, like public domain: [legal code](https://creativecommons.org/publicdomain/zero/1.0/))
* You warrant and represent that the Contribution is your original creation,
that you have the authority to grant this license to anyone, and that this
license does not require the permission of any third party.  Otherwise, you
provide your Contribution "as is" without warranties. 

### How to translate? ###
### Achievements
Create a copy of file 563560_loc_english.vdf and rename it to your language(e.g. 563560_loc_ukrainian.vdf).
Translate the contents of this file similarly to how 563560_loc_ukrainian.vdf is translated.
> **Tip:** Most of achievements are same with Alien Swarm and they were already translated. You can find the translated achievements into for example reactivedrop_ukrainian.txt starting at ASW_KILL_WITHOUT_FRIENDLY_FIRE_NAME.
### Mail and News
Create a copy of each mail and news file and rename it to your language sufix(e.g. labsmail1_russian.txt). Translate the contents of each file. See labsmail1_russian.txt as an example.
### BaseModUi, CloseCaption, GameUi, ReactiveDrop
In these files(e.g. basemodui_czech.txt) the untranslated strings are indented by two tabs to be easily visible. Translate each indented line and remove the two tabs in front of it.
### How to test your translation before submitting it ###
* Set your Steam's interface language to the language you are translating into
* Go into Steam - Library - Alien Swarm: Reactive Drop, right click and choose Properties. In the Language tab select the language you are translating into. 
> **Note:** Changing Steam's interface langue will affect game's interface language. Changing language in game's properties tab will affect voice over sounds, subtitles, mail and news language.

* Copy your files into respective folders in C:\Program Files (x86)\Steam\steamapps\common\Alien Swarm Reactive Drop
### Got Questions? ###
Feel free to ask them [here](http://steamcommunity.com/app/563560/discussions/1/). Prefix your question with [Translation] tag. 