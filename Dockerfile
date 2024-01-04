# Use Ubuntu 24.04 as the base image.
FROM ubuntu:24.04

#! Set the default version of Cloney to install.
#! You can override this value using the --build-arg flag.
#! Example: docker build --build-arg CLONEY_VERSION=1.1.0 .
ARG CLONEY_VERSION=1.1.0
#! --------------------------------------------------------

# Update the package repository on the container.
RUN apt-get update

# Install essential packages (curl and unzip) for the installation process.
RUN apt-get install -y \
  curl \
  unzip

# Download Cloney using the specified version from the official repository.
RUN curl -sSL \
  "https://raw.githubusercontent.com/ArthurSudbrackIbarra/cloney/${CLONEY_VERSION}/installation/install.sh" | bash

# Grant execute permissions to the Cloney binary.
RUN chmod +x /usr/local/bin/cloney

# Set the default command to sleep indefinitely, providing a placeholder for future commands.
CMD ["sleep", "infinity"]
