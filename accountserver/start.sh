#!/bin/bash
set -ex
make
./accountserver --config=config/config.conf
