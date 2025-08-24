import argparse
from enum import Enum
import json
import logging
from pathlib import Path
from typing import NamedTuple

import apache_beam as beam
from apache_beam.io import (
    FileBasedSink,
    ReadFromText,
    filesystems,
)
from apache_beam.io.fileio import FileMetadata, WriteToFiles, destination_prefix_naming

from apache_beam.io.filesystem import CompressionTypes
from apache_beam.options.pipeline_options import PipelineOptions, SetupOptions


class EventType(Enum):
    RATED_RAPID_GAME = "rapid"
    RATED_CLASSICAL_GAME = "classical"
    OTHER = "other"


class ResultType(Enum):
    DRAW = "1/2-1/2"
    WIN = "1-0"
    LOSS = "0-1"
    INCOMPLETE = "*"


class Game(NamedTuple):
    ply: int
    event: str
    white: str
    black: str
    result: str
    white_elo: int
    black_elo: int
    moves: str
    black_clock: str
    white_clock: str
    termination: str
    date: str


class GameTsvSink(FileBasedSink):
    def __init__(self, **kwargs):
        super().__init__(**kwargs, coder=None)
        self._header_written = False

    def open(self, temp_path):
        """Open the file and create CSV writer."""
        self._file_handle = temp_path
        self._header_written = False
        return self._file_handle

    def flush(self):
        self._file_handle.flush()

    def write(self, value):
        """Write header once, then write element's chess game fields."""
        # Write header row once
        if not self._header_written:
            header_line = (
                "\t".join(
                    [
                        "ply",
                        "event",
                        "white",
                        "black",
                        "result",
                        "white_elo",
                        "black_elo",
                        "moves",
                        "black_clock",
                        "white_clock",
                        "termination",
                        "date",
                    ]
                )
                + "\n"
            )
            self._file_handle.write(header_line.encode("utf-8"))
            self._header_written = True

        # Write element data
        data_line = (
            "\t".join(
                [
                    str(value.ply),
                    str(value.event),
                    str(value.white),
                    str(value.black),
                    str(value.result),
                    str(value.white_elo),
                    str(value.black_elo),
                    str(value.moves),
                    str(value.black_clock),
                    str(value.white_clock),
                    str(value.termination),
                    str(value.date),
                ]
            )
            + "\n"
        )
        self._file_handle.write(data_line.encode("utf-8"))

    def close(self, file_handle):
        """Close the file handle."""
        file_handle.close()

    def create_metadata(self, destination, full_file_name):
        """Return metadata for the sink. Required by FileBasedSink."""
        return FileMetadata(
            mime_type="text/csv",
            compression_type=CompressionTypes.UNCOMPRESSED,
        )


def game_from_line(line, indexes) -> list[Game]:
    cols = line.split("\t")
    if cols[indexes["Ply"]] == "Ply":
        return []

    event_raw = cols[indexes["Event"]]
    event = EventType.OTHER
    if event_raw == "Rated Rapid game":
        event = EventType.RATED_RAPID_GAME
    elif event_raw == "Rated Classical game":
        event = EventType.RATED_CLASSICAL_GAME

    result = ResultType(cols[indexes["Result"]])

    return [
        Game(
            ply=int(cols[indexes["Ply"]]),
            event=event.value,
            white=cols[indexes["White"]],
            black=cols[indexes["Black"]],
            result=result.value,
            white_elo=int(cols[indexes["WhiteElo"]]),
            black_elo=int(cols[indexes["BlackElo"]]),
            moves=cols[indexes["Moves"]],
            black_clock=cols[indexes["BlackClock"]],
            white_clock=cols[indexes["WhiteClock"]],
            termination=cols[indexes["Termination"]],
            date=cols[indexes["UTCDate"]],
        )
    ]


def filter_game(game: Game):
    return game.ply > 3 and game.event != EventType.OTHER.value


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--input",
        dest="input",
        required=True,
        help="Input file to process.",
    )
    parser.add_argument(
        "--output",
        dest="output",
        required=True,
        help="Output file to write results to.",
    )

    args, pipeline_args = parser.parse_known_args()

    output_dir = Path(args.output)

    pipeline_options = PipelineOptions(pipeline_args)
    pipeline_options.view_as(SetupOptions).save_main_session = True

    with open("tools/config.json") as f:
        source_config = json.load(f)

    headers = source_config["sources"]["lichess-structured"]["headers"]
    indexes = {k: v["idx"] for k, v in headers.items()}

    with beam.Pipeline(options=pipeline_options) as p:
        _ = (
            p
            | ReadFromText(args.input)
            | "ToGame" >> beam.FlatMap(game_from_line, indexes)
            | "FilterGames" >> beam.Filter(filter_game)
            | "WriteGames"
            >> WriteToFiles(
                path=output_dir,
                destination=lambda g: f"{g.event} - {g.date}",
                sink=lambda _: GameTsvSink(file_path_prefix="lichess - "),
                file_naming=destination_prefix_naming(),
            )
        )


if __name__ == "__main__":
    logging.getLogger().setLevel(logging.INFO)
    main()
