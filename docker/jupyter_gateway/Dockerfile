FROM python:3.11-slim

WORKDIR /app

# Install the Jupyter Kernel Gateway and the necessary client libraries
RUN pip install --no-cache-dir jupyter_kernel_gateway jupyter_client ipykernel

# Set the entry point to run the kernel gateway application
# It uses the default port 8888, which we configure in docker-compose.yml
ENTRYPOINT ["jupyter", "kernelgateway"]
