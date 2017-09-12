#!/bin/bash
# This script helps TeamCity identify issues in the database definition.
# It should be run from the root directory of the project like:
# './database/test_database_deployment.sh'

# The number of seconds to wait before checking the container status
# We shouldn't expect to wait long. Extended setup routines are best
# left to other containers that would interact with this one. Just
# establish the schema and be done with it.
TIMEOUT=10

docker-compose up --build -d database
sleep $TIMEOUT
docker-compose ps | grep 'mutably_database.* Up'
docker-compose down -v
