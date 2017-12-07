# Archived

This thing is useless compared to [rsync](https://en.wikipedia.org/wiki/Rsync) and friends.

-------------------------------

# gofilesync

**gofilesync** is a complete rewrite of [FolderSync](https://legolord208.github.io/software#foldersync) (sort of accidentally stolen name from an android app, I know)  
gofilesync is both a *command line* tool, a *GUI*, AND a library!  
If you start it normally, it simply puts itself in your system tray, and does nothing...  
*until* you click "Configure", and the menu pops up.

Set a schedule, add folders, et.c  
It all just works flawlessly!

## "Lazy" sync

The auto-scheduled sync only syncs modified files!

## Command line

`gofilesync --help`  
```
Usage of gofilesync:
  -dst string
    	The destination folder to paste.
  -lazy
    	Whether or not to only sync necessary files.
  -src string
    	The source folder to copy.
```

### Example
```
gofilesync --src folder1 --dst folder2 --lazy
```

## API

Heck that's right! As if the command line tools wasn't enough,  
you can also use the `api` folder  
to make your completely custom sync application in Go!

And don't worry, gofilesync automatically makes sure you don't try to sync the same thing twice at the same time.

### Example

```Go
err := gofilesync.LazySync("folder1", "folder2")
if err != nil {
	// Handle errors
}
```
