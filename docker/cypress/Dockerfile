FROM cypress/included:15.11.0

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
