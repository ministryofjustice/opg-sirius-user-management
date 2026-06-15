FROM cypress/included:15.17.0@sha256:9378699679f7d6ddb63c572866d4f11ec822abacb099345e89f5b86381efd53f

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
