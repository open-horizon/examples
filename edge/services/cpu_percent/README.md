# Horizon CPU Percent Service

The CPU Percent service provides the edge host's current CPU percentage being used to other Horizon services. It is a sample service useful when experimenting with Horizon, because it does not require any special hardware and produces constantly changing values.

## Input Values

This service takes no input values.

## RESTful API

Other Horizon services can use the CPU Percent service by requiring it in its own service definition, and then in its code accessing the CPU Percent REST APIs with the URL:
```
http://cpu:8347/v1/<api-from-the-list-below>
```

### **API:** GET /cpu
---

#### Parameters:
none

#### Response:

code: 
* 200 -- success
* other http codes TBD

body:


| Name | Type | Description |
| ---- | ---- | ---------------- |
| cpu | float | the cpu percent currently being used on this edge node host |


#### Example:
```
curl -sS -w "%{http_code}" http://cpu:8347/v1/cpu | jq .
{
  "cpu": 5.05
}
200
```
