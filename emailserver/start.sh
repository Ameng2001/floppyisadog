#!/bin/bash
set -ex
make
./emailserver --config=config/config.conf
