// ecosystem.config.cjs
module.exports = {
    apps: [
        {
            name: "app",                    // Веб бэкенд
            script: "./app",                 // Скомпилированный бинарник из cmd/app
            instances: 1,
            exec_mode: "fork",
            watch: false,

            // Настройка логов
            error_file: "/dev/stderr",
            out_file: "/dev/stdout",
            log_file: "/dev/stdout",
            merge_logs: true,
            log_date_format: "YYYY-MM-DD HH:mm:ss",

            // Переменные окружения для веб бэкенда
            env: {
                APP_TYPE: "bot"  // Для идентификации типа приложения
            }
        },
        {
            name: "cron",                    // Крон приложение
            script: "./cron",                 // Скомпилированный бинарник из cmd/cron
            instances: 1,
            exec_mode: "fork",
            watch: false,

            // Настройка логов
            error_file: "/dev/stderr",
            out_file: "/dev/stdout",
            log_file: "/dev/stdout",
            merge_logs: true,
            log_date_format: "YYYY-MM-DD HH:mm:ss",

            // Переменные окружения для крона
            env: {
                APP_TYPE: "cron",
            }
        }
    ]
};