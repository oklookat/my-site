#!/bin/sh

npm update -g npm
npm install

exec "$@"