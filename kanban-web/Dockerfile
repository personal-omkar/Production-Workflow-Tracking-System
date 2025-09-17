FROM scratch

# Copy the source code into the container
COPY . .
 
# Add your executable to the container
ADD kanban-web ./

# Set the entrypoint to your executable
ENTRYPOINT ["/kanban-web"]