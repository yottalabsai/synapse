# Use Paketo Buildpacks to build the application
FROM alpine:3.20

# Set the working directory
WORKDIR /app

# Copy the Go application source code
COPY ./synapse /app
COPY ./resources /app/resources

# Expose the port the application runs on
EXPOSE 8002 50051

# Define an argument that can be passed at build time
ENV PROFILE k8s

# Set the entrypoint
CMD ["sh", "-c", "/app/synapse"]