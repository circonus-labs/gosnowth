# gosnowth
The Golang snowth client

## State and Topology

```
LocateMetric(uuid string, metric string, node *SnowthNode) (*DataLocation, error)
GetTopologyInfo(node *SnowthNode) (*Topology, error)
LoadTopology(hash string, topology *Topology, node *SnowthNode) error
ActivateTopology(hash string, node *SnowthNode) error
GetNodeState(node *SnowthNode) (*NodeState, error)
GetGossipInfo(node *SnowthNode) (*Gossip, error)
GetTopoRingInfo(hash string, node *SnowthNode) (*TopoRing, error)
```

LocateMetric - used to locate a particular metric within a topology. Need to supply the check uuid string, and metric name string as well as a node.

GetTopologyInfo - get the topology that the passed in node belongs to

LoadTopology - load a new topology.  This does not activate, activate is a seperate call.  Provide the new topology hash, as well as a representation of the new topology and the node to perform the call against.

ActivateTopology - given the hash provided as a parameter, activate this topology.  This is Dangerous.

GetNodeState - get the state of the node that is provided as a parameter.

GetGossipInfo - get the gossip information from the node.  Will include all of the other node's gossip information that this node knows about.

GetTopoRingInfo - get the topology information for another topology.  Look up performed with the hash parameter against the node parameter.
