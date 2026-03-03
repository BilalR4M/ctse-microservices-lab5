FROM python:3.9-slim

WORKDIR /app

# Install dependencies (libpq-dev/gcc required for psycopg2)
RUN apt-get update \
    && apt-get install -y gcc libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py ./

EXPOSE 8083

# Unbuffered to show logs properly in docker
CMD ["python", "-u", "app.py"]
