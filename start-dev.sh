#!/bin/bash

# Start development servers for AWS Go Server + React frontend

echo "🚀 Starting AWS Go Server + React Development Environment"
echo "=========================================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    exit 1
fi

# Check if Node is installed
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is not installed"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo "✅ Node version: $(node --version)"
echo "✅ npm version: $(npm --version)"
echo ""

# Load environment variables from .env file
if [ -f .env ]; then
    echo "📝 Loading environment variables from .env..."
    export $(cat .env | grep -v '^#' | grep -v '^$' | xargs)
    echo "✅ Environment variables loaded"
else
    echo "❌ Error: .env file not found"
    echo "   Create .env file from .env.example:"
    echo "   cp .env.example .env"
    exit 1
fi

# Verify Cognito configuration
if [ -z "$AWS_COGNITO_USER_POOL_ID" ]; then
    echo "❌ Error: AWS_COGNITO_USER_POOL_ID not set in .env"
    exit 1
fi

if [ -z "$AWS_COGNITO_CLIENT_ID" ]; then
    echo "❌ Error: AWS_COGNITO_CLIENT_ID not set in .env"
    exit 1
fi

if [ -z "$AWS_COGNITO_CLIENT_SECRET" ]; then
    echo "❌ Error: AWS_COGNITO_CLIENT_SECRET not set in .env"
    exit 1
fi

echo "✅ Cognito configuration verified"
echo ""

# Function to cleanup background processes
cleanup() {
    echo ""
    echo "🛑 Shutting down servers..."
    kill $(jobs -p) 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM

# Start Go server in background
echo "📦 Starting Go server on http://localhost:8080..."
cd "$(dirname "$0")"
go run cmd/server/main.go &
GO_PID=$!

# Wait a bit for Go server to start
sleep 3

# Check if Go server is running
if ! curl -s http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "⏳ Waiting for Go server to start..."
    sleep 2
    if ! curl -s http://localhost:8080/healthz > /dev/null 2>&1; then
        echo "❌ Error: Go server failed to start"
        echo "   Check the logs above for errors"
        kill $GO_PID 2>/dev/null
        exit 1
    fi
fi

echo "✅ Go server is running"
echo ""

# Start React dev server
echo "⚛️  Starting React dev server on http://localhost:5173..."
cd web
npm run dev &
REACT_PID=$!

echo ""
echo "=========================================================="
echo "🎉 Development servers are running!"
echo ""
echo "  📱 React App:  http://localhost:5173"
echo "  🔧 Go API:     http://localhost:8080"
echo "  💚 Health:     http://localhost:8080/healthz"
echo ""
echo "Press Ctrl+C to stop all servers"
echo "=========================================================="

# Wait for background processes
wait
