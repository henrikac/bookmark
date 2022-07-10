# CLI Bookmark

## Requirements
- [jq](https://stedolan.github.io/jq/)

## Usage
Bookmarks are stored in `$HOME/.bookmarks.json`.

#### Save bookmark
```
$ bm save syu sudo pacman -S jq
```
This will bookmark `syu` with the command `sudo pacman -S jq`.
To bookmark the most recent run command simply run
```
$ bm save <bookmark> !!
```
or if you want to bookmark a command that you ran earlier
```
$ history
...
844 docker compose up -d
...
$ bm save <bookmark> !844
```

#### List bookmarks
```
$ bm list
```
Lists all your saved bookmarks.

#### Search bookmark
```
$ bm search <bookmark>
```
Search for a specific `<bookmark>` and return the command if it exists.

#### Execute bookmark
```
$ bm exec <bookmark>
```
This will execute the command saved in `<bookmark>`.

#### Remove bookmark
```
$ bm remove <bookmark>
```
This will remove `<bookmark>` if it exists.

#### Clear bookmarks
```
$ bm clear
```
This will delete *ALL* your bookmarks including the `.bookmarks.json` file.
