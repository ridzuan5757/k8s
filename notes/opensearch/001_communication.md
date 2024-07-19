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


