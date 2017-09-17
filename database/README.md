The collection of documents in this path comprise the application database. It
uses PostgreSQL 9.6 in an environment defined in `Dockerfile`. Other parts of
the system, like anvil and the docker-compose setup, may expect an XML archive
in `./data`. The archive can be retrieved after installing `wget` and `lbzip2`,
followed by running the bash script `get-archive.sh`.
