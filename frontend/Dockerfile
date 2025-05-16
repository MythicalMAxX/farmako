FROM node:18-alpine as builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./
RUN npm ci

# Copy the source code
COPY . .

# Build the application
RUN npm run build

# Serve with lightweight web server
FROM nginx:alpine

# Copy build output from builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy custom nginx config
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"] 