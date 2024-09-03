index template

```json
PUT /_index_template/otlp-metrics-template
{
    "index_patterns": ["otlp-metrics-*"],
    "priority": 100,
    "template":{
        "aliases":{
          "otlp-metrics":{
            "is_write_index": true
          }
        },
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 1,
            "codec": "best_compression"
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
PUT /_index_template/otlp-logs-template
{
    "index_patterns": ["otlp-logs-*"],
    "priority": 100,
    "template":{
        "aliases":{
          "otlp-logs":{
            "is_write_index": true
          }
        },
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 1,
            "codec": "best_compression"
        }
    }
}
```

geo - ip data

```json
PUT otlp-metrics
{
  "mappings": {
    "properties": {
      "attributes.resource.attributes.position": {
        "type": "geo_point"
      },
      "position": {
        "type": "geo_point"
      }
    }
  }
}
```