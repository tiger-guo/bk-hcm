#!/usr/bin/env bash

cd ../data-service && sh data-service.sh stop
cd ../api-server && sh api-server.sh stop
cd ../auth-server && sh auth-server.sh stop
