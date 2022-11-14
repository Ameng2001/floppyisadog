#!/bin/bash
set -ex
make
./webportalserver --config=config/config.conf
