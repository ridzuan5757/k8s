# OpenSearch Data Ingestion

There are seceral ways to ingest data into OpenSearch:
- Ingest individual documents.
- Index multiple document in bulk.
- Use Data Prepper - Server-side data collector that can enrich data for
  downstream analysis and visualisation.

## Bulk Indexing

To index documents in bulk, we can use the Bulk API:

```bash
POST _bulk
{ "create": { "_index": "students", "_id": "2" } }
{ "name": "Jonathan Powers", "gpa": 3.85, "grad_year": 2025 }
{ "create": { "_index": "students", "_id": "3" } }
{ "name": "Jane Doe", "gpa": 3.52, "grad_year": 2024 }
```
