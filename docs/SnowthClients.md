# Snowth Client Functionality

| Category | Function | libsnowth | gosnowth | perl-snowth |
|:---|:---|:---:|:---:|:---:|
| Client Robustness | Ability to bootstrap nodes from single node | - | Yes (uses topology and gossip) | Yes (polls topology to isolate nodes) |
|  | Ability to only use "working" nodes in event of failure | - | Yes (uses gossip data to determine health of nodes) | No |
|  |  |  |  |  |
| State and Topology APIs | Get Node State | Yes (snowth_client_get_active_hash) | Yes (SnowthClient.GetNodeState) | Yes* (Snowth::Topo->initial_state) |
|  | Retrieve Gossip Data | No | Yes (SnowthClient.GetGossipInfo) | No |
|  | Load New Topology | - | Yes (SnowthClient.LoadTopology) | Yes (Snowth::Topo->load_remote_topo) |
|  | Activating a New Topology | - | Yes (SnowthClient.ActivateTopology) | No |
|  | Retrieve Topology Data | Yes (snowth_client_get_topo) | Yes (SnowthClient.GetTopologyInfo) | Yes* (Snowth::Topo->compile_topo) |
|  | Retrieve Toporing Data | - | Yes (SnowthClient.GetTopoRingInfo) | No |
|  | Retrieve Data Location | - | Yes (SnowthClient.LocateMetric) | Yes (Snowth::Topo->find) |
|  |  |  |  |  |
| Rebalancing APIs | Starting a Rebalance to the New Topology | - | No | No |
|  | Retrieve the rebalance state | - | No | No |
|  | Aborting a rebalance | - | No | No |
|  |  |  |  |  |
| Data Submission APIs | Write Raw Data | - | Yes (SnowthClient.WriteRaw) | Yes (Snowth::Pusher::RawNumeric) |
|  | Write Numeric Data | - | Yes (SnowthClient.WriteNNT) | Yes (Snowth::Pusher::NNT) |
|  | Write Text Data | - | Yes (SnowthClient.WriteText) | Yes (Snowth::Pusher::Text) |
|  | Write Histogram Data | - | Yes (SnowthClient.WriteHistogram) | Yes (Snowth::Pusher::Histogram) |
|  |  |  |  |  |
| Data Retrieval APIs | Retrieve Numeric Data | - | Yes (SnowthClient.ReadNNTAllValues & SnowthClient.ReadNNTValues) | Yes (Snowth::Fetcher::NNT) |
|  | Retrieve Text Data | - | Yes (SnowthClient.ReadTextValues) | Yes (Snowth::Fetcher::Text) |
|  | Retrieve Histogram Data | - | Yes (SnowthClient.ReadHistogramValues) | Yes (Snowth::Fetcher::Histogram) |
|  |  |  |  |  |
| Data Deletion APIs | Delete Numeric for Metric/Check | - | No | No |
|  | Delete Text for Metric/Check | - | No | No |
|  | Delete Histogram for Metic/Check | - | No | No |
|  | Delete all data prior to date | - | No | No |
|  |  |  |  |  |
| Lua Extensions API | Get List of Lua Extensions | - | Yes (SnowthClient.GetLuaExtensions) | No |
|  | Execute a Lua Extension | - | Yes (SnowthClient.ExecLuaExtension) | No |

\* Only works against the list of bootstrap nodes, until a state response is encountered, unable to target a specific node
