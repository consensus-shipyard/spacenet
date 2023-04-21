# This script generates a membership ("validator set") understood by Eudico
# from validator addresses in the format output by
#
# eudico net listen
#
# Input is read from stdin and output written to stdout.
# Configuration number of the membership is always set to 0.
#
# Example input:
# t1dgw4345grpw53zdhu75dc6jj4qhrh4zoyrtq6di@/ip4/172.31.39.78/tcp/43077/p2p/12D3KooWNzTunrQtcoo4SLWNdQ4EdFWSZtah6mgU44Q5XWM61aan
# t1a5gxsoogaofa5nzfdh66l6uynx4m6m4fiqvcx6y@/ip4/172.31.33.169/tcp/38257/p2p/12D3KooWABvxn3CHjz9r5TYGXGDqm8549VEuAyFpbkH8xWkNLSmr
# t1q4j6esoqvfckm7zgqfjynuytjanbhirnbwfrsty@/ip4/172.31.42.15/tcp/44407/p2p/12D3KooWGdQGu1utYP6KD1Cq4iXTLV6hbZa8yQN34zwuHNP5YbCi
# t16biatgyushsfcidabfy2lm5wo22ppe6r7ddir6y@/ip4/172.31.47.117/tcp/34355/p2p/12D3KooWEtfTyoWW7pFLsErAb6jPiQQCC3y3junHtLn9jYnFHei8
#
# Example output:
# {
#     "configuration_number": 0,
#     "validators": [
#         {
#             "addr": "t1dgw4345grpw53zdhu75dc6jj4qhrh4zoyrtq6di",
#             "net_addr": "/ip4/172.31.39.78/tcp/43077/p2p/12D3KooWNzTunrQtcoo4SLWNdQ4EdFWSZtah6mgU44Q5XWM61aan",
#             "weight": "0"
#         },
#         {
#             "addr": "t1a5gxsoogaofa5nzfdh66l6uynx4m6m4fiqvcx6y",
#             "net_addr": "/ip4/172.31.33.169/tcp/38257/p2p/12D3KooWABvxn3CHjz9r5TYGXGDqm8549VEuAyFpbkH8xWkNLSmr",
#             "weight": "0"
#         },
#         {
#             "addr": "t1q4j6esoqvfckm7zgqfjynuytjanbhirnbwfrsty",
#             "net_addr": "/ip4/172.31.42.15/tcp/44407/p2p/12D3KooWGdQGu1utYP6KD1Cq4iXTLV6hbZa8yQN34zwuHNP5YbCi",
#             "weight": "0"
#         },
#         {
#             "addr": "t16biatgyushsfcidabfy2lm5wo22ppe6r7ddir6y",
#             "net_addr": "/ip4/172.31.47.117/tcp/34355/p2p/12D3KooWEtfTyoWW7pFLsErAb6jPiQQCC3y3junHtLn9jYnFHei8",
#             "weight": "0"
#         }
#     ]
# }

import sys
import json

membership = {
    "configuration_number": 0,
    "validators": [],
}

def parse_validator(line: str):
    tokens = line.split("@")
    membership["validators"].append({
        "addr": tokens[0],
        "net_addr": tokens[1],
        "weight": "0",
    })

for line in sys.stdin.readlines():
    line = line.strip()
    if line != "": # Skip empty lines
        parse_validator(line)

# Printing the output of json.dumps instead of using directly json.dump to stdout,
# since the latter seems to append extra characters to the output.
print(json.dumps(membership, indent=4))
