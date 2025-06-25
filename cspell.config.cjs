"use strict";

const fs = require("fs");
const path = require("path");

const minWordLength = 4;
const wordFiles = ["go.mod", "go.sum"];

function loadFile(cwd, filename) {
  const filepath = path.join(cwd, filename);
  if (!fs.existsSync(filepath)) {
    return "";
  }

  return fs.readFileSync(filepath, { encoding: "utf8", flag: "r" });
}

function customWords(cwd = process.cwd()) {
  let result = [];

  wordFiles.forEach((filename) => {
    const fileData = loadFile(cwd, filename);

    const words = fileData
      .match(/(\w+)/g)
      .map((w) => w.trim())
      .filter((w) => !!w && w.length >= minWordLength);

    result = result.concat(words);
  });

  return [...new Set(result)];
}

function expand(
  pattern,
  options = { begin: "(", end: ")", sep: "|" },
  start = 0
) {
  const len = pattern.length;
  const parts = [];
  function push(word) {
    if (Array.isArray(word)) {
      parts.push(...word);
    } else {
      parts.push(word);
    }
  }
  let i = start;
  let curWord = "";
  while (i < len) {
    const ch = pattern[i++];
    if (ch === options.end) {
      break;
    }
    if (ch === options.begin) {
      const nested = expand(pattern, options, i);
      i = nested.idx;
      curWord = nested.parts.flatMap((p) =>
        Array.isArray(curWord) ? curWord.map((w) => w + p) : [curWord + p]
      );
      continue;
    }
    if (ch === options.sep) {
      push(curWord);
      curWord = "";
      continue;
    }
    curWord = Array.isArray(curWord)
      ? curWord.map((w) => w + ch)
      : curWord + ch;
  }
  push(curWord);
  return { parts, idx: i };
}

function expandWord(pattern, options = { begin: "(", end: ")", sep: "|" }) {
  return expand(pattern, options).parts;
}

function expandWords(wordList) {
  const words = wordList.flatMap((w) => expandWord(w));
  return words;
}

/** @type { import("@cspell/cspell-types").CSpellUserSettings } */
const config = {
  version: "0.2",
  language: "en,en-gb",
  caseSensitive: true,
  minWordLength: minWordLength,
  userWords: [],
  useGitignore: true,
  spellCheckOnlyWorkspaceFiles: true,
  ignorePaths: [
    ".cspell.json",
    ".golangci.yml",
    ".vscode",
    ".git",
    "go.mod",
    "go.sum",
  ],
  ignoreRegExpList: [
    // ignore urls
    "Urls",
    // ignore urls
    "HexValues",
    // ignore urls
    "Base64",
    // ignore email addresses
    "Email",
  ],
  languageSettings: [
    {
      languageId: "go",
      caseSensitive: false,
      allowCompoundWords: true,
      ignoreRegExpList: [
        // ignore multiline imports
        /import\s*\((.|(?:\r)?\n)*?\)/g,
        // ignore single line imports
        /import\s*.*".*?"/g,
        // ignore go generate directive
        /\/\/\s*go:generate.*/g,
        // ignore nolint directive
        /\/\/\s*nolint:.*/g,
        // ignore logging methods
        "(print|debug|info|warn|error|panic|fatal)(ln|f|j)",
      ],
    },
    {
      languageId: "makefile",
      caseSensitive: false,
      allowCompoundWords: true,
      dictionaries: ["go", "makefile"],
    },
    {
      languageId: "javascript",
      caseSensitive: true,
      allowCompoundWords: false,
    },
  ],
  enabledLanguageIds: [
    "bash",
    "go",
    "javascript",
    "json",
    "jsonc",
    "makefile",
    "markdown",
    "yaml",
    "yml",
  ],
  dictionaryDefinitions: [
    {
      name: "custom-words",
      words: customWords(),
    },
  ],
  dictionaries: ["custom-words", "go", "makefile"],
  words: expandWords([
    "golangci",
    "requestid",
    "prefixer(s|)",
    "fsync",
    "nolint",
    "libstd(c|)",
    "(ext|)ldflags",
    "txlock",
    "uow",
    "algorand",
  ]),
};

module.exports = config;
