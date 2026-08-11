FROM cypress/included:15.20.0@sha256:5baa6e2f1831e7b9823246d22b5914fe77187a9bab030aa2f44fceabde17c9d1

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
