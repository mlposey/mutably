#!/bin/bash
# Performs a benchmark of $LIMIT page insertions.
#
# You'll want to set the values of DB, USER, and PASSWORD
# in this file before using.
#
# The result of the time command is output to a text file
# that has the name of your current git branch. Profile
# information is output to the default destination
# specified by the anvil source code.

LIMIT=100000

# DB=your-db
# USER=your-db-user
# PASSWORD=your-db-user-password

if [[ -z ${DB+x} || -z ${USER+x} || -z ${PASSWORD+x} ]]; then
    echo "Modify the database connection parameters in this file."
    exit 1
fi

CURRENT_BRANCH=$(git branch | grep '[*] .*$' | cut -c3-)

# We don't know if it's empty when we start, but we will clear
# verbs when we're done.
psql -U $USER -c "delete from verbs;"

go build
/usr/bin/time -o time_${CURRENT_BRANCH}.txt \
  ./anvil -import -profile -limit=$LIMIT -d=$DB -u=$USER -p=$PASSWORD \
  ~/data/wiktionary/enwiktionary_latest/enwiktionary-latest-pages-meta-current.xml

psql -U $USER -c "delete from verbs;"
