FROM cypress/included:15.21.1@sha256:8d874a8e5457459c7b8f185f797e46c299f7d0584c9fbda0b0ca23b3e85be749

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
