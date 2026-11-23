FROM cypress/included:15.21.2@sha256:ed569d6dd485deec4aed62269687c166c34d0fb78153d90b694d9b044a1a9450

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
