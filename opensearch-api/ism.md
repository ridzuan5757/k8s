otlp-policy

```json
{
    "policy": {
        "policy_id": "otlp-policy",
        "description": "OTLP policy for hot-warm-cold-delete workflow.",
        "last_updated_time": 1725270084128,
        "schema_version": 21,
        "error_notification": null,
        "default_state": "hot",
        "states": [
            {
                "name": "hot",
                "actions": [
                    {
                        "retry": {
                            "count": 3,
                            "backoff": "exponential",
                            "delay": "1m"
                        },
                        "index_priority": {
                            "priority": 100
                        }
                    }
                ],
                "transitions": [
                    {
                        "state_name": "warm",
                        "conditions": {
                            "min_size": "100gb"
                        }
                    }
                ]
            },
            {
                "name": "warm",
                "actions": [
                    {
                        "retry": {
                            "count": 5,
                            "backoff": "exponential",
                            "delay": "5m"
                        },
                        "index_priority": {
                            "priority": 50
                        }
                    },
                    {
                        "retry": {
                            "count": 3,
                            "backoff": "exponential",
                            "delay": "1m"
                        },
                        "rollover": {
                            "min_doc_count": 5,
                            "min_size": "100gb",
                            "copy_alias": false
                        }
                    }
                ],
                "transitions": [
                    {
                        "state_name": "cold",
                        "conditions": {
                            "min_index_age": "1d"
                        }
                    }
                ]
            },
            {
                "name": "cold",
                "actions": [
                    {
                        "retry": {
                            "count": 3,
                            "backoff": "exponential",
                            "delay": "1m"
                        },
                        "close": {}
                    }
                ],
                "transitions": [
                    {
                        "state_name": "delete",
                        "conditions": {
                            "min_index_age": "1h"
                        }
                    }
                ]
            },
            {
                "name": "delete",
                "actions": [
                    {
                        "retry": {
                            "count": 3,
                            "backoff": "exponential",
                            "delay": "1m"
                        },
                        "delete": {}
                    }
                ],
                "transitions": []
            }
        ],
        "ism_template": [
            {
                "index_patterns": [
                    "otlp-*"
                ],
                "priority": 100,
                "last_updated_time": 1725268896123
            }
        ]
    }
}
```