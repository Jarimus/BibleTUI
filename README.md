# BibleTUI

![BibleTUI video](assets/BibleTUI.gif)

## Description
- Get different Bible translations from [API.Bible](https://scripture.api.bible/) in a multitude of languages.
- Read them in the command line with an intuitive interface.
- Keep a record of your favourite translations.
- Listen to the Bible in select languages using text-to-speech. (NOTE: Only works on Windows)
- Multiple users - each with their own set of translations.

## Why
Bibles come in many forms: physical books, e-books, audio books, etc.

But command line? With an easy-to-use interface? Welcome to BibleTUI!

## How to run

### (a) Download executable
1.  Go to the repository's [GitHub Actions page](https://github.com/Jarimus/BibleTUI/actions) and choose the '[Build Executable Artifacts](https://github.com/Jarimus/BibleTUI/actions/workflows/Upload_artifact.yml)' action.
2.  Click the latest workflow to access the executables (which are hopefully not expired).
3.  Download the proper executable for your OS and architecture.

### (b) Build from code
1. git clone the repository. On Windows, use the main branch. On Linux/Darwin, use the 'no-tts' branch. The text-to-speech dependencies have deprecated and do not work at the moment for Linux/Darwin.
2. run 'go build .' to build the executable.

### Setting the API key
4.  Go to the [API Bible page](https://scripture.api.bible/)
5.  Create a new account and get a new API key.
6.  Run the executable and enter the API key in the main menu.

## Usage
- Use arrow keys to navigate the menus, 'Enter' to confirm, 'Esc' to go back.
- 'Random verse' produces a random verse from the Bible (World English Bible only).
- 'Read the Bible' lets you choose a book and chapter, and start reading in a scrolling viewport.
- 'Change translation' allows you to change the current translation or add new translations to your personal list from the hundreds available in the API.Bible api.
- You can have multiple users, each with their own list of Bible translations.
- Press 'p' while reading to listen to the chapter in select languages.

## Potential new features
- Better support for different alphabets
- Search function for the longer lists (languages, translations)

## Contributing
All help is greatly appreciated! Contribute by forking the repo and opening pull requests.

All pull requests should be submitted to the main branch.

Did you find any bugs? 🐛 Let me know!
