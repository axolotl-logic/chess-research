#!/usr/bin/env bash

volume_path=`dirname "$0"`/volume
base_name=lichess_db_standard_rated_2025-07

remote_file="https://database.lichess.org/standard/$base_name.pgn.zst"
input_file="$volume_path/$base_name.pgn.zst"
output_file="$volume_path/$base_name.tsv"

wget -c $remote_file -O $input_file

zstdcat $input_file  | pgne3k lichess.tagspec  | tqdm > $output_file
