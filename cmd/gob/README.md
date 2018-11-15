# Gob

The Gob service (think: "go back") keeps track of what messages have been sent from the Hub and allows Apps to query what has been sent. Typically, you would want to do this during startup and if packet loss has occurred.

```text
  +-------------+                         +-------------+
  |             |                         |             |
  |     App     |                         |     Hub     |
  |             |                         |             |
  +------^------+                         +-----+-------+
         |                                      |
+---------------------------------------------------------+
|        |                                      |         |
| +------+------+          Gob            +-----v-------+ |
| |     TCP     |                         |    Hub      | |
| |  playback   |                         |  receiver   | |
| |             |                         |             | |
| +------^------+                         +-----+-------+ |
|        |                                      |         |
|        |                                      |         |
|        |                                      |         |
|        |            +-------------+     +-----v-------+ |
|        |            | Gob append- |     | Application | |
|        +------------+ only event  <-----+   logic     | |
|                     |    store    |     |             | |
|                     +-------------+     +-------------+ |
|                                                         |
+---------------------------------------------------------+
```

## Usage

This is used in two situations:

- When an App is starting up and needs to read up on what messages have been sent to build internal state for a session. Typically, it broadcast "who has sequence number 0, for the latest SessionID (0xffffffffffffffff)", and when it gets a response from a Gob, it will connect to it via TCP and request messages with sequence numbers 0 to the largest possible sequence number (0xffffffffffffffff). The Gob will send as many packets as it has, and then closes down the connection, leaving the App to resume normal online operation. The App should keep track of what message sequence numbers have been sent out already for its' own AppID, so that it doesn't re-send messages to the Hub uselessly.
- When an App experiences a gap in sequence numbers from the Hub. The App then asks the Gob for the messages with the missing sequence numbers, for the current SessionID.

## Internals

The Gob append-only event store is simply a struct with the following format.

```golang
type gobStore struct {
    data         map[uint64][][]byte
    lastSequence map[uint64]uint64
    lastSession  uint64
}
```
