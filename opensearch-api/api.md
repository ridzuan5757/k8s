index template

```json
PUT /_index_template/otlp-metrics-template
{
    "index_patterns": ["otlp_metrics*"],
    "template":{
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 1,
            "codec": "best_compression"
        },
        "mappings":{
            "properties":{
                "position": {
                    "type": "geo_point"
                }
            }
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