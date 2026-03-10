#!/bin/sh
# start.sh - Скрипт запуска Go бэкенда (app + cron) через PM2

echo "=================================================="
echo "🚀 Freelance Backend (Go) - Starting with PM2"
echo "=================================================="
echo ""

# Выводим информацию о версиях
echo "$📦 PM2 version: $(pm2 --version)"
echo "📦 Node version: $(node --version)"
echo ""

# Проверка наличия бинарников
echo "🔍 Checking Go binaries..."
if [ ! -f /app/app ]; then
    echo "❌ ERROR: app binary not found!"
    exit 1
fi
if [ ! -f /app/cron ]; then
    echo "❌ ERROR: cron binary not found!"
    exit 1
fi
echo "✅ All binaries found"
echo ""

# Проверка переменных окружения
echo "🔍 Checking environment variables..."

# Основные переменные для веб приложения
REQUIRED_VARS="TOKEN, MYSQL_HOST"
MISSING_VARS=""

for VAR in $REQUIRED_VARS; do
    if [ -z "$(eval echo \$$VAR)" ]; then
        MISSING_VARS="$MISSING_VARS $VAR"
    fi
done

if [ ! -z "$MISSING_VARS" ]; then
    echo "❌ Missing required variables:$MISSING_VARS"
    exit 1
fi

if [ ! -z "$MYSQL_HOST" ]; then
    echo "   ✅ MYSQL_HOST=${MYSQL_HOST}"
fi
echo ""

echo "=================================================="
echo "🟢 Starting processes with PM2..."
echo ""

# Выводим информацию о запускаемых процессах
echo "📋 Processes to start:"
echo "   - app: telegram bot"
echo "   - cron: scheduled tasks"
echo ""

# Запускаем PM2
exec pm2-runtime /app/ecosystem.config.cjs \
    --output /dev/stdout \
    --error /dev/stderr