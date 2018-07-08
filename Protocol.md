Redox protocol specification
===

The redox protocol operates in these layers:

- TCP
  - SSL (optional)
    - redox/dict/v1
      - redox/protobuf


## redox/dict
The redox/dict layer is used to flatten a stream of protobuf messages into a continuous byte stream to send through the TCP or TCP/SSL layer.

Both parties of the connection track a trace of the last protobuf messages to build a dynamically-updated varuint-to-string dictionary encoding protobuf message types.

redox/dict is a single-direction protocol, where one party only sends (the sender) and one party only receives (the recipient). In the actual implementation, two instances of redox/dict are sent/received in opposite directions over the TCP/SSL socket. In other words, two dictionaries are built for different directions in a redox connection.

There are 4 types of transmission units:
- `HEADER`: a header showing that this is the redox/dict layer and exchanging the version. Format:
  - byte: const 0x00 (`UNIT_HEADER`)
  - byte array: const x'ff 72 65 64 6f 78 00 64 69 63 74 fe' (`HEADER_MAGIC`, 12 bytes)
  - uint16be: const 0x0001 (`HEADER_VERSION`)
- `SET_DICT_SPEC`: changes the specification for dictionary building. Format:
  - byte: const 0x01 (`UNIT_SET_DICT_SPEC`)
  - varuint: queue size
  - varuint: dictionary rebuild frequency
- `DEFINE_WORD`: defines a new dictionary entry. Format:
  - byte: const 0x02 (`UNIT_DEFINE_WORD`)
  - varuint: message type length in bytes
  - byte array: message type encoded in UTF-8
- `MESSAGE`: carries an actual protobuf message. Format:
  - byte: const 0x03 (`UNIT_MESSAGE`)
  - varuint: dictionary ID of the message type
  - varuint: message length in bytes
  - byte array: message contents encoded in redox/protobuf

A redox/dict connection first sends a `HEADER` unit and a `SET_DICT_SPEC` unit. Both parties then allocates
- `dictionary`, an empty tail-expandable string list (can be coupled with a hashtable from the opposite direction),
- `trace`, an empty uint queue with the length specified in the `SET_DICT_SPEC` unit, and
- `counter`, a uint initialized as 0.

When a `DEFINE_WORD` unit is sent, the message type ("word") is inserted into `dictionary`. The word's index in the dictionary (starting with 0) is the message ID.

When a `MESSAGE` unit is sent, the following is executed in the exact order:
- the message type is pushed to `trace`.
- `counter` increases by 1.
- the message type is resolved from the dictionary.
- if `counter` is equal to the last-received `SET_DICT_SPEC`'s frequency (if it is greater, the program should identify as a bug),
  - `counter` is reset to 0.
  - a message type frequency table is built from `trace`.
  - `dictionary` is resorted according to the following order:
    - the word's frequency in `trace`, in descending order
    - the word's original index in `dictionary`, in ascending order

If a `SET_DICT_SPEC` is sent again, `counter` and `trace` are reset.

## redox/protobuf
The redox/protobuf encodesmessages using the protobuf format and sends them through the redox/dict layer.
