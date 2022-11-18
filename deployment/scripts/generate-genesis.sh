#!/bin/bash

genesis_car_file_name=$1
genesis_template=$2

cd lotus || exit

./lotus-seed genesis car --out "$genesis_car_file_name" "$genesis_template"
