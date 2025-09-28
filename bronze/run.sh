#!/usr/bin/env bash

set -e

volume_path=`dirname "$0"`/volume
date_suffix=`date --date="$(date +%Y-%m-15) -1 month" +%Y-%m`
base_name=lichess_db_standard_rated_$date_suffix

remote_file="https://database.lichess.org/standard/$base_name.pgn.zst"
input_file="$volume_path/$base_name.pgn.zst"
output_file="$volume_path/$base_name.tsv"
output_tmp_file="$output_file.tmp"

wget -c $remote_file -O $input_file

zstdcat $input_file  | pgne3k lichess.tagspec  | tqdm > $output_tmp_file
mv $output_tmp_file $output_file
