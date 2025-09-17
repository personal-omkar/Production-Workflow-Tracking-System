#!/bin/bash

# Entries to be added to /etc/hosts
entries=(
"0.0.0.0       kanban.irpl.com"
)

# Function to check if entry exists
entry_exists() {
    local entry="$1"
    grep -q "$entry" /etc/hosts
}

# Add a newline only if the file is not empty and if the previous line is not empty
[ -s /etc/hosts ] && [ "$(tail -n 1 /etc/hosts)" ] && echo "" >> /etc/hosts

# Check and add entries to hosts file
for entry in "${entries[@]}"; do
    if ! entry_exists "$entry"; then
        # Append the entry
        echo "$entry" >> /etc/hosts
    fi
done