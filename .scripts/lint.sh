#!/usr/bin/env bash

# Copyright 2017 Typeform SL. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

set -e

for d in $(go list ./... | grep -v vendor); do
  gometalinter $GOPATH/src/$d
done
