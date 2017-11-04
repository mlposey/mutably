#!/bin/bash

REMOTE=marcusposey@srv.marcusposey.com

docker login

docker build -t mlposey/mutably-db -f database/Dockerfile database/
docker push mlposey/mutably-db

docker build -t mlposey/mutably-api -f api/Dockerfile api/
docker push mlposey/mutably-api

docker build -t mlposey/mutably-docs -f api/swagger/Dockerfile api/swagger/
docker push mlposey/mutably-docs

scp docker-compose.yaml $REMOTE:/home/marcusposey/mutably/

ssh $REMOTE << EOF
    cd ~/mutably
    sudo docker-compose down
    sudo docker-compose pull
    sudo -E docker-compose up -d
EOF