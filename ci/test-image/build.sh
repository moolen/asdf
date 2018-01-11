#!/bin/bash

git -C kitt pull || git clone git@github.com:figome/kitt.git

docker build -t eu.gcr.io/figo-v1/smcl-test .
docker push eu.gcr.io/figo-v1/smcl-test