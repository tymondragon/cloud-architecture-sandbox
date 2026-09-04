# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:alpine AS backend-builder

WORKDIR /backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Copy frontend dist
COPY --from=frontend-builder /frontend/dist ./pkg/api/frontend/dist

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /cluster-viz .

# Stage 3: Final distroless image
FROM gcr.io/distroless/static-debian12

COPY --from=backend-builder /cluster-viz /cluster-viz

EXPOSE 8080

ENTRYPOINT ["/cluster-viz"]
