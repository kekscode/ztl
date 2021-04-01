#!/bin/sh

DATE=$(date "+%Y%m%d%H%M")

# create proper file
touch "$DATE Test File.md" && sleep 5 && \
echo "# $DATE test file" >"$DATE Test File.md" && sleep 5 &&\
cp "$DATE test file.md" "$DATE test file copied.md" && sleep 5 &&\
mv "$DATE test file copied.md" "$DATE test file renamed.md" && sleep 5 &&\
rm -f "$DATE test file renamed.md" && \
rm -f "$DATE test file copied.md" && \
rm -f "$DATE test file.md"
