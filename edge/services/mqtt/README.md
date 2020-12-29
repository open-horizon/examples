# Horizon MQTT Broker Service

The MQTT Broker service provides an mqtt broker and publisher for inter-container commmunication. 

### **API:** 
---

#### Publishing:

```
mosquitto_pub [-d] [-h host] [-p port] {-f file | -m message} [-t topic]
```

#### Parameters:

| Name | Description |
| ---- | ---------------- |
| -h | host name | 
| -d | put the broker into the background after starting |
| -p | publish on the specified port |
| -f | send the contents of a file as the message | 
| -m | message payload to send | 
| -t | mqtt topic to publish to |

#### Subscribing:

```
mosquitto_sub [-h host] [-p port] [-t topic]
```

#### Parameters:

| Name | Description |
| ---- | ---------------- |
| -h | host name |  
| -p | subscribe on the specified port |
| -t | mqtt topic to subscribe to |
