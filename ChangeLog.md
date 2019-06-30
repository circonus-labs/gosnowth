# Changelog

We document mentionable user facing changes to the gosnowth library here. We
structure these changes according to gosnowth releases. Release versions adhere
to [Semantic Versioning](http://semver.org/) rules.

## [Next Release]

### Added

- A new field has been added to the FindTagsItem structure returned by calls to
SnowthClient.FindTags(). The field is called Activity (JSON: `activity`), and
contains the activity data returned by the IRONdb find tags API.

## [v1.2.0] - 2019-06-25

### Added

- Adds SnowthClient.GetStats() functionality. This retrieves metrics and stats
data about an IRONdb node via the /stats.json API endpoint.
- The Stats type is defined to hold the metric data returned by the GetStats()
operation. It stores the data in a map[string]interface{}, allowing the metrics
exposed by IRONdb to change without breaking gosnowth.
- Helper methods are defined on the Stats type to check and retrieve commonly
used information, such as IRONdb version and identification information.
- Adds an assignable middleware function that can run during the
SnowthClient.WatchAndUpdate() process. This allows downstream users of this
library to implement inspections and activate or deactivate node use according
to node information.

### Changed

- The code that creates and updates SnowthNode values has been changed to pull
information via GetStats() instead of GetState(), so that additional information
about the version of IRONdb running on a node can be obtained using the
SnowthNode value.

## [v1.1.3] - 2019-04-03

### Added

- Adds support for new check tags data returned from IRONdb to the SnowthClient.FindTags() methods.

## [v1.1.2] - 2019-03-13

### Added

- Adds context aware versions of all methods exposed by SnowthClient values.
These methods all contain a context.Context value as the first parameter, and
have the same name as their non-context variant with Context appended to the
end. These methods allow full support for IRONdb request cancellation via
context timeout or cancellation.
- Implements a Config type that can be used to pass configuration data when
creating new SnowthClient values. The examples provided in the [/examples]
folder demonstrate use of a Config type to configure SnowthClient values.

### Changed

- Includes account and check information in the data sent to IRONdb when
writing to histogram endpoints.

### Fixed

- Bug: Code in SnowthClient.WatchAndUpdate() could fire continuously, without
any delay, once started. Created: 2019-03-12. Fixed: 2019-03-13.

[Next Release]: https://github.com/circonus-labs/gosnowth
[v1.2.0]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.2.0
[v1.1.3]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.1.3
[v1.1.2]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.1.2
[/examples]: https://github.com/circonus-labs/gosnowth/tree/master/examples
