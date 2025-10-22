FROM python:3.11-slim

# Set the working directory inside the container
WORKDIR /app

# Install core Jupyter dependencies and your project-specific libraries
# --no-cache-dir reduces the final image size
RUN pip install --no-cache-dir \
    jupyter_client \
    ipykernel \
    deap 

# Note: We do not specify a CMD here. 
# The Docker image only needs to contain the necessary libraries.
# The Jupyter Kernel Gateway (running in its own container) 
# will execute the specific command to start the kernel process (ipykernel_launcher) 
# inside a new instance of this image when requested.
