#!/bin/bash

function makefile() {
	linesuffix=""
	filesuffix=""
	if $1; then
		linesuffix="Err"
		filesuffix="_err"
	fi
	2goarray "icon$linesuffix" main < icon.ico > "icon$filesuffix.go"
	sed -i "s/ $//; 5s/^.*$/var icon$linesuffix = []byte{/" "icon$filesuffix.go"
}

makefile false
makefile true
