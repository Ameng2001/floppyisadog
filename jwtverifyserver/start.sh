#!/bin/bash
set -ex
make
./jwtverifyserver --config=config/config.conf
