#!/bin/bash
# Скрипт автоматической пересборки и перезапуска бота Chatan
# Разместить в ~/chatan/scripts/auto-rebuild.sh

# Конфигурация
PROJECT_DIR="/home/kkorsakov/chatan"
LOG_FILE="$PROJECT_DIR/logs/rebuild.log"
BOT_PID_FILE="$PROJECT_DIR/bot.pid"
CONFIG_FILE="$PROJECT_DIR/config.yaml"

# Создаем директорию логов при необходимости
mkdir -p "$PROJECT_DIR/logs"
touch "$LOG_FILE"
chown kkorsakov:kkorsakov "$LOG_FILE"
chmod 644 "$LOG_FILE"

# Функция для логирования
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Функция остановки бота
stop_bot() {
    if [ -f "$BOT_PID_FILE" ]; then
        local pid=$(cat "$BOT_PID_FILE")
        if ps -p "$pid" > /dev/null; then
            log "Останавливаем бота (PID: $pid)..."
            kill "$pid" && rm "$BOT_PID_FILE"
            sleep 2  # Даем время на корректное завершение
        else
            rm "$BOT_PID_FILE"
        fi
    fi
}

# Функция запуска бота
start_bot() {
    log "Запускаем бота..."
    cd "$PROJECT_DIR" || exit 1
    
    # Запускаем и сохраняем PID
    nohup ./chatan >> "$PROJECT_DIR/logs/bot.log" 2>&1 &
    echo $! > "$BOT_PID_FILE"
    log "Бот запущен с PID: $!"
}

# Основной процесс
log "=== Начало пересборки ==="

# 1. Проверяем конфигурацию
if [ ! -f "$CONFIG_FILE" ]; then
    log "Ошибка: config.yaml не найден!"
    exit 1
fi

# 2. Останавливаем текущий процесс бота
stop_bot

# 3. Переходим в директорию проекта
cd "$PROJECT_DIR" || exit 1

# 4. Обновляем зависимости
log "Обновляем зависимости..."
go get -u gopkg.in/telebot.v3
go get -u gopkg.in/yaml.v3
go mod tidy

# 5. Пересобираем проект
log "Пересобираем проект..."
if ! go build -o chatan cmd/chatan/main.go; then
    log "Ошибка сборки! Подробности:"
    go build -o chatan cmd/chatan/main.go 2>&1 | tee -a "$LOG_FILE"
    log "Попробуйте выполнить вручную:"
    log "1. cd $PROJECT_DIR"
    log "2. go mod tidy"
    log "3. go build -o chatan cmd/chatan/main.go"
    exit 1
fi

# 6. Запускаем бота
start_bot

# 7. Проверяем статус
sleep 1
if ! ps -p $(cat "$BOT_PID_FILE") > /dev/null; then
    log "Ошибка: бот не запустился! Проверьте логи:"
    tail -n 20 "$PROJECT_DIR/logs/bot.log" | tee -a "$LOG_FILE"
    exit 1
fi

log "=== Пересборка успешно завершена ==="
log "Логи бота: $PROJECT_DIR/logs/bot.log"
log "Логи пересборки: $LOG_FILE"