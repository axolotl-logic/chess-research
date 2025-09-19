"""
Spark application to process lichess monthly database dump into
measures of individual moves.
"""

import collections
from typing import Optional

import chess
import networkx as nx
from pyspark.sql import SparkSession
from pyspark.sql.types import (
    FloatType,
    StructType,
    StructField,
    StringType,
    IntegerType,
)


def main():
    spark = SparkSession.builder.appName("IngestLichess").getOrCreate()

    input_path = get_lichess_path()

    df = spark.read.csv(input_path, sep="\t", header=True)

    schema = StructType(
        [
            StructField("move_ply", IntegerType(), True),
            StructField("game_url", StringType(), True),
            StructField("game_tc", StringType(), True),
            StructField("game_tc_base", IntegerType(), True),
            StructField("game_tc_inc", IntegerType(), True),
            StructField("player_elo", IntegerType(), True),
            StructField("move_clock_start", IntegerType(), True),
            StructField("move_clock_end", IntegerType(), True),
            StructField("move_clock_diff", IntegerType(), True),
            StructField("game_utc_date", StringType(), True),
            StructField("game_utc_time", StringType(), True),
            StructField("move_from_fen", StringType(), True),
            StructField("move_from_fragility", FloatType(), True),
            StructField("move_from_fragility_diff", FloatType(), True),
        ]
    )

    df = (
        df.rdd.filter(filter_games)
        .flatMap(measure_game)
        .map(lambda row: list(row[c] for c in schema.fieldNames()))
        .toDF(schema)
    )

    df.write.mode("overwrite").partitionBy("game_tc").parquet(get_output_path())

    spark.stop()


def get_lichess_path():
    return "bronze/volume/lichess_db_standard_rated_2025-07.tsv"


def get_output_path():
    return "bronze/volume/lichess_db_standard_rated_2025-07_moves"


def parse_time_control(tc: str):
    base, inc = tc.split("+")

    return int(base), int(inc)


def filter_games(row):
    event = row["Event"]
    return event.startswith("Rated Classical") and int(row["Ply"]) >= 4


def measure_game(row):
    board = chess.Board()

    moves = row["Moves"].split(" ")
    clocks = {
        chess.WHITE: collections.deque([int(c) for c in row["WhiteClock"].split(" ")]),
        chess.BLACK: collections.deque([int(c) for c in row["BlackClock"].split(" ")]),
    }
    elos = {chess.WHITE: int(row["WhiteElo"]), chess.BLACK: int(row["BlackElo"])}
    base_time, increment = parse_time_control(row["TimeControl"])

    clock_lens_diff = len(clocks[chess.WHITE]) - len(clocks[chess.BLACK])
    assert clock_lens_diff >= 0 and clock_lens_diff <= 1, (
        "Unexpected number of clock entries"
    )

    last_clock: dict[chess.Color, Optional[int]] = {
        chess.WHITE: None,
        chess.BLACK: None,
    }

    last_fragility = None
    last_clock_diff = None
    for idx, move in enumerate(moves):
        start_clock = last_clock[board.turn]
        end_clock = clocks[board.turn].popleft()
        if start_clock is None:
            start_clock = end_clock

        assert start_clock is not None
        assert end_clock is not None

        last_clock[board.turn] = end_clock
        clock_diff = end_clock - start_clock

        move = move.replace("?", "").replace("!", "")

        measures = measure_board(board)
        fragility = measures["move_from_fragility"]
        fragility_diff = None
        if last_fragility is not None:
            fragility_diff = fragility - last_fragility

        last_fragility = fragility

        yield {
            "event": row["Event"],
            "move_ply": idx + 1,
            "game_url": row["Site"],
            "game_tc": row["TimeControl"],
            "game_tc_base": base_time,
            "game_tc_inc": increment,
            "player_elo": elos[board.turn],
            "move_clock_start": start_clock,
            "move_clock_end": end_clock,
            "move_clock_diff": clock_diff,
            "move_from_fragility_diff": fragility_diff,
            "move_from_clock_diff": last_clock_diff,
            "move": move,
            "game_utc_date": row["UTCDate"],
            "game_utc_time": row["UTCTime"],
            **measures,
        }

        board.push_san(move)

        last_clock_diff = clock_diff


def measure_board(board: chess.Board):
    g = nx.DiGraph()
    for square, _ in board.piece_map().items():
        for sees in board.attacks(square):
            g.add_edge(square, sees)

    centrality = nx.betweenness_centrality(g, normalized=True)
    fragility: float = 0
    for sq, value in centrality.items():
        if board.attackers(not board.color_at(sq), sq):
            fragility += value

    return {
        "move_from_fen": board.fen(),
        "move_from_fragility": fragility,
    }


if __name__ == "__main__":
    main()
