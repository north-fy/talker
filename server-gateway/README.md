# Server Gateway

REST API шлюз для gRPC-микросервисов Talker на **Gin**.

## 🎯 Назначение

Единая точка входа для клиентов: принимает REST-запросы (`/api/v1`), проксирует их в gRPC-сервисы и маппит gRPC-ошибки в HTTP-статусы.

Бэкенды (gRPC):

| Сервис | Адрес |
|---|---|
| User | `:50051` |
| Message | `:50052` |
| Chat | `:50053` |

## 🔐 Аутентификация

Шлюз валидирует `Authorization: Bearer <JWT>` локально (дешифрует `JWT_SECRET`) и пробрасывает токен в gRPC metadata для сервисов. `user_id` извлекается из JWT и подставляется во внутренние запросы.

Публичные эндпоинты: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`. Остальные требуют валидный JWT.

## 📦 Эндпоинты

### Auth
| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/auth/register` | Регистрация (public) |
| POST | `/api/v1/auth/login` | Вход, получение JWT (public) |
| GET | `/api/v1/auth/me` | Текущий пользователь |
| POST | `/api/v1/auth/validate` | Валидация токена |

### Chats
| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/chats` | Создать чат |
| GET | `/api/v1/chats` | Список чатов (query: `type`, `search`, `include_archived`) |
| GET | `/api/v1/chats/:chat_id` | Получить чат |
| PATCH | `/api/v1/chats/:chat_id` | Обновить чат |
| GET | `/api/v1/users/me/chats` | Чаты текущего пользователя |

### Members
| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/chats/:chat_id/members` | Добавить участника |
| GET | `/api/v1/chats/:chat_id/members` | Участники чата (query: `role`, `search`) |
| GET | `/api/v1/chats/:chat_id/members/:user_id` | Получить участника |
| PATCH | `/api/v1/chats/:chat_id/members/:user_id/role` | Изменить роль |
| DELETE | `/api/v1/chats/:chat_id/members/:user_id` | Удалить участника |

### Invites
| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/chats/:chat_id/invites` | Создать ссылку-приглашение |
| DELETE | `/api/v1/chats/:chat_id/invites/:invite_id` | Отозвать ссылку |
| POST | `/api/v1/invites/join` | Вступить по ссылке |

### Messages
| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/chats/:chat_id/messages` | Отправить сообщение |
| GET | `/api/v1/chats/:chat_id/messages` | Сообщения чата (query: `limit`, `before`, `after`) |
| GET | `/api/v1/messages/search` | Поиск (query: `query`, `chat_id`, `limit`, `before`) |
| GET | `/api/v1/messages/:message_id` | Получить сообщение |
| PATCH | `/api/v1/messages/:message_id` | Изменить сообщение |
| DELETE | `/api/v1/messages/:message_id` | Удалить сообщение |
| POST | `/api/v1/messages/:message_id/reactions` | Добавить реакцию |
| DELETE | `/api/v1/messages/:message_id/reactions/:reaction` | Удалить реакцию |
| POST | `/api/v1/chats/:chat_id/read` | Отметить прочитанным |
| GET | `/api/v1/chats/:chat_id/unread-count` | Непрочитанные |

### Прочее
| Метод | Путь | Назначение |
|---|---|---|
| GET | `/healthz` | Проверка живости |

## ⚠️ Ошибки

Единый формат ошибок:
```json
{ "error": "описание", "code": 404 }
```

Маппинг gRPC → HTTP: `Unauthenticated→401`, `PermissionDenied→403`, `NotFound→404`, `InvalidArgument/FailedPrecondition→400`, `AlreadyExists→409`, `DeadlineExceeded→504`, `Unavailable/Aborted→503`.

## 🚀 Запуск

### Локально

```bash
cp .env.example .env
make run
```

Требуются запущенные user/message/chat сервисы и заданный `JWT_SECRET` (должен совпадать с user-сервисом).

### Docker

```bash
make up
```

## 🔧 Makefile

- `make up/down/rebuild/logs` — Docker
- `make run` — локальный запуск
- `make build/test/vet/fmt` — проверки