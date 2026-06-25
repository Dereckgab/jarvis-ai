#!/bin/bash

# ============================================================
# JARVIS Full IA - Quick Start Script
# ============================================================
# This script sets up and starts the entire JARVIS ecosystem.
# Prerequisites: Docker and Docker Compose installed.
# ============================================================

set -e

echo "🤖 JARVIS Full IA - Starting Up..."
echo "======================================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    echo "   Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

# Create .env file from .env.example if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from .env.example..."
    cp backend/.env.example .env
    echo "⚠️  Please edit the .env file with your API keys before proceeding."
    echo "   Required: AI_DEEPSEEK_API_KEY or AI_OPENAI_API_KEY"
    echo ""
    read -p "Press Enter to continue after editing .env, or Ctrl+C to abort..."
fi

# Build and start all services
echo ""
echo "🔨 Building Docker images..."
docker compose build

echo ""
echo "🚀 Starting all services..."
docker compose up -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service health
echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "======================================"
echo "✅ JARVIS Full IA is now running!"
echo ""
echo "🌐 Frontend:  http://localhost:3000"
echo "🔧 Backend:   http://localhost:8080"
echo "📊 Qdrant UI: http://localhost:6333/dashboard"
echo ""
echo "📋 Useful commands:"
echo "   make logs          - View all logs"
echo "   make down          - Stop all services"
echo "   make clean         - Remove everything"
echo "   make status        - Check service status"
echo "======================================"
