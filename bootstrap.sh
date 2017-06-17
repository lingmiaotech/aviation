#!/bin/bash

export http_proxy=http://localhost:8123
export https_proxy=http://localhost:8123

glide install

echo "Finished. Now go hack!"
