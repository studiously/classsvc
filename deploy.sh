#!/bin/bash

export TAG=`if [ "$TRAVIS_PULL_REQUEST"=="false" ]; then echo "latest"; else echo $TRAVIS_PULL_REQUEST_BRANCH ; fi`

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
docker build -f Dockerfile -t $REPO:$COMMIT .
docker tag $REPO:$COMMIT $REPO:$TAG
docker tag $REPO:$COMMIT $REPO:travis-$TRAVIS_BUILD_NUMBER

docker push $REPO