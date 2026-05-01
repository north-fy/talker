# Chat Service

Сервис управления чатами, участниками и приглашениями.

## 🎯 Функционал

### Управление чатами
- `CreateChat` — создание нового чата
- `GetChat` — получение информации о чате
- `GetChats` — список чатов с фильтрацией
- `UpdateChat` — обновление чата (название, аватар)

### Управление участниками
- `AddMember` — добавление участника в чат
- `RemoveMember` — удаление участника из чата
- `GetMembers` — список участников с пагинацией
- `UpdateMemberRole` — изменение роли участника
- `GetMember` — получение информации об участнике

### Проверки для других сервисов
- `IsMember` — проверка, является ли пользователь участником
- `GetUserChats` — получение всех чатов пользователя
- `GetChatInternal` — внутреннее получение чата
- `ValidateMemberAccess` — проверка прав доступа

### Приглашения
- `CreateInviteLink` — создание ссылки-приглашения
- `JoinChatByInvite` — вступление в чат по ссылке
- `RevokeInviteLink` — отзыв ссылки-приглашения

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