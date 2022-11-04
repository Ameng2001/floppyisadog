#!/bin/bash
set -ex
make
./smsserver --config=config/config.conf
