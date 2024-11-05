index template

```json
DELETE /_index_template/otlp-metrics-template
PUT /_index_template/otlp-metrics-template
{
    "index_patterns": ["otlp-metrics*"],
    "priority": 100,
    "template":{
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 2,
            "codec": "best_compression",
            "plugins.index_state_management.rollover_alias": "otlp-metrics"
        },
        "mappings": {
            "properties": {
                "position": {
                    "type": "geo_point"
                }
            }
        }
    }
}
PUT /otlp-metrics-000001
POST /otlp-metrics-000001/_alias/otlp-metrics

DELETE /_index_template/otlp-logs-template
PUT /_index_template/otlp-logs-template
{
    "index_patterns": ["otlp-logs*"],
    "priority": 100,
    "template":{
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 2,
            "codec": "best_compression",
            "plugins.index_state_management.rollover_alias": "otlp-logs"
        }
    }
}
PUT /otlp-logs-000001
POST /otlp-logs-000001/_alias/otlp-logs
```