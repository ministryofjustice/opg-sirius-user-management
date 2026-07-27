FROM cypress/included:15.19.0@sha256:fe59f2cac5c7451644c8780cb2fb6735fde512d06642117199a7479269702f03

WORKDIR /e2e

COPY package.json .
COPY package-lock.json .
RUN npm ci --include=dev --ignore-scripts
