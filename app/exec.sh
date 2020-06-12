#!/usr/bin/env bash
docker build -t guard-my-app .
docker run --init -p 3000:3000 -it guard-my-app