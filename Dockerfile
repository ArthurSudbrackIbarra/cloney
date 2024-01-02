# Stage 1:

# Use a recent version of Alpine Linux as the base image.
FROM alpine:3.9.6 as builder

# Update the package repository on the Alpine system.
RUN apk update

# Install the required packages.
RUN apk add --no-cache \
  bash \
  curl \
  unzip

#! Set the version of Cloney to install.
#! Use the --build-arg flag to override this value.
#! Example: docker build --build-arg CLONEY_VERSION=1.1.0 .
ARG CLONEY_VERSION=1.1.0

# Download Cloney.
RUN curl -sSL \
  "https://raw.githubusercontent.com/ArthurSudbrackIbarra/cloney/${CLONEY_VERSION}/installation/install.sh" | bash

# Stage 2:

# Use a recent version of Alpine Linux as the base image.
FROM alpine:3.9.6 as cloney

# Update the package repository on the Alpine system.
RUN apk update

# Copy the application binary from the builder stage.
COPY --from=builder /usr/local/bin/cloney /usr/local/bin

# Configure the permission to execute the binary.
RUN chmod +x /usr/local/bin/cloney

# Set the entrypoint to the Cloney binary.
ENTRYPOINT ["/usr/local/bin/cloney"]
