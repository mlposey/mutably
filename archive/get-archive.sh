#!/bin/bash
# The docker-compose file requires a pages dump in this directory.
# This script should be run before bringing up the containers.
wget https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-meta-current.xml.bz2
lbzip2 -d -k enwiktionary-latest-pages-meta-current.xml.bz2
