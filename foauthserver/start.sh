#!/bin/bash
set -ex
make
./foauthserver --config=config/config.conf
