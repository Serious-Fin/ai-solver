FROM node:latest
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npx svelte-kit sync
RUN npm run build
EXPOSE 3000
CMD ["node", "build"]