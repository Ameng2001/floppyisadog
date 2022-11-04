#!/bin/bash
set -ex
make
./appcommon --config=config/config.conf
