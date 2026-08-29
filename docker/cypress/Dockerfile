FROM cypress/included:15.21.0@sha256:af15ca5f3b851e32360f37e1d446d6e6821fd11928c87a2fd7fad12e8495e63b

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
