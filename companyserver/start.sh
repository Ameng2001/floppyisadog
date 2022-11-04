#!/bin/bash
set -ex
make
./companyserver --config=config/config.conf
