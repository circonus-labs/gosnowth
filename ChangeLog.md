# Changelog

We document mentionable user facing changes to the gosnowth library here. We
structure these changes according to gosnowth releases. Release versions adhere
to [Semantic Versioning](http://semver.org/) rules.

## [Next Release]

## [v1.5.2] - 2020-05-01

* Adds additional debug logging, if configured, to support better tracking of
IRONdb request retries.
* Fixes a bug that would cause failures on POST requests to IRONdb if the
request needed to be retried on more than one node.

## [v1.5.1] - 2020-04-30

* Adds an option, not used by default, to retry requests to IRONdb that fail
for reasons that might be resolved by retrying. The number of attempts can be
set using the SnowthClient.SetRetries() method. Delay will increase between
each successive retry attempt.
* Adds an option, off by default, but can be set with
SnowthClient.SetConnectRetries(), that allows the client to retry requests 
to IRONdb that fail due to network errors on other available nodes, up to a 
specified number of times. This can be used in conjuction with
SnowthClient.SetRetries() or on its own.
* Added a `Limit` field to `FindTagsOptions` struct for specifying the maximum
number of metric results returned from the IRONdb find call.

## [v1.5.0] - 2020-04-16

* Adds linting configuration to the project and includes linting cleanup.
* Changes the API for most of the library functions to make the parameter
specifying a snowth node variadic and optional. This is possible because
gosnowth can now correctly determine which node to use itself.
* Adds internal implementation of topology location services matching the
logic used by snowth clusters.

## [v1.4.6] - 2020-03-06

* Adds support for the Explain parameter and results for CAQL requests.
* Adds support for tracing and dumping request payload to stdout for diagnostic
purposes.

## [v1.4.5] - 2020-01-10

### Added

- Added a new FindTagsLatest type and a new 'Latest' field to the `FindTagsItem`
type to support `SnowthClient.FindTags()` returning the latest data values when
requested from the IRONdb find call.

### Changed

- Changed the signature of the `SnowthClient.FindTags()` and
`SnowthClient.FindTagsContext()` methods to accept a `*FindTagsOptions`
argument. This argument contains the values for the supported IRONdb find
operation query parameters.

## [v1.4.4] - 2019-11-21

### Added

- Added SnowthClient.DoRequest() to issue a custom HTTP request to IRONdb.
- Added SnowthClient.RebuildActivity() to request a rebuild of IRONdb activity
tracking data for a list of metrics supplied in the new type RebuildActivityRequest.
- Added SnowthClient.WriteRawMetricList() convenience function to support writing
raw metric data directly to IRONdb via FlatBuffers Objects.

### Changed

- Changed the signature for SnowthClient.WriteRaw() to return the status
result of the /raw operation with the new type IRONdbPutResponse.

### Fixed

- Bug (severe): An issue was fixed that resulted in panics when rollup result
data, containing null values in specific places, was decoded from JSON format.
Created: 2019-09-16. Fixed: 2019-11-15.

## [v1.4.3] - 2019-09-20

### Changed

- Modified the internal structure of the RollupValue and RollupAllValue data
types, which are returned by the SnowthClient.ReadRollupValues() and
SnowthClient.ReadRollupAllValues() methods, to better express results returned
by IRONdb that contain `null` data values.
- Modified the internal structure of the TextValue data type, which is returned
by the SnowthClient.ReadTextValues() methods, to better express results returned
by IRONdb that contain `null` data values.

## [v1.4.2] - 2019-09-20

### Changed

- Changed the signature for the SnowthClient.GetCAQLQuery() methods to use the
new CAQLQuery type as a parameter. This allows all available parameters to be
used when executing CAQL queries.
- Added the CAQLError type which may be returned by the
SnowthClient.GetCAQLQuery() methods as an error if the error returned from the
corresponding IRONdb API call can be represented from this type. This allows
retrieval of extended error information when CAQL query requests fail.
- The SnowthClient.GetCAQLQuery() methods now send CAQL query requests to IRONdb
via a POST request. This prevents potential problems with query string encoding.

## [v1.4.1] - 2019-09-19

### Fixed

- Bug (moderate): The encoding used for CAQL queries was causing a parsing error
when the queries contained spaces. Created: 2019-09-15 Fixed: 2019-09-19.

## [v1.4.0] - 2019-09-15

### Added

- Adds the FetchQuery type and the SnowthClient.FetchValues and
SnowthClient.FetchValuesContext() methods to support fetching data, in the DF4
format, using the IRONdb /fetch API.
- New functionality has been added to read histogram data using the SnowthClient
ReadHistogramValues() and ReadHistogramValuesContext() methods. The new
HistogramValue type has been added to represent the histogram data returned by
the new methods.
- New ReadRollupAllValues() and ReadRollupAllValuesContext() methods have been
added to SnowthClient. These methods return slices of RollupAllValue values
representing IRONdb rollup response data in the legacy / type=all format.
- Adds new SnowthClient methods for listing Lua extensions and calling any Lua
extension via the IRONdb extension APIs: SnowthClient.GetLuaExtensions() and
SnowthClient.ExecLuaExtension().
- Adds specific support for performing CAQL queries using the Lua extension with
the SnowthClient.GetCAQLQuery() method.
- Adds support for a Go data structure representation of the DF4 data format.
This is needed to represent CAQL query results and will also allow for future
support for IRONdb /fetch API requests.

### Changed

- The RollupValues type has been replaced by the RollupValue and RollupAllValue
types. These types are better able to represent all possible rollup data formats
that can be returned by IRONdb. This is an API breaking change that modifies the
signature of the SnowthClient ReadRollupValues() and ReadRollupValuesContext()
methods.
- The SnowthClient ReadTextValues() and ReadTextValuesContext() have updated
signatures to match the parameters of the other data retrieval methods. This is
an API breaking change for these methods.
- These new types allow true support for IRONdb formatted timestamps for rollup
data retrieval methods by changing previous Timestamp integer field to a field
named Time containing a Go time.Time value. This is translated to/from the
IRONdb timestamp format during JSON encoding/decoding and the IRONdb timestamp
can be retrieved (as as string) by calling the Timestamp() method on values of
the new types.
- The new methods also support retrieving all types of rollup data. They are no
longer restricted to only the average type data.

## [v1.3.2] - 2019-08-26

### Changed

- The integer size of the activity data and account ID data in the FindTagsItem
structure returned by calls to FindTags() have been changed from `int32` to
`int64`. As has the type of the account ID parameter in calls to FindTags().
- The integer size in the rollup values returned by calls to GetNodeState() have
been changed from `uint32` to `uint64`.

### Fixed

- Bug (severe): It is possible for some types of IRONdb data to deserialize into
values that overflow the 32-bit variables used to hold them. Created: 2019-07-01
Fixed: 2019-08-27.

## [v1.3.1] - 2019-08-26

### Changed

- The signature of the FindTags() and FindTagsContext() methods have changed.
This is a breaking change to the API. The results are now returned wrapped in
a \*FindTagsResults value. This allows the total results count value returned
by the IRONdb request to be returned to the gosnowth user. Upgrading to this
release will require modifying any use of these methods in your code to reflect
this change.
- The internal functionality of the client do() method has been modified. It no
longer attempts to decode the contents of a response within this method. It
returns the response body data back to the caller to be handled there.
Additionally, it now also passes response headers back to the caller, so that
if they contain any needed information, it can now be used.

## [v1.2.1] - 2019-07-01

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

- Adds support for new check tags data returned from IRONdb to the
SnowthClient.FindTags() methods.

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
[v1.5.2]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.5.2
[v1.5.1]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.5.1
[v1.5.0]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.5.0
[v1.4.6]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.6
[v1.4.5]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.5
[v1.4.4]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.4
[v1.4.3]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.3
[v1.4.2]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.2
[v1.4.1]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.1
[v1.4.0]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.4.0
[v1.3.2]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.3.2
[v1.3.1]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.3.1
[v1.2.1]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.2.1
[v1.2.0]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.2.0
[v1.1.3]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.1.3
[v1.1.2]: https://github.com/circonus-labs/gosnowth/releases/tag/v1.1.2
[/examples]: https://github.com/circonus-labs/gosnowth/tree/master/examples
