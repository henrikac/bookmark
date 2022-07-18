# CLI Bookmark

A bookmarker for your terminal.

## Installation
Run `go install github.com/henrikac/bookmark@latest`.

## Usage
#### Add bookmark
```
$ bookmark add speak echo \"hello world\"
```
This will bookmark `speak` with the command `echo \"hello world\"`.

*Notes*:
- if you want to use quotes you need to escape them.
- if your command contains flags e.g. `-h` you need to add `--` between the bookmark name and the command
```
$ bookmark add du -- docker compose up -d
```

#### List bookmarks
```
$ bookmark list
```
Lists all your saved bookmarks.

#### Search bookmark
```
$ bookmark search <bookmark>
```
Search for a specific `<bookmark>` and return the command if it exists.

#### Execute bookmark
```
$ bookmark exec <bookmark>
```
This will execute the command saved in `<bookmark>`.

#### Remove bookmark
```
$ bookmark remove <bookmark>
```
This will remove `<bookmark>` if it exists.

#### List configurations
```
$ bookmark config list
``'
THis will list all the configurations.

#### Update configuration
```
$ bookmark config set storePath ~/.config/bookmark/
```
This will set the `storePath` to `~/.config/bookmark/bookmarks.json`.
If the last part of the given path is a folder like in the example above then it ***MUST*** end with a `/`.
