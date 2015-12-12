# piken #
`piken` is a CLI tool to search unicode data backed by SQLite3.

The name is in homage to the creators of UTF-8: Rob **Pike** and **Ken**
Thompson.

## Install ##

```console
% go get eriol.xyz/piken
```

After installation use `piken update` to get unicode data from
http://www.unicode.org/Public/UNIDATA/UnicodeData.txt.

## Usage ##

```console
% piken -h
NAME:
   piken - unicode search tool backed by SQLite3

USAGE:
   ./piken [global options] command [command options] [arguments...]

VERSION:
   0.1

AUTHOR(S):
   Daniele Tricoli <eriol@mornie.org>

COMMANDS:
   update       Update unicode data
   search, s    Search for unicode
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```

## Search ##

`piken` exposes [SQLite3's full-text search](https://www.sqlite.org/fts3.html)
on command line.

**NOTE:** if you are using *zsh*, it's useful to disable glob expressions with
`noglob` (using an alias for example) instead of quote searches.

Default fields showed are *codepoint* and *name*, but it easy to override this:
```console
% piken search red apple
1F34E -- RED APPLE -- 🍎

piken search --fields Name --fields Category --fields CodePoint red apple
RED APPLE -- So -- 1F34E -- 🍎
```

The following fields can be used:
- CodePoint
- Name
- Category
- CanonicalClass
- BidiClass
- DecompositionType
- NumericType
- NumericDigit
- NumericValue
- BidiMirrored
- Unicode1Name
- ISOComment
- SimpleUppercaseMapping
- SimpleLowercaseMapping
- SimpleTitlecaseMapping


### AND, OR, NOT operators ###

```console
% piken search heart AND (blue OR green)
1F499 -- BLUE HEART -- 💙
1F49A -- GREEN HEART -- 💚
```

**Note:** operators must be uppercase.

### * and ^ prefixes ###

```console
% piken search mass*
1D340 -- TETRAGRAM FOR MASSING -- 𝍀
1F486 -- FACE MASSAGE -- 💆
```

```console
% piken search ^cat
1F408 -- CAT -- 🐈
1F431 -- CAT FACE -- 🐱
1F639 -- CAT FACE WITH TEARS OF JOY -- 😹
1F63C -- CAT FACE WITH WRY SMILE -- 😼
```

### NEAR operator ###

```console
% piken search crying face
1F622 -- CRYING FACE -- 😢
1F62D -- LOUDLY CRYING FACE -- 😭
1F63F -- CRYING CAT FACE -- 😿

% piken search crying NEAR/0 face
1F622 -- CRYING FACE -- 😢
1F62D -- LOUDLY CRYING FACE -- 😭
```

### Specify column used to search ###

```console
% piken search codepoint:1F602
1F602 -- FACE WITH TEARS OF JOY -- 😂

% piken search doo* codepoint:A*
A585 -- VAI SYLLABLE DOO -- ꖅ
A61B -- VAI SYMBOL DOONG -- ꘛ

% piken search category:Ll | wc -l
1984

% piken search math* frak* cap* NOT bold codepoint:1D50*
1D504 -- MATHEMATICAL FRAKTUR CAPITAL A -- 𝔄
1D505 -- MATHEMATICAL FRAKTUR CAPITAL B -- 𝔅
1D507 -- MATHEMATICAL FRAKTUR CAPITAL D -- 𝔇
1D508 -- MATHEMATICAL FRAKTUR CAPITAL E -- 𝔈
1D509 -- MATHEMATICAL FRAKTUR CAPITAL F -- 𝔉
1D50A -- MATHEMATICAL FRAKTUR CAPITAL G -- 𝔊
1D50D -- MATHEMATICAL FRAKTUR CAPITAL J -- 𝔍
1D50E -- MATHEMATICAL FRAKTUR CAPITAL K -- 𝔎
1D50F -- MATHEMATICAL FRAKTUR CAPITAL L -- 𝔏

# Note: Capital C and Capital H are missing from UnicodeData.txt
```

### Copy into clipboard ###

Option `--copy` (compact form `-c`) can be used to copy the glyph directly into
the clipboard when search return only one row:

```console
% piken search --copy door
WARN[0000] Copy to clipboard not allowed for multiple rows.
2F3E -- KANGXI RADICAL DOOR -- ⼾
1F6AA -- DOOR -- 🚪

% piken search --copy door NOT k*
1F6AA -- DOOR -- 🚪
```
