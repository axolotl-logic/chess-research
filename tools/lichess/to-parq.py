#!/usr/bin/env python
"""
Convert TSV to Parquet format without loading entire file into memory.
Uses chunked processing for memory efficiency with large files.
"""

import argparse
from pathlib import Path
from typing import Optional

import pandas as pd
import pyarrow as pa
import pyarrow.parquet as pq


def tsv_to_parquet(
    input_file: str, output_file: Optional[str] = None, chunk_size: int = 10000
):
    """
    Convert TSV to Parquet format using chunked processing.

    Args:
        input_file: Path to input TSV file
        output_file: Path to output Parquet file (optional)
        chunk_size: Number of rows to process at a time
    """
    input_path = Path(input_file)
    if not input_path.exists():
        raise FileNotFoundError(f"Input file not found: {input_file}")

    # Default output filename
    if output_file is None:
        output_file = str(input_path.with_suffix(".parquet"))

    # Initialize parquet writer
    parquet_writer = None
    schema = None

    try:
        # Process TSV in chunks
        for i, chunk in enumerate(
            pd.read_csv(input_file, sep="\t", chunksize=chunk_size)
        ):
            # Convert to PyArrow table
            table = pa.Table.from_pandas(chunk, preserve_index=False)

            # Initialize writer with schema from first chunk
            if parquet_writer is None:
                schema = table.schema
                parquet_writer = pq.ParquetWriter(output_file, schema)

            # Write chunk to parquet file
            parquet_writer.write_table(table)

            if i == 0:
                print(f"Processing {input_file} -> {output_file}")
                print(
                    f"Schema: {len(table.column_names)} columns, chunk size: {chunk_size}"
                )

        print(f"Conversion complete: {output_file}")

    finally:
        if parquet_writer:
            parquet_writer.close()


def main():
    parser = argparse.ArgumentParser(description="Convert TSV to Parquet format")
    parser.add_argument("input", help="Input TSV file path")
    parser.add_argument("-o", "--output", help="Output Parquet file path")
    parser.add_argument(
        "-c",
        "--chunk-size",
        type=int,
        default=10000,
        help="Chunk size for processing (default: 10000)",
    )

    args = parser.parse_args()

    try:
        tsv_to_parquet(args.input, args.output, args.chunk_size)
    except Exception as e:
        print(f"Error: {e}")
        return 1

    return 0


if __name__ == "__main__":
    exit(main())
