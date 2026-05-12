# AGENTS.md

## Команды

- **Сборка фронтенда**: `cd web && npm run build`
- **Запуск сервера**: `go run ./cmd/server`
- **Сборка Go**: `go build -o server.exe ./cmd/server`
- **Docker**: `docker build -t invest-calc .`

## Архитектура

- Clean architecture: `domain` → `usecases` → `interfaces` → `infrastructure`
- Фронтенд собирается в `internal/interfaces/http/static/` и embed-ится в Go бинарник
- Все расчёты на сервере (Go), фронтенд только отображает и отправляет параметры

## Конвенции

- Go: стандартный `gofmt`, без комментариев без запроса
- TypeScript/React: без комментариев, следовать стилю существующих компонентов
- Форматирование денег: `formatMoney()` в `web/src/utils.ts`
- Деньги хранятся в строках с пробелами-разделителями тысяч, парсятся через `parseMoney()`
- Коммиты только по явному запросу пользователя
- После изменений в коде запускать `npm run build` и `go build` для проверки

## Важное

- `web/vite.config.ts` настроен на вывод в `../internal/interfaces/http/static/`
- `Dockerfile` multi-stage: node → golang → alpine
- Не коммитить `.vscode/`, `*.exe`, `node_modules/`, `__debug_bin*`
