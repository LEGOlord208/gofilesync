#!/bin/bash

function fixfile() {
	linesuffix=""
	filesuffix=""
	if [ $1 ]; then
		linesuffix="Err"
		filesuffix="_err"
	fi
	sed -i "s/ $//; 5s/^.*$/var icon$linesuffix = []byte{/" "icon$filesuffix.go"
}

2goarray icon main < icon.ico > icon.go
fixfile false

2goarray iconErr main < icon_err.ico > icon_err.go
fixfile true
