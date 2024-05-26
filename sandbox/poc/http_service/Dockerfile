FROM node:latest
WORKDIR /usr/src/app
COPY package*.json ./
USER root
RUN npm install
COPY . .
RUN npm run build
CMD ["node", "-r", "/usr/src/app/dist/instrumentation.js", "dist/main.js"]

