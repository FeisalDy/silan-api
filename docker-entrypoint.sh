#!/bin/sh
set -e

echo "🚀 Starting production deployment..."

# Run database seeder
echo "🌱 Running database seeders..."
./seeder

# Start the main application
echo "▶️  Starting main application..."
exec ./main
