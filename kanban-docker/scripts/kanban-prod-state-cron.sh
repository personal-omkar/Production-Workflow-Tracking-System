#!/bin/bash

# Infinite loop to execute modify_db.sh every 5 seconds
while true; do
    ./prod_state_update_script.sh
    sleep 5
done

