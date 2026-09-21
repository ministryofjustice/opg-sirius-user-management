FROM cypress/included:15.19.0@sha256:418f70c94f6d727580da697faf8f4ca50c40fda78a57590708fa57772237697c

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
