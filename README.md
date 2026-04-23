# Chirpy

Бэкенд-сервис для соцплатформы в духе Twitter: пользователи публикуют короткие посты (*chirps*), авторизуются через JWT, оформляют премиум-подписку. Написан как pet-проект, чтобы на одном сервисе пройти полный цикл REST-бэкенда — от схемы БД до refresh-токенов и health-check'ов.

## Стек

- **Go** (`net/http`) — без фреймворков, на стандартной библиотеке
- **SQLite** + [**sqlc**](https://sqlc.dev) — типобезопасные запросы к БД из Go-кода
- **JWT** (access + refresh пара)
- **bcrypt** — хранение паролей
- **goose** / `.sql` миграции — версионирование схемы

## Возможности

- Регистрация и логин с выдачей пары access + refresh токенов
- Защищённые маршруты под JWT-middleware
- CRUD постов: создание, чтение списком и по id, удаление (только автором)
- Валидация длины и содержимого поста
- Отзыв refresh-токена (logout)
- Премиум-подписка через webhook от платёжного провайдера
- `/api/healthz` — readiness-проба
- `/admin/metrics` — счётчик обращений для наблюдаемости
- `/admin/reset` — сброс состояния в dev-окружении

## Структура

```
Chirpy/
├── sql/                     # SQL-миграции и запросы для sqlc
├── internal/                # Внутренние пакеты: auth, database, ...
├── assets/                  # Статика для index.html
├── handler_users_create.go  # POST /api/users
├── handler_login.go         # POST /api/login
├── handler_refresh.go       # POST /api/refresh
├── handler_chirps_get.go    # GET  /api/chirps
├── handler_delete_chirp.go  # DELETE /api/chirps/{id}
├── handler_subscribe.go     # POST /api/polka/webhooks
├── handler_validate.go      # Валидация тела поста
├── metrics.go               # /admin/metrics
├── readiness.go             # /api/healthz
├── reset.go                 # /admin/reset (dev-only)
├── json.go                  # Хелперы respondWithJSON / respondWithError
├── main.go
└── sqlc.yaml
```

## Запуск

**Требования:** Go 1.22+

```bash
git clone https://github.com/ThroughTheThornsToTheStarss/Chirpy.git
cd Chirpy
cp .env.example .env   # заполнить секреты (см. ниже)
go run .
```

Сервис поднимется на `http://localhost:8080`.

## Переменные окружения

| Переменная   | Назначение                                           |
|--------------|------------------------------------------------------|
| `DB_URL`     | Путь к SQLite-файлу                                  |
| `JWT_SECRET` | Секрет для подписи JWT                               |
| `POLKA_KEY`  | Ключ для webhook-авторизации от платёжного провайдера|
| `PLATFORM`   | `dev` включает эндпоинт `/admin/reset`               |

## Примеры запросов

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"hunter2"}'
```

**Логин — получить пару токенов:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"hunter2"}'
```

Ответ:
```json
{
  "id": "a1b2...",
  "email": "user@example.com",
  "token": "<access_jwt>",
  "refresh_token": "<refresh>"
}
```

**Создать пост:**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H 'Authorization: Bearer <access_jwt>' \
  -H 'Content-Type: application/json' \
  -d '{"body":"Первый chirp!"}'
```

**Обновить access-токен:**
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H 'Authorization: Bearer <refresh>'
```

**Отозвать refresh (logout):**
```bash
curl -X POST http://localhost:8080/api/revoke \
  -H 'Authorization: Bearer <refresh>'
```

## Решения, которые были в центре внимания

- **Access + refresh пара, а не вечный токен.** Access живёт 1 час, refresh — 60 дней, отзывается по `/api/revoke`. Если access украдут — окно атаки ограничено, refresh можно отозвать целенаправленно.
- **`sqlc` вместо ORM.** Пишется чистый SQL в `sql/queries/`, `sqlc` генерирует типобезопасные Go-обёртки. В итоге — без скрытых N+1 и рантайм-сюрпризов, при этом код хендлеров остаётся чистым.
- **Middleware-авторизация.** JWT проверяется в одном месте, `user_id` прокидывается в контекст запроса — хендлеры не дёргают заголовки руками.
- **Централизованные JSON-ответы.** Все хендлеры ходят через `respondWithJSON` / `respondWithError` из `json.go` — единый формат ошибок, никаких ручных `w.Write([]byte(...))`.
- **Readiness и metrics с первого коммита.** Держал в голове, что сервис должен уметь жить за балансировщиком: `healthz` для проверок, `metrics` для наблюдаемости.

## Статус

Pet-проект. Используется как песочница для обкатки паттернов, которые потом применяю в продакшн-коде — в том числе на стажировке Go-разработчиком в amoCRM.
