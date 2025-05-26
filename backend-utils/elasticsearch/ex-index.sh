#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


addr=${1:-"http://localhost:9200"}
auth=${auth:-""} # "-u elastic:password"

####
curl $auth $addr -f

curl $auth "$addr/_cat/indices?v" -f

curl $auth $addr/_cat/nodes -f

curl $auth "$addr/_cluster/health" -f

#### create index
curl $auth -X PUT "$addr/test-accounts" -H "Content-Type: application/json" -d @index.test-accounts.fulltext.json
# curl $auth -X PUT "$addr/test-accounts" -H "Content-Type: application/json" -d @index.test-accounts.fulltext.json

#### create a new index test-accounts_v2
# curl $auth -X PUT "$addr/test-accounts_v2"

#### copy data from test-accounts to test-accounts_v2
#curl $auth -X POST "$addr/_reindex" -H "Content-Type: application/json" \
#  -d '{"source": {"index": "test-accounts"},"dest": {"index": "test-accounts_v2"}}'

#### delete old index
# curl $auth -X DELETE "$addr/test-accounts"

#### alias test-accounts_v2 to test-accounts
#curl $auth -X POST "$addr/_aliases" \
#  -H "Content-Type: application/json" -d '{"actions":[{"add":{"index":"accounts_v2","alias": "accounts" }}]}'

#### auto generate id
# curl $auth -X PUT "$addr/test-accounts" -H "Content-Type: application/json" -d @docs.test-accounts.0001.json

#### put a document
curl $auth -X POST "$addr/test-accounts/_doc/0001" -f -H 'Content-Type: application/json' -d $'
{
  "firstname": "小明-上海",
  "lastname": "John",
  "email": "xiaoming@example.com",
  "phone": "13800138000",
  "level": "admin",
  "labels": ["test", "fake"]
}'

curl $auth -X GET "$addr/test-accounts/_search?pretty"

curl $auth -X GET "$addr/test-accounts/_mapping?pretty"

curl $auth "$addr/_analyze" -H "Content-Type: application/json" \
-d '{ "analyzer": "ik_smart", "text": "小明" }'

#### put multiply documents
curl $auth -X POST "$addr/_bulk" -H "Content-Type: application/x-ndjson" \
-d $'
{ "index": { "_index": "test-accounts", "_id": "002" } }
{ "firstname": "小蓝", "lastname": "Smith", "email": "xiaolan@example.com", "phone": "13600136000", "level": "editor", "labels": ["fake", "user"] }
{ "index": { "_index": "test-accounts", "_id": "003" } }
{ "firstname": "John", "lastname": "Smith", "email": "john@example.com", "phone": "1234567890", "level": "reviewer", "labels": ["fake", "tester"] }
'

#### search
curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "match": { "firstname": "小明" } }
}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "match": { "firstname": "蓝" }
} }'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "match": { "firstname": "John" }
} }'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "multi_match": {"query": "John", "fields": ["firstname", "lastname"] }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "multi_match": {"query": "小明", "fields": ["firstname", "lastname"]}
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "multi_match": {"query": "小明", "fields": ["firstname", "lastname"], "type": "phrase" }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "multi_match": {"query": "小明", "fields": ["firstname^3", "lastname"] }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "multi_match": {"query": "Smith", "fields": ["firstname", "lastname"] }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "term": { "email": "xiaoming@example.com" }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "match": { "firstname": { "query": "小明", "fuzziness": "AUTO" } }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" -d '{ "query": {
  "term": { "labels": "fake" }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" \
-d '{ "from": 0, "size": 2, "query": {
  "term": { "labels": "fake" }
}}'

curl $auth "$addr/test-accounts/_search" -H "Content-Type: application/json" \
-d '{ "from": 2, "size": 2, "query": {
  "term": { "labels": "fake" }
}}'


curl -X GET "$addr/test-accounts/_search" -H "Content-Type: application/json" \
-d '{ "from": 0, "size": 10,
  "sort": [{ "createdAt": { "order": "desc" } } ],
  "query": { "match_all": {} }
}'

exit
curl -X GET "$addr/test-accounts/_search" -H "Content-Type: application/json" \
-d '{ "from": 0, "size": 10,
  "sort": [{ "createdAt": { "order": "desc" } } ],
  "query": { "match_all": {} }
}'

curl -X GET "$addr/test-accounts/_search" -H "Content-Type: application/json" \
-d '{ "from": 0, "size": 10,
  "sort": [{ "createdAt": { "order": "desc" } } ],
  "query": { "match": { "fulltext": "John" } }
}'

exit
./bin/elasticsearch-plugin install \
  https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v8.12.2/elasticsearch-analysis-ik-8.12.2.zip
