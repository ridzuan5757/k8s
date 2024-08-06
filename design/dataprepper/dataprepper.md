# Data Prepper capabilities

## Anomaly Detection

Data Prepper can be used to train machine learning models and generate 
anomalies in near real-time on time-series aggregated events. Anomalies can be
generated either on:
- Events generated within the pipeline
- Events coming directly into the pipeline, like OpenTelemetry metrics.

These aggregated time-series events can be fed to the `anomaly_detector`
processor, which trains a machine learning model and generates anomalies with a
grade score.

Pipeline can then be configured to write the anomalies to a separate index to
create document monitors and trigger fast alerting.

## Metrics generation from traces

Data Prepper can be used to derive metrics from OpenTelemetry traces. For
example, consider a
pipeline that receives incoming traces and extract a metric called
`durationInNanos`, aggregated over a window of 30 seconds. It then derives a
histogram from the incoming traces.

## Event aggregation

We can use Data Prepper to aggregate data from different events over a period of
time. Aggregating events can help to reduce unnecessary log volume and manage
use cases like multiline logs that are received as separate events. The
`aggregate` processor is a stateful processor that groups events based on the
values for a set of specified identification keys and performs a configurable
action for each group.

The `aggregate` processor state is stored in memory. For example, in order to
combine four events into one, the processor needs to retain pieces of the first
three events. The state of an aggregate group of events is kept for a
configurable amount of time. Depending on the log data, the aggregate action
being used, and the number of memory options in the processor configuration, the
aggregation could take place over long period of time.

### Basic Usage

Consider a pipeline that extracts:
- `sourceIp`
- `destinationIp`
- `port`

Using the `grok` processor and then aggregates on those fields over a period of
30 seconds using the `aggregate` processor and the `put_all` action. At the end
of the 30 seconds period, the aggregated log is sent to the OpenSearch sink.

Say, if the following batch of logs is sent to the pipeline:

```json
{ "log": "127.0.0.1 192.168.0.1 80", "status": 200 }
{ "log": "127.0.0.1 192.168.0.1 80", "bytes": 1000 }
{ "log": "127.0.0.1 192.168.0.1 80" "http_verb": "GET" }
```

The `grok` processor will extract keys such that the log events will look like
the following. These events now have the data that the `aggregate` processor
will need for the `identification_keys`.

```json
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "port": 80, "status": 200 }
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "port": 80, "bytes": 1000 }
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "port": 80, "http_verb": "GET" }
```

After 30 seconds, the `aggregate` processor writes the following aggregated log
to the sink:

```json
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "port": 80, "status": 200, "bytes": 1000, "http_verb": "GET" }
```

### Removing duplicates

Duplicate entries can be removed by deriving keys from incoming events and
specifying the `remove_duplicates` option for the `aggregate` processor. This
action immediately processes the first event for a group and drops all the
following events in that group.

Say, the first event is processed with the identification keys `sourceIp` and
`destinationIp`:

```json
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "status": 200 }
```

The pipeline will then drop the following event because it has the same keys:

```json
{ "sourceIp": "127.0.0.1", "destinationIp": "192.168.0.1", "bytes": 1000 }
```

The pipeline processes this event and creates a new group because the `sourceIp`
is different.

```json
{ "sourceIp": "127.0.0.2", "destinationIp": "192.168.0.1", "bytes": 1000 }
```

### Log aggregation and conditional routing

Multiple plugins can be used to combine log aggregation with conditional
routing. For example, this pipeline `log-aggregate-pipeline` receives logs by
using HTTP client and extracts important values from the logs by matching the
value in the `log` key.

Two of the values that the pipeline extracts from the logs with `grok` pattern
include `response` and `clientip`. The `aggregate` processor then uses the
`clientip` value, along with the `remove_duplicates` option, to drop any logs
that contain a `clientip` that has already been processed within the given
`group_duration`.

Three routes, or conditional statements exist in this pipeline. These routes
separate the value of the response into `2xx`, `3xx`, `4xx` and `5xx` responses.
Logs with a `2xx` or `3xx` status are sent to the `aggregated_2xx_3xx` index,
logs with a `4xx` status are sent to the `aggregated_4xx` index and logs with a
`5xx` status are sent to the `aggregated_5xx` index.

## Log enrichment

Different types of log enrichment can be performed with Data Prepper:
- Filtering
- Extracting key-value pairs from strings.
- Mutating strings.
- Converting lists to maps.
- Processing incoming timestamps.

### Filtering

`drop_events` processor can be used to filter specific log events before sending
them to a sink. For example, if we are collecting web request logs and only want
to store unsuccessful requests, we can create a pipeline that drops any requests
for which the response is less than 400 so that only log events with HTTP status
codes of 400 and higher remain.

### Extracting key-value pairs from strings

Log data often includes strings of key value pairs. For example, if a user
queries a URL that can be paginated, the HTTP logs might contain the following
HTTP query string.

```bash
page=3&q=my-search-term
```

To perform analysis using the search term, we can extract value of `q` from the
query string. The `key_value` processor provides robust support for extracting
keys and values from strings.

### Mutating events

There are several processors that allow event to be mutated:

| Mutate Event Processor | Usage                                                                                                       |
|------------------------|-------------------------------------------------------------------------------------------------------------|
| `add_entries`          | Add entries to an event                                                                                     |
| `convert_entry_type`   | Convert value types in an event                                                                             |
| `copy_values`          | Copy values within an event                                                                                 |
| `delete_entries`       | Delete entries from an event                                                                                |
| `list_to_map`          | Convert list of objects from an event where each object contains a `key` field into a map of target keys    |
| `map_to_list`          | Convert a map of object from an event, where each object contains a `key` field, into a list of target keys |
| `rename_keys`          | Rename keys in an event                                                                                     |
| `select_entries`       | Select entries from an event                                                                                |

For example, consider the following incoming event:

```json
{
   "key_one": "value_one",
   "key_two": "value_two"
}
```

Say, if we are using `add_entries` to add new `"key_three"` as the entries with
format `"${key_one}-${key_two}"`, the processor will transforms it into an event
with a new key `key_three`, which combines values of other keys in the original
event:

```json
{
   "key_one": "value_one",
   "key_two": "value_two",
   "key_three": "value_one-value_two"
}
```

### Mutating strings

We can change the way that a string appears by using mutate string processor.
For example, we can use the `uppercase_string` processor to convert a string to
uppercase, and we can use the `lowercase_string` processor to convert a string
to lowercase. The following is a list of processors that allow us to mutate a
string:
- `substitute_string`
- `split_string`
- `uppercase_string`
- `lowercase_string`
- `trim_string`

### Converting lists to maps

The `list_to_map` processor, which is one of the mutate event processors,
converts a list of objects in an event to a map. For example, consider the
following processor configuration:

```yaml
processor:
  - list_to_map:
      key: "name"
      source: "A-car-as-list"
      target: "A-car-as-map"
      value_key: "value"
      flatten: true
```

This processor will convert an event that contains a list of objects to a map
like this:

```json
{
  "A-car-as-list": [
    {
      "name": "make",
      "value": "tesla"
    },
    {
      "name": "model",
      "value": "model 3"
    },
    {
      "name": "color",
      "value": "white"
    }
  ]
}
```

```json
{
  "A-car-as-map": {
    "make": "tesla",
    "model": "model 3",
    "color": "white"
  }
}
```

### Processing incoming timestamps

The `date` processor can be used to generate timestamp for incoming events if we
specify `@timestamp` as the `destination` option. It is also capable to parses the 
`timestamp` key from incoming events by converting it to ISO 8601 format. Say, if 
the preceding pipeline processes the following event:

```json
{"timestamp": "10/Feb/2000:13:55:36"}
```

This processor can convert the event to the following format:

```json
{
  "timestamp":"10/Feb/2000:13:55:36",
  "@timestamp":"2000-02-10T15:55:36.000-06:00"
}
```

## Data Sampling

The following sampling capabilities is provided by Data Prepper:
- Time sampling
- Percentage sampling
- Tail sampling

### Time sampling

`rate_limiter` can be used together with `aggregate` processor to limit the
number of events that can be processed per second. We can choose to either drop
excess events or carry them forward to the next time period.

### Percentage sampling

`percent_sampler` can be used with `aggregate` processor to limit the number of
event that will be sent to the sink (OpenSearch). All excess events will be
dropped.

### Tail sampling

`tail_sampler` with `aggregate` processor can be used to sample events based on
a set of defined policies. This implementation waits an aggregation to complete
across aggregation periods based on the configured wait period.

When an aggregation is complete, and if it matches the speficic error condition,
it is sent tot the sink. Otherwise, only a configured percentage of events is
sent to the sink.

