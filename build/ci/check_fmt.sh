#!/usr/bin/env bash

set -ex

[[ $(gofmt -d . | wc -c) -eq 0 ]]
