# Search Data

There are several ways to search data:
- Query Domain Specific Language `DSL` - The primary OpenSearch query language,
  which we can use to create complex, fully customizable queries.
- Query string query language - Scaled down query language that we can use in
  a query parameter of a search request or in OpenSearch Dashboard.
- `SQL`
- Piped Processing Language `PPL` - Primary language used for observability in
  OpenSearch. PPL uses a pipe syntax that chains commands into a query.
- Dashboards Query Language `DQL` - Simple text based query language for
  filtering data in OpenSearch Dashboards.

## Retrieve All Documents in an Index

To retrieve all documents in an index, use the following requests:

```bash
GET /students/_search
```

The preceding request is equivalent to the `match_all` query, which matches all
documents in an index:

```bash
GET /students/_search
{
  "query": {
    "match_all": {}
  }
}
```

OpenSearch returns the matching documents:

```json
{
  "took": 12,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 3,
      "relation": "eq"
    },
    "max_score": 1,
    "hits": [
      {
        "_index": "students",
        "_id": "1",
        "_score": 1,
        "_source": {
          "name": "John Doe",
          "gpa": 3.89,
          "grad_year": 2022
        }
      },
      {
        "_index": "students",
        "_id": "2",
        "_score": 1,
        "_source": {
          "name": "Jonathan Powers",
          "gpa": 3.85,
          "grad_year": 2025
        }
      },
      {
        "_index": "students",
        "_id": "3",
        "_score": 1,
        "_source": {
          "name": "Jane Doe",
          "gpa": 3.52,
          "grad_year": 2024
        }
      }
    ]
  }
}
```

## Response Fields

The preceding response contains the following fields:
- `took` contains the amount of time the query took to run in milliseconds.
- `timed_out` indicates whether request timed out. If a request timed out, then
  OpenSearch returns the results that were gathered before the timeout. We can
  set the desired timeout by providing the `timeout` query parameter:

  ```bash
  GET /students/_search?timeout=20ms
  ```
- `_shards` specifices total number of shards on which the query ran as well as
  well as the number of shards that succeeded or failed. A shard may fail if the
  shard itself and all its replicas are unavailable. If any of the involved
  shards fail, OpenSearch continues to run the query on the remaining shards.
- `hits` contains the total number of matching documents and the documents
  themselves (listed in the `hits` array). Each matching document contains the
  `_index` and `_id` fields as well as the `_source` field, which contains the
  complete originally indexed document.
- Each document is given a relvance score in the `_score` field. Because we ran
  a `match_all` search, all document scores are set to `1`. There is no
  difference in their relevance. The `max_score` field contains the highest
  score of any matching document.
