#!/bin/bash
set -ex
make
./botserver --config=config/config.conf
