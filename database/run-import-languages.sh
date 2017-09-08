#!/bin/bash
# Run the import languages python script.
# The official image for PostgreSQL automatically runs
# .sql and .sh files but not .py. That design makes
# this file necessary.
python3 /docker-entrypoint-initdb.d/import-languages.py /data/language-subtag-registry
