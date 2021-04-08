#!/bin/sh

DATE=$(date "+%Y%m%d%H%M")

n=6; 
NAME1=$(gtr -cd '[a-zA-Z0-9]' < /dev/urandom | head -c$n); echo $output
NAME2=$(gtr -cd '[a-zA-Z0-9]' < /dev/urandom | head -c$n); echo $output
NAME3=$(gtr -cd '[a-zA-Z0-9]' < /dev/urandom | head -c$n); echo $output

NAME="$NAME1 $NAME2 $NAME3"

TIMEOUT=1

# create proper file
touch "$DATE $NAME.md" && sleep $TIMEOUT && \
echo "# $DATE new markdown header" >"$DATE $NAME.md" && sleep $TIMEOUT &&\
cp "$DATE new markdown header.md" "$DATE new markdown header copied.md" && sleep $TIMEOUT &&\
mv "$DATE new markdown header copied.md" "$DATE new markdown header renamed.md" && sleep $TIMEOUT &&\
mv "$DATE new markdown header renamed.md" "$DATE new markdown header renamed a.g.ain.md" && sleep $TIMEOUT &&\
rm -f "$DATE new markdown header renamed.md" && \
rm -f "$DATE new markdown header renamed a.g.ain.md" && \
rm -f "$DATE new markdown header copied.md" && \
rm -f "$DATE new markdown header.md"

