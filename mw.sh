#!/bin/bash

# Check if file is provided
if [ $# -eq 0 ]; then
  echo "Usage: mw <word>"
  exit 1
fi

curl "https://www.merriam-webster.com/dictionary/$1" \
| grep -o '<span class="dtText">.*</span>' \
| sed 's/.*<\/strong>\(.*\)<.*/-- \1/' \
| sed 's/.*<span class=.text-uppercase.>\([^<]*\)<\/span>.*/-\* \1/'

