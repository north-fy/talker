# Chat Service

Сервис управления чатами, участниками и приглашениями.

## 🎯 Функционал

### Управление чатами
- `CreateChat` — создание нового чата (создатель становится администратором)
- `GetChat` — получение информации о чате
- `GetChats` — список чатов с фильтрацией и пагинацией
- `UpdateChat` — обновление чата (название, аватар)
- `DeleteChat` — удаление чата (с удалением всех сообщений через Message Service)

### Управление участниками
- `AddMember` — добавление участника в чат
- `RemoveMember` — удаление участника из чата
- `GetMembers` — список участников с пагинацией
- `UpdateMemberRole` — изменение роли участника (owner/admin/member)
- `GetMember` — получение информации об участнике

### Приглашения
- `CreateInviteLink` — создание ссылки-приглашения
- `JoinChatByInvite` — вступление в чат по ссылке
- `RevokeInviteLink` — отзыв ссылки-приглашения

### Внутренние методы (для User / Message сервисов)
- `IsMember` — проверка, является ли пользователь участником
- `GetUserChats` — получение всех чатов пользователя
- `GetChatInternal` — внутреннее получение чата
- `ValidateMemberAccess` — проверка прав доступа участника

## 🗄️ База данных

### Таблицы

```sql
-- chats
CREATE TABLE chats (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    type VARCHAR(50) NOT NULL,  -- private/group/channel
    created_by UUID NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- chat_members
CREATE TABLE chat_members (
    chat_id UUID REFERENCES chats(id),
    user_id UUID NOT NULL,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),
    last_read_at TIMESTAMP,
    PRIMARY KEY (chat_id, user_id)
);

-- invite_links
CREATE TABLE invite_links (
    id UUID PRIMARY KEY,
    chat_id UUID REFERENCES chats(id),
    code VARCHAR(100) UNIQUE NOT NULL,
    max_uses INT DEFAULT 0,
    used_count INT DEFAULT 0,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID NOT NULL,
    is_active BOOLEAN DEFAULT true
);
```

## 🧠 Кэширование (Redis)

- `GetChat` кэширует данные чата (ключ `chat:{id}`, TTL 10 мин)
- `GetChatInternal` использует тот же кэш
- `IsMember` кэширует membership-статус (ключ `member:{chat}:{user}`, TTL 5 мин)
- Инвалидация при обновлении чата, добавлении/удалении участников

## 🚀 Запуск

### Docker

```bash
docker-compose up -d --build
```

### Локально

```bash
cp .env.example .env
go mod download
go run ./cmd/server
```

Требуются запущенные Postgres (порт 5433), Redis (порт 6380) и сервисы user (50051), message (50052).

## 🔧 Makefile

- `make up` — запуск в Docker
- `make down` — остановка
- `make rebuild` — пересборка и запуск
- `make logs` — логи
- `make migrations` — применение миграций
