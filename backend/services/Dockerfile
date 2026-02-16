# Step 1 - Start with the offical Python base image
FROM python:3.13-slim

# Step 2 - Set working directory inside the container
WORKDIR /app

# Step 3 - Copy the requirements first and install the dependencies
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt

# Step 4 - Copy the rest of the code
COPY . .

# Step 5 - Tell the container what to run on startup
CMD ["python3", "server.py"]