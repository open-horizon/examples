#!/bin/bash

OBJECT_TYPES_STRING=$MMS_OBJECT_TYPES # operator need to put MMS_OBJECT_TYPES in the deployment as env from the userinput configmap, eg: "model model1 model2"
OBJECT_ID=config.json # this is the file we are watching in this example mms consumer service
                    # this file is in shared volume. Operator should bind as /ess-store/<objectType>-<objectID>

echo "Object types to check: $OBJECT_TYPES_STRING"
IFS=' ' read -r -a objecttypes <<< "$OBJECT_TYPES_STRING"

while true; do
    for objecttype in "${objecttypes[@]}"
    do
        MMS_FILE_NAME="/ess-store/$objecttype-$OBJECT_ID"
        echo "MMS_FILE_NAME: $MMS_FILE_NAME"
        if [[ -f $MMS_FILE_NAME ]]; then
            eval $(jq -r 'to_entries[] | .key + "=\"" + .value + "\""' $MMS_FILE_NAME)
            echo "$MMS_FILE_NAME: Hello from ${HW_WHO} from objectType: $objecttype, objectId: $OBJECT_ID!"
        else
            echo "file $MMS_FILE_NAME not found"
        fi
        
    done
    sleep 20 
done