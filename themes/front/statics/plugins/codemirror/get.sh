#!/bin/bash

[ -f codemirror.js ] && exit 0

curl -L http://codemirror.net/codemirror.zip > codemirror.zip
unzip codemirror.zip
cd codemirror-*

cp lib/codemirror.css ..
cp lib/codemirror.js ..

cp mode/xml/xml.js ..
cp mode/css/css.js ..
cp mode/javascript/javascript.js ..
cp mode/htmlmixed/htmlmixed.js ..

cp mode/markdown/markdown.js ..

cp addon/dialog/dialog.css ..
cp addon/dialog/dialog.js ..
cp addon/display/fullscreen.css ..
cp addon/display/fullscreen.js ..

cp keymap/vim.js ..

cp LICENSE ..

cd ..
rm codemirror.zip
rm -rf codemirror-*

