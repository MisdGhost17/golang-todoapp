# golang-todoapp

REST API для управления задачами на Go и PostgreSQL. Проект включает регистрацию, вход с JWT, ротацию refresh-токенов, личные задачи и статистику выполнения.

HTTP-слой построен на стандартном `net/http`, доступ к PostgreSQL — через `pgx`. Код разделён по функциональности и слоям: HTTP → сервис → репозиторий.

## Возможности

- Регистрация с уникальным логином и хешированием пароля через bcrypt.
- Access JWT с подписью HS256 и сроком действия 15 минут.
- Refresh-сессии на 30 дней: хранение SHA-256 хешей, одноразовое использование токенов, отзыв при повторном использовании.
- Выход из отдельной сессии.
- Создание, просмотр, изменение и удаление своих задач.
- Пагинация списка задач и проверка версии при обновлении.
- CRUD профилей пользователей и статистика задач.
- Swagger UI, SQL-миграции, структурированные логи и graceful shutdown.

## Требования

- Go версии, указанной в [go.mod](go.mod): `1.27.1`.
- Docker и Docker Compose v2.
- Make.
- OpenSSL для генерации JWT-секрета.

В Docker Compose используется PostgreSQL 18.1. Все команды ниже выполняются из корня репозитория.

## Локальный запуск

### 1. Настрой окружение

```sh
cp .env.example .env
openssl rand -hex 32
```

В `.env` заполни:

```dotenv
POSTGRES_USER=todoapp
POSTGRES_PASSWORD=replace_with_local_password
POSTGRES_DB=todoapp
JWT_SECRET=replace_with_generated_secret
```

В `JWT_SECRET` вставь результат `openssl rand -hex 32`. Не используй приведённую строку-заглушку как секрет. `.env` исключён из Git.

Остальные настройки можно оставить как в [.env.example](.env.example).

### 2. Запусти PostgreSQL

```sh
make env-up
make env-port-forward
```

База работает в Docker. Отдельный `port-forwarder` делает её доступной приложению на `127.0.0.1:5432`.

Дождись готовности PostgreSQL:

```sh
docker compose logs --tail=30 todoapp-postgres
```

### 3. Примени миграции

```sh
make migrate-up
```

Миграции создают схему `todoapp`, таблицы пользователей, задач, данных входа и refresh-сессий. Приложение не применяет их автоматически.

### 4. Запусти API

```sh
make todoapp-run
```

По умолчанию:

- API: `http://localhost:5050/api/v1`
- Swagger UI: `http://localhost:5050/swagger/index.html`
- Swagger JSON: `http://localhost:5050/swagger/doc.json`

Команда `todoapp-run` загружает переменные через Makefile, устанавливает `POSTGRES_HOST=localhost` и каталог логов `out/logs`, выполняет `go mod tidy`, затем запускает приложение. Само приложение файл `.env` не читает.

Для остановки API нажми `Ctrl+C`. Остановить локальное окружение:

```sh
make env-port-close
make env-down
```

## Конфигурация

| Переменная | Назначение |
| --- | --- |
| `HTTP_ADDR` | Адрес HTTP-сервера, например `:5050` |
| `HTTP_SHUTDOWN_TIMEOUT` | Время graceful shutdown, по умолчанию `30s` |
| `POSTGRES_HOST` | Хост PostgreSQL; при локальном запуске задаётся Makefile |
| `POSTGRES_PORT` | Порт PostgreSQL, по умолчанию `5432` |
| `POSTGRES_USER` | Пользователь базы |
| `POSTGRES_PASSWORD` | Пароль базы |
| `POSTGRES_DB` | Имя базы |
| `POSTGRES_TIMEOUT` | Таймаут операций с базой, в примере `10s` |
| `JWT_SECRET` | Случайный секрет подписи HS256; проверяется минимум 32 байта |
| `LOGGER_LEVEL` | Уровень логирования, в примере `DEBUG` |
| `LOGGER_FOLDER` | Каталог логов; локально задаётся Makefile |
| `TIME_ZONE` | Часовой пояс приложения, по умолчанию `UTC` |

Издатель JWT (`todoapp`), срок access-токена (15 минут) и срок refresh-сессии (30 дней) сейчас заданы в коде.

## Пример работы с API

### Регистрация

```sh
curl -i http://localhost:5050/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
    "full_name": "Test User",
    "login": "auth_test",
    "password": "learning-go-auth-2026"
  }'
```

Успех — `201 Created` с профилем пользователя. Регистрация не выдаёт токены. Повторный логин — `409 Conflict`.

Логин нормализуется: пробелы по краям удаляются, регистр становится нижним. Допустимы 3–64 символа `a-z`, `0-9`, `_`. Пароль должен содержать минимум 15 символов и максимум 72 байта UTF-8; он не нормализуется.

### Вход

```sh
curl -i http://localhost:5050/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"auth_test","password":"learning-go-auth-2026"}'
```

Ответ `200 OK`:

```json
{
  "access_token": "<access JWT>",
  "refresh_token": "<refresh token>",
  "token_type": "Bearer"
}
```

Каждый вход создаёт отдельную сессию. Сохрани полученные значения локально; не публикуй их в логах, issue или коммитах.

### Создание и получение задач

Вместо `<access JWT>` подставь токен из ответа входа:

```sh
curl -i http://localhost:5050/api/v1/tasks \
  -H 'Authorization: Bearer <access JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"title":"Разобраться с JWT","description":"Добавить авторизацию в API"}'

curl -i 'http://localhost:5050/api/v1/tasks?limit=20&offset=0' \
  -H 'Authorization: Bearer <access JWT>'
```

Автор определяется по проверенному JWT. Поле `author_user_id` в запросе создания и параметр `user_id` в запросе списка не позволяют выбрать другого пользователя. Без `limit` возвращаются все свои задачи.

Частичное обновление:

```sh
curl -i -X PATCH http://localhost:5050/api/v1/tasks/1 \
  -H 'Authorization: Bearer <access JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"completed":true}'
```

Непереданные поля сохраняются. `description: null` удаляет описание; `title` и `completed` не допускают `null`. Для чужой или отсутствующей задачи возвращается `404`.

### Обновление токенов

```sh
curl -i http://localhost:5050/api/v1/auth/refresh \
  -H 'Content-Type: application/json' \
  -d '{"refresh_token":"<refresh token>"}'
```

Ответ содержит новую пару токенов. Старый refresh-токен больше нельзя использовать. Обновление не продлевает абсолютный срок сессии.

Повторное использование старого refresh-токена отзывает всю сессию, включая её новый refresh-токен. Клиент должен выполнять только одно обновление одновременно и не повторять запрос автоматически при потере ответа: может потребоваться повторный вход.

### Выход

```sh
curl -i http://localhost:5050/api/v1/auth/logout \
  -H 'Content-Type: application/json' \
  -d '{"refresh_token":"<refresh token>"}'
```

Ответ — `204 No Content`. Повторный выход также успешен. Клиенту нужно удалить оба токена.

**Выход отзывает refresh-сессию, но уже выданный access JWT остаётся действительным до своего `exp`.** Refresh и logout не требуют access JWT.

## Маршруты

Все пути имеют префикс `/api/v1`.

| Метод | Путь | Доступ / назначение |
| --- | --- | --- |
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход по логину и паролю |
| POST | `/auth/refresh` | Обновление по refresh-токену |
| POST | `/auth/logout` | Отзыв сессии по refresh-токену |
| POST | `/tasks` | JWT, создание своей задачи |
| GET | `/tasks` | JWT, список своих задач |
| GET | `/tasks/{id}` | JWT, чтение своей задачи |
| PATCH | `/tasks/{id}` | JWT, изменение своей задачи |
| DELETE | `/tasks/{id}` | JWT, удаление своей задачи |
| POST | `/users` | Создание профиля без логина и пароля; сейчас публичный |
| GET | `/users` | Список профилей; сейчас публичный |
| GET | `/users/{id}` | Чтение профиля; сейчас публичный |
| PATCH | `/users/{id}` | Изменение профиля; сейчас публичный |
| DELETE | `/users/{id}` | Удаление профиля; сейчас публичный |
| GET | `/statistics` | Статистика; сейчас публичный |

Статистика поддерживает `user_id`, `from` и `to`. Даты передаются как `YYYY-MM-DD` и фильтруют время создания задач: `from` включительно, `to` исключительно. Без `user_id` вычисляется общая статистика.

Типичные статусы: `400` — некорректные данные, `401` — неуспешная аутентификация, `404` — ресурс не найден или чужая задача, `409` — конфликт логина или версии задачи, `500` — внутренняя ошибка. Полный контракт находится в Swagger.

## Swagger

Аннотации находятся рядом с HTTP-хендлерами. Для обновления файлов в `docs/`:

```sh
make swagger-gen
```

Затем перезапусти приложение. В Swagger UI нажми **Authorize** и введи `Bearer <access_token>` целиком. Refresh-токен передаётся в JSON соответствующих запросов.

Если появляется **Failed to fetch**, проверь Request URL, доступность сервера и адрес Swagger UI. Сейчас в аннотации `@host` закреплён `127.0.0.1:5050`, а CORS разрешает `http://localhost:5050`. При другом порте, HTTPS или удалённом доступе адрес API и CORS нужно согласовать с окружением. Для локальной проверки можно открыть UI на `http://127.0.0.1:5050/swagger/index.html`, чтобы UI и API использовали один origin. Само сообщение не доказывает, что причина именно в CORS.

## Структура проекта

```text
cmd/todoapp/                 Точка входа, сборка зависимостей, Dockerfile
internal/core/              Домен, конфигурация, логирование, HTTP и пул БД
internal/features/
  auth/                     Регистрация, JWT, refresh-сессии и logout
  tasks/                    Личные задачи
  users/                    Профили пользователей
  statistics/               Статистика выполнения
migrations/                 SQL-миграции up/down
docs/                       Сгенерированные Swagger Go/JSON/YAML
```

Внутри фич HTTP-слой преобразует запросы и ответы, сервис выполняет прикладную логику, репозиторий работает с PostgreSQL. Интерфейсы зависимостей объявляются в использующем их слое. Ротация refresh-токенов выполняется в транзакции с блокировкой строки сессии.

## Команды разработки

| Команда | Действие |
| --- | --- |
| `make env-up` | Запустить PostgreSQL |
| `make env-port-forward` | Открыть БД на локальном порту 5432 |
| `make migrate-up` | Применить миграции |
| `make migrate-create seq=change_name` | Создать пару файлов миграции |
| `make todoapp-run` | Запустить Go-приложение локально |
| `make swagger-gen` | Перегенерировать Swagger |
| `make ps` | Посмотреть состояние Compose-сервисов |
| `go test ./...` | Запустить существующие тесты и проверить сборку пакетов |
| `go vet ./...` | Выполнить статический анализ |

`make migrate-down` запускает откат миграций, а не обязательно только последней. Для одного шага используй `make migrate-action action="down 1"`.

`make env-cleanup` удаляет локальные данные PostgreSQL после подтверждения. `make logs-cleanup` удаляет логи. Эти команды не нужны для обычного перезапуска.

## Запуск приложения в Docker

В текущем `docker-compose.yaml` сервис `todoapp` ещё не получает `JWT_SECRET`. Перед запуском добавь в его `environment`:

```yaml
- JWT_SECRET=${JWT_SECRET}
```

После настройки `.env`, запуска БД и применения миграций:

```sh
make todoapp-deploy
```

Порт контейнера опубликован как `5050:5050`. Для остановки приложения:

```sh
make todoapp-undeploy
```

Локальные логи пишутся в `out/logs`, контейнерные согласно текущему mount — в `app/out/logs`. База хранит данные в `out/pgdata`.

## Что ещё нужно до публичного развёртывания

- Ограничить доступ к профилям пользователей и статистике.
- Добавить ограничение частоты входа и регистрацию попыток; выровнять проверку пароля для неизвестного логина.
- Настроить HTTPS и способ хранения токенов браузерным клиентом. Сейчас оба токена возвращаются в JSON, HttpOnly cookie не используется.
- Проверить ротацию и отзыв сессий интеграционными тестами с PostgreSQL, включая конкурентные запросы.
- Настроить очистку истёкших сессий и хранение секретов для окружения развёртывания.
