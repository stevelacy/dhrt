# dhrt db
> distributed horizontal real time database

## Goals
> Distributed fault-tolerant horizontal scaling real-time database

This database should:
- be horizontally scalable -- adding new nodes to a cluster should have minimal overhead
- distributed
  - a single create request to any of the nodes should propagate to all nodes
  - not have any single point of failure, no master node
- primarily key/value
- fast -- this must be absolutely fast
- minimal -- the smallest number of features for use
- memory efficient
- memory mode and persisted mode, should default to persisted
- accept a custom path for data directory
- split data at "collections" (a collection is a set of similar data, such as a table)
  - collections could be placed on particular nodes to improve speed of read and writes.
  - each collection would be stored separately on it's own flat file


## Query commands

The CLI and REPL will support this syntax

- [ ] `set <key> <value>`
- [ ] `set <collection> <key> <value>`
- [ ] `get <key>`
- [ ] `get <collection> <key>`
- [ ] `del <key>`
- [ ] `stats`

## HTTP interface

```
GET / -- stats
GET /:key
  -- returns { Key, Value }
  -- gets value by key
POST /:key/:value
  -- returns { Key, Value }
  -- sets value by key
GET /:collection
  -- returns [ { Key, Value } ]
  -- gets all key/value pairs for a collection
GET /:collection/:key
  -- returns { Key, Value }
  -- gets value by key from collection (improves performance?)
POST /:collection/:key/:value
  -- returns { Key, Value }
  -- sets value by key for collection
DELETE /:key
  -- returns Key
  -- deletes key/value pair
DELETE /:collection/:key
  -- returns Key
  -- deletes key/value pair
DELETE /:collection
  -- returns [ Keys ]
  -- deletes entire key/value pair collection

```


#### Distributed searching
- get "item-1"
  - search first node [0]
    - if unable to find key, start propagating through the other nodes for the key until either key is found or return nil
    - sync.ErrGroup https://www.oreilly.com/learning/run-strikingly-fast-parallel-file-searches-in-go-with-sync-errgroup


#### Distributed networking
- Each node will have an exposed endpoint for tcp/http connections.
- On Open each node uses the address passed to it to connect and initiate a data handoff. Eg. `dhrt join 10.2.2.10:7453` or `./dhrt --node=10.2.2.10:7453`
  - This lets the initial node know that there is a new node on the cluster, it will then broadcast the new node's config to the existing cluster nodes including the new node.
- Cluster reads:
  - node1 will do a internal search.
  - if the datum is not found node1 will send a query to each of the cluster nodes until the datum is found or no results are found. Timeout in config will be followed
- Cluster writes:
  - node1 will do an internal save.
  - node1 will then broadcast a `set` command to each node in the cluster. Timeout in config will be followed.
- Cluster Node removal:
  - node1 will broadcast a node:leave command to the cluster upon a `SIGHUP`, `SIGINT`, `SIGTERM`, `SIGQUIT` system command.
    - all nodes in the cluster will then remove node1 from their node list
- Cluster Node fail:
  - node1 tries to make a request to node2 and fails (timeout)
  - node1 will then mark it as unhealthy until one of the following:
    - node1 attempts more requests to node2, if the requests succeed node2 is marked healthy
    - node1 attempts additional requests and failes, node2 is then marked as dead.
    - node2 will only be marked healthy once dead when it re-syncs with one of the nodes in the cluster. The re-sync will trigger a cluster wide broadcast of node2.




### Potential issues

- When a new node enters the cluster it may take a while for it to catch up with the existing nodes in the cluster when the dataset is large. There would potentially be a LRU cache in the handoff instead of the entire dataset, additional data would be passed over at a slower rate.
