# Message Service

Сервис управления сообщениями, реакциями и WebSocket соединениями.

## 🎯 Функционал

### CRUD сообщений
- `SendMessage` — отправка сообщения
- `GetMessages` — получение сообщений чата (пагинация)
- `GetMessage` — получение одного сообщения
- `EditMessage` — редактирование сообщения
- `DeleteMessage` — удаление сообщения

### Реакции
- `AddReaction` — добавление реакции (❤️ 👍 😂)
- `RemoveReaction` — удаление реакции

### Поиск
- `SearchMessages` — поиск по сообщениям

### Непрочитанные
- `MarkAsRead` — отметка о прочтении
- `GetUnreadCount` — количество непрочитанных

### WebSocket
- `ConnectWebSocket` — стриминг сообщений в реальном времени

### Внутренние методы
- `GetLastMessage` — последнее сообщение чата (для Chat Service)
- `DeleteChatMessages` — удаление всех сообщений чата

## 🗄️ База данных

### Таблицы

```sql
-- messages (партиционирование по created_at)
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    content TEXT,
    type VARCHAR(50) DEFAULT 'text',
    reply_to UUID,
    attachments TEXT[],
    reactions JSONB DEFAULT '{}',
    is_edited BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
) PARTITION BY RANGE (created_at);

-- read_receipts
CREATE TABLE read_receipts (
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    last_read_message_id UUID,
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);