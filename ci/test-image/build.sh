#!/bin/bash

git -C kitt pull || git clone git@github.com:figome/kitt.git

docker build -t figo/smcl-test .