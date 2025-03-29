# BibleTUI

![BibleTUI video](assets/BibleTUI.gif)

## Description
- Get different Bible translations from API.Bible in a multitude of languages!
- Read them in the command line with an intuitive interface!
- Keep a record of your favourite translations.

## Why
Bibles come in many forms: physical books, e-books, audio books, etc.

But command line? With a easy-to-use interface? Welcome to BibleTUI!

## How to use
1.  Install the [Go toolchain](https://go.dev/dl/).
2.  Go to the [API Bible page](https://scripture.api.bible/)
3.  Create a new account and get a new api key.
4.  Create a 'getApiKey' function in the internal/api_query package. The function should return your api key as a string.
5.  Use 'go build .' to build the executable.
6.  Run your BibleTUI!

## Usage
- Use arrow keys to navigate the menus, 'Enter' to confirm, 'Esc' to go back.
- 'Random verse' produces a random verse from the Bible (World English Bible only).
- 'Read the Bible' lets you choose a book and chapter, and start reading in a scrolling viewport.
- 'Change translation' allows you to change the current translation or add new translations to your personal list from the hundreds available in the API.Bible api.
- While reading use up and down arrows to scroll throught the chapter. Use left and right arrows to navigate between chapters.

## Upcoming features
- Switch between users who can have their own list of translations.
- Implement a proper database instead of simple json files.
- Enter your api-key within the interface.

## Contributing
All help is greatly appreciated! Contribute by forking the repo and opening pull requests.

All pull requests should be submitted to the main branch.