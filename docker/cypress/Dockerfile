FROM cypress/included:15.18.0@sha256:5673defbe5ce3caeaf5eabbc6077c46d448e11fd7a64d4d5bfd269ec63c2a9dd

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
