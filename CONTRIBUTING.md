Thanks for your interest in the Alien Swarm: Reactive Drop project. When you make a
contribution to the project (e.g. create an Issue or submit a Pull Request)
(a "Contribution"), the Reactive Drop Team wants to be able to use your Contribution to improve
the game.

As a condition of providing a Contribution, you agree that:

- You irrevocably grant anyone the right to use your work under the following license: Creative Commons CC0 Waiver (release all rights, like public domain: [legal code](https://creativecommons.org/publicdomain/zero/1.0/))
- You warrant and represent that the Contribution is your original creation, that you have the authority to grant this license to anyone, and that this license does not require the permission of any third party. Otherwise, you provide your Contribution "as is" without warranties.

# How to translate

There are multiple ways to add or edit your translation of the game:

* via the GitHub web editor [\[jump\]](#using-the-github-web-editor)
* by downloading and editing an archive of the repository [\[jump\]](#using-a-repository-archive)
* by forking the Git repository, editing your fork, and then submitting a pull request [\[jump\]](#using-a-github-fork)

## Using the GitHub web editor

0. Ensure you are [logged in](https://github.com/login) to GitHub.

1. Open the [GitHub web editor](https://github.dev/ReactiveDrop/reactivedrop_translations) by clicking on this link.

2. When prompted, install the suggested extension.

   ![Recommended Extensions tooltip, suggesting the installation of "Valve KeyValue File Support".](https://i.imgur.com/2YjF5ue.png)

3. Use the File Explorer panel on the left to find a file to translate.

   ![The File Explorer activity sidebar.](https://i.imgur.com/bGqR0js.png)

4. Find some untranslated lines – compared to the translated ones, they are **indented with two tabs**!

   ![Two lines to translate, with an orange square emphasising the indentation.](https://i.imgur.com/Dmh3MYi.png)

5. Translate the lines, paying attention to not remove the quotes `"` at the start and at the end of the string.
   Don't edit the lines starting with `[english]`: they're just there to provide some context for the translation!

   ![The previous two lines are now translated into Italian.](https://i.imgur.com/fmD73HV.png)

6. Outdent (remove the indentation from) the now-translated lines.  
   If you have many lines, you can select them all and then press <kbd>Shift</kbd> + <kbd>Tab</kbd>!

   https://user-images.githubusercontent.com/1540885/210732809-08c0ad52-8416-46ac-9bbe-6a90da8e2533.mp4

7. Once you've translated enough lines, open the Source Control panel on the left.

   ![The Source Control activity sidebar.](https://i.imgur.com/6Q0O7DB.png)

8. You will be prompted to enter a message describing your changes; add the name of the language you translated strings to between parenthesis, then describe very briefly what these strings are related to.
   
   For example, a good message might be `(italian) Translate Steam Rich Presence strings`.

   ![The Source Control activity sidebar, now having "(italian) Translate Steam Rich Presence strings" as commit message.](https://i.imgur.com/TMlIhTb.png)

9. Create a new **pull request** by pressing <kbd>Ctrl</kbd> + <kbd>Enter</kbd> while the cursor is in the message box.
   
   If it's your first time submitting translations, you might be prompted to *fork* the repository: please do so!
   
   ![Window prompting to fork the repository, with a orange 9 suggesting to press the "Fork Repository" button.](https://i.imgur.com/P8Vwy3K.png)

   You'll then be prompted to enter the title of the pull request, which will be pre-filled with the message you entered earlier; press <kbd>Enter</kbd> to continue.

   ![Window prompting to enter a pull request title, with "(italian) Translate Steam Rich Presence strings" prefilled.](https://i.imgur.com/MNgRTBU.png)

   You will also be prompted to enter a branch name; leave whatever the editor suggested unchanged, and press <kbd>Enter</kbd> again.

   ![Window prompting to enter a branch name, with "patch-1" prefilled.](https://i.imgur.com/LjubYZn.png)

10. Done! Your translation should now appear in the [list of pull requests of the repository](https://github.com/ReactiveDrop/reactivedrop_translations/pulls), and can be approved, edited, or rejected!

    ![List of the open pull requests of the ReactiveDrop repository, displaying 1 open with the title "(italian) Translate Steam Rich Presence strings".](https://i.imgur.com/7XNWLNo.png)

## Using a repository archive

[Download](https://github.com/ReactiveDrop/reactivedrop_translations/archive/refs/heads/master.zip) the repository, edit it on your computer and then submit it through the [Issues tab](https://github.com/ReactiveDrop/reactivedrop_translations/issues) by providing an archive.

## Using a GitHub fork

[Fork the repository](https://github.com/ReactiveDrop/reactivedrop_translations/fork), edit & commit it and then push your changes back to the origin.

Use an editor to make your changes and integrate them through github. Read more about it [on our wiki](https://github.com/ReactiveDrop/reactivedrop_translations/wiki).

### Pre-requisites
- Any unicode plain text editor. E.g. [Notepad++](https://notepad-plus-plus.org/) or [Visual Studio Code](https://code.visualstudio.com/).
- An [GIT installation](https://git-scm.com/downloads)

# Files

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
6. Start game, and test your translation :D

- ATTENTION: Every file in the resource folder is UTF-16LE encoded, whereas all files in this repo are UTF-8 encoded. You may need to change the encoding of the files before you copy them over.
