# open-interlocking
An open source self hosted interlocking (simulation). For example you can use opinterlockingerk to connect it to your model railway or railway in minecraft.

## Translation table
The language of this project is english but you could know some german words:
| German word | English word                                                                            |
| ----------- | --------------------------------------------------------------------------------------- |
| Stellwerk   | interlocking                                                                            |
| Vorsignal   | distant signal                                                                          |
| Weiche      | railroad switch                                                                         |
| Bloecke     | blocks (blocks are occupiable parts of a track)                                         |
| Subbloecke  | subblocks (one subblock per block for each direction (the same block will be reserved)) |
| Fahrstrasse | a connection between two signals over railroad switches and blocks                      |
| Gleisbild   | track diagram                                                                           |


## Templating

You need to create a template file in which you explain how your track diagram should look like.

Give attention to this:
- Your distant signal names must start with 'V', your signal names must start with 'S', your railroad switch names must start with 'W', your block names must start with 'B' and the name of your subblocks must be the block name with lowercase letters added
- letter V, S, W and B must be followed by a number

The file has the following syntax:

```yaml

# define distant signals
distant_signals:
    V1:

# define signals
signals:
    S1:
        distant_signals:
            - V1 # that means distant signal V1 is at the same location as S1
    S2:

# define railroad switches
switches:
    W1:
    W2:

# define blocks
blocks:
    B1:

subblocks:
    B1a: # subblock for Block B1
        start: W1 # the subblock starts with railroad switch W1
        end: S1 # the subblock ends with signal S1
    B1b: # subblock for Block B1
        start: W2
        end: S2

relations:
    signals:
        S1:
            following: W1
        S2:
            following: W2
            previous: B1a
    switches:
        W1:
            previous: S1
            following_straight_blade: B1
            following_bending_blade: B2


```
This example won't work as it's only to show how the file should be written. You want an example? Go to [examples directory](examples).


Also we have a yaml schema:
1. Your yaml file must have the `open-interlocking.yaml` or `opinterlockingerk.yml` file suffix.
2. If this doesn't work (because your client won't load the schema from schemastore.org) you can tell your yaml autocomplete extension the schema manually: [template.schema.json](template.schema.json)

## Api

### Authorization

You must authorizate you with a Bearer token. Generate a new token by running open-interlocking with `-g -c [GIVE_THE_TOKEN_A_COMPUTER_NAME]`. You can supply a permission by adding ``-p permission`` (regex). The default permission is `.+`.
Put the token in the ``Authorization`` header field.

The following permissions are available:
| Permission         | Description                                                          |
| ------------------ | -------------------------------------------------------------------- |
| state/ack          | Acknowledge a state of a signal, distant signal or switch with POST. |
| occupy             | Occupy blocks                                                        |
| connection/set     | Set a connection between two signals                                 |
| connection/desolve | Desolve a connection between two signals                             |


There is a swagger/openapi specification: open-api.yml.

| Method | Path                            | Description                                                                                              | Response                                   |
| ------ | ------------------------------- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| GET    | /:kind/:name                    | Get the current state of a signal, distant_signal or switch                                              | true or false                              |
| POST   | /:kind/:name                    | Acknowledge the current_state. (first GET the state, be sure you set the signal before acknowledging it) | true or false (state)                      |
| GET    | /connection/:signal1/:signal2   | Get the connection between signal 1 and signal 2. (signal1 and signal2 are the name of the signal)       | take a look at the open api specification. |
| POST   | /connection/:signal1/:signal2   | Set a connection between these two signals.                                                              | take a look at the open api specification. |
| DELETE | /connection/:signal1/:signal2   | Desolve the connection between these two signals.                                                        | take a look at the open api specification. |
| POST   | /block/occupy/:from/:to/:action | Occupy the block 'to' (switch or block e.g. W1 or B1) (action = join or leave)                           | 'success'                                  |

### RailroadConnection state
connection.state is integer with this value:
| Integer | State                                                             |
| ------- | ----------------------------------------------------------------- |
| 0       | Connection ain't set                                              |
| 1       | Connection waiting until the switches acknowledged                |
| 2       | Connection waiting until the signals/distant signals acknowledged |
| 3       | Connection set                                                    |
| 4       | Connection desolving                                              |

### Block / Switch reserved
block.reserved / switch.reserved is integer with this value:
| Integer | State                                                     |
| ------- | --------------------------------------------------------- |
| 0       | Block is not reserved                                     |
| 1       | Block is reserving (waiting until the connection is set ) |
| 2       | Block is reserved for a connection                        |
| 3       | Block is occupied by a train                              |