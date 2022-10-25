Thanks for your interest in the Alien Swarm: Reactive Drop project. When you make a
contribution to the project (e.g. create an Issue or submit a Pull Request)
(a "Contribution"), the Reactive Drop Team wants to be able to use your Contribution to improve
the game.

As a condition of providing a Contribution, you agree that:

- You irrevocably grant anyone the right to use your work under the following license: Creative Commons CC0 Waiver (release all rights, like public domain: [legal code](https://creativecommons.org/publicdomain/zero/1.0/))
- You warrant and represent that the Contribution is your original creation, that you have the authority to grant this license to anyone, and that this license does not require the permission of any third party. Otherwise, you provide your Contribution "as is" without warranties.


# Ways to translate
## Which way
There are multiple ways to add or edit your translation of the game:
1. The easy way
  * Edit the repository directly through GitHub's web editor. Navigate to a file, click on the edit-button on the top far right, make your changes, add a description click on the 'Propose changes'.
  * Pro: fast and easy.
  * Con: not suited for bigger changes, as the web-editor has much less tools than a full fledged text editor.
2. The manual way
  * [Download](https://github.com/ReactiveDrop/reactivedrop_translations/archive/refs/heads/master.zip) the repository, edit it on your computer and then submit it through the [Issues tab](https://github.com/ReactiveDrop/reactivedrop_translations/issues) by providing an archive.
  * Pro: you do not have to understand how to work with [Git](https://en.wikipedia.org/wiki/Git).
  * Con: the admins have to do the git-work for you, and as a result your changes will be logged as theirs.
3. The professional way
  * [Fork the repository](https://github.com/ReactiveDrop/reactivedrop_translations/fork), edit & commit it and then push your changes back to the origin.
  * Pro: full control over big sets of changes, use functions like search & replace in files, regex, machine translation.
  * Con: steeper learning curve using the git system.
  * Use an editor to make your changes and integrate them through github. Read more about it [on our wiki](https://github.com/ReactiveDrop/reactivedrop_translations/wiki).

## Pre-requisites
- Any unicode plain text editor. E.g. [Notepad++](https://notepad-plus-plus.org/) or [Visual Studio Code](https://code.visualstudio.com/).
- An [GIT installation](https://git-scm.com/downloads)

## File Encoding
Please make sure that your text editor preserves file's encoding.  
The new files you create need to have same encoding as their english counterpart (e.g. labsmail1_czech.txt needs to be UTF-8 encoded just like labsmail1_english.txt). If you use Notepad++ you can see the file's encoding in the menu Encoding.

## File type
### Achievements
These files are automatically created based on the `resource/reactivedrop_*.txt` files. You can manually trigger an update of the files through executing `utils/translation-sync-tool.exe`.
* sync-tool: Very useful verification tool! 
 1. It can help solving issues with ci-checks. After you have manually modified the content about achievements, it is highly recommended to run the tool to check for problems. Avoid Pull Request Failure.
 2. It can automatically synchronize achievement fields from `reactivedrop_*.txt` to `563560_loc_*.vdf`! This is very useful to avoid the painful achievement of comparing and duplicating work.

### Mail and News
Create a copy of each mail and news file, and rename it to your language sufix(e.g. labsmail1_russian.txt). Translate the contents of each file. See labsmail1_russian.txt as an example.

### BaseModUi, CloseCaption, GameUi, ReactiveDrop, and Community VDF files
In these files(e.g. basemodui_czech.txt) the untranslated strings are indented by two tabs to be easily visible. Translate each indented line and remove the two tabs in front of it.

### Workshop
Create a copy of the English file and rename it to your language suffix (e.g. workshop_tags_schinese.json). Translate the contents of the file. In JSON files, only translate text on the right side of the `:`.

### Item Schema
Copy the lines with `_english` for your language. Try to keep each block of translations in alphabetical order by language name. If there are words between `%percent_signs%`, leave them as-is (they are replaced with numbers by the game client). If there is a colon in the text between `%percent_signs:like this%`, translate only the part after the colon (it is used as a replacement if the number is missing from the item data). Don't change any lines that aren't labelled with a language.

## Changing English strings
If you change an English string in a way that does not require editing other languages (such as fixing a typo that doesn't change the meaning), you need to also change the `[english]` copy of the string in each of the other language files that have it translated. Adding new English strings or changing strings in a way that does require re-translation is handled by the sync script/executable in the utils folder.

## How to test your translation before submitting it
First of all, Go into Steam - Library - Alien Swarm: Reactive Drop, right click and choose 'Properties'.
1. Select 'Language' tab: Select the language you are translating into.
2. Select 'Betas' (Recommended): Join the beta version, this way you can see the real effect of the translated text in the game. Otherwise, the game text may not match the translated text.
3. Select 'Local files': Click 'Browse...' to open your game folder.
4. Enter the subfolders '.\reactivedrop\resource', like: `Steam\steamapps\common\Alien Swarm Reactive Drop\reactivedrop\resource`
5. Copy your translate files into the respective folders (.\resource).
6. Start game, and test your translation  :D

- ATTENTION: Every file in the resource folder is UTF-16LE encoded, whereas all files in this repo are UTF-8 encoded. You may need to change the encoding of the files before you copy them over.
