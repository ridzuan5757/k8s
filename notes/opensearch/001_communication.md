# REST API

HTTP requests can be sent via terminal or in the Dev Tools Console in OpenSearch
Dashboards.

## Sending requests in a terminal

When sending cURL requests in a terminal, the request format varies depending on
whether Security plugin is being used or not.

If the Security plugin is not used:

```bash
curl -X GET "http://localhost:9200/_cluster/health"
```

If the Security plugin is used, username and password need to be provided in the
requests.

```bash
curl -X GET "http://localhost:9200/_cluster/health" -ku admin:<custom-admin-password>
```

The default username is `admin`, and the password is set in the
`docker-compose.yaml` file in the 
`OPENSEARCH_INITIAL_ADMIN_PASSWORD=<custom-admin-password>` setting. OpenSearch
generally returns response in a flat JSON format by default. For a
human-readable response body, provide the `pretty` query parameter.

```bash
curl -XGET "http://localhost:9200/_cluster/health?pretty"
```

## Indexing documents

To add JSON document to an OpenSearch index (this is called as indexing a
document), we send HTTP request with the following header:

```bash
PUT https://<host>:<port>/<index-name>/_doc/<document-id>

curl -XPUT "http://localhost:9200/students/_doc/1" -H 'Content-Type: application/json' -d'
{
  "name": "John Doe",
  "gpa": 3.89,
  "grad_year": 2022
}'
```

Once the request is sent, OpenSearch creates an index called `students` and
stores the ingested document in the index. If `document-id` is not provided for
the docuemnt, OpenSearch will generate the document Id.


## Dynamic mapping

When document is indexed, OpenSearch infers the field types from the JSON types
submitted in the document. This process is called dynamic mapping. To view the
inferred field types, send a request to the `_mapping` endpoint.

```bash
GET /students/_mapping
```

OpenSearch responds with the field `type` for each field

```json
{
    "students": {
        "mappings": {
            "properties": {
                "gpa": {
                    "type": "float"
                },
                "grad_year": {
                    "type": "long"
                },
                "name": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "query": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                }
            }
        }
    }
}
```

OpenSearch mapped:
- Numeric fields to the `float` and `long` types.
- `name` text is mapped to `text` and added `name.keyword` subfield mapped to
  `keyword`.
- `grad_year` is mapped to `long`.

Fields mapped to `text` are analyzed (lowercased and split into terms) and can
be used for full-text search. Fields mapped to `keyword` are used for exact term
search.

For `grad_year`, if we want to map it to the `date` type instead, we need to
delete the index and recreated it, explicitly specifying the mappings.

## Searching for documents

To run a search for the document, specify the index that we are searching and
query that will be used to match documents. The simplest query is the `match_all`
query, which matches all documents in an index:

```bash
GET /students/_search
{
  "query": {
    "match_all": {}
  }
}
```

The resulting value would be:

```json
{
    "took": 4,
    "timed_out": false,
    "_shards": {
        "total": 1,
        "successful": 1,
        "skipped": 0,
        "failed": 0
    },
    "hits": {
        "total": {
            "value": 1,
            "relation": "eq"
        },
        "max_score": 1.0,
        "hits": [
            {
                "_index": "students",
                "_id": "1",
                "_score": 1.0,
                "_source": {
                    "name": "John Doe",
                    "gpa": 3.89,
                    "grad_year": 2022
                }
            }
        ]
    }
}
```

## Updating Documents

In OpenSearch, documents are immutable. However, we can update a document by
retrieving it, updating its information and reindexing it. We can update entire
document using the Index Document API, providing values for all existing and
added fields in the document.

For example, to udpate the `gpa` field and add an `address` field to the
previously indexed document, we can use the following request:

```bash
PUT /students/_doc/1
{
  "name": "John Doe",
  "gpa": 3.91,
  "grad_year": 2022,
  "address": "123 Main St."
}
```

Alternatively, we can update parts of a document by calling the Update Document
API:

```bash
POST /students/_update/1/
{
  "doc": {
    "gpa": 3.91,
    "address": "123 Main St."
  }
}
```

The results would be:

```json
{
    "_index": "students",
    "_id": "1",
    "_version": 4,
    "result": "updated",
    "_shards": {
        "total": 2,
        "successful": 2,
        "failed": 0
    },
    "_seq_no": 3,
    "_primary_term": 1
}
```

## Deleting Document

To delete document, use delete request and provide the document ID.

```bash
DELETE /students/_doc/1
```

The results would be:

```json
{
    "_index": "students",
    "_id": "1",
    "_version": 5,
    "result": "deleted",
    "_shards": {
        "total": 2,
        "successful": 2,
        "failed": 0
    },
    "_seq_no": 4,
    "_primary_term": 1
}
```

## Deleting Index

To delete an index, use the following request:

```bash
DELETE /students
```

The results would be:

```bash
{
    "acknowledged": true
}
```

## Index Mappings and Settings

OpenSearch indexes are configured with mappings and settings:
- `mapping` is a collection of fields and the types of those fields.
- `settings` include index data like the index name, creation date and number of
  shards.

We can specify mapping and settings in one request. For example, the following
request:
- Specifies the number of index shards
- Maps the `name` field to `text`
- Maps the `grad_year` to `date`

```bash
PUT /students
{
  "settings": {
    "index.number_of_shards": 1
  }, 
  "mappings": {
    "properties": {
      "name": {
        "type": "text"
      },
      "grad_year": {
        "type": "date"
      }
    }
  }
}
```
The response would be:

```json
{
    "acknowledged": true,
    "shards_acknowledged": true,
    "index": "students"
}
```

Now we can index the same document that we indexed in the previous section:

```bash
PUT /students/_doc/1
{
  "name": "John Doe",
  "gpa": 3.89,
  "grad_year": 2022
}
```
