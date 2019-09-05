# Horizon Stop Word Removal Service

The Stop Word Removal service runs as a WSGI server that can take a JSON object such as {"text": "how are you today"} and will remove common stop words and return {"result": "how you today"}.

## Input Values

This service takes no input values.

## RESTful API

Other Horizon services can use the Stop Word Removal service by requiring it in its own service definition, and then in its code accessing the Stop Word Removal REST APIs with the URL:
```
http://ibm.stopwordremoval:80/remove_stopword 
```

### **API:** POST /remove_stopword
---

#### Parameters:
none

#### Response:

body:


| Name | Type | Description |
| ---- | ---- | ---------------- |
| result | string | text received with common stop words removed |


#### Example:
```
curl -X POST http://ibm.stopwordremoval:80/remove_stopword -H 'Content-Type: application/json' -H 'cache-control: no-cache' -d '{"text": "how are you today"}'
{
  "result": "how you today"
}
```
