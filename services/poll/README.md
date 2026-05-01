# Poll Service

Сервис создания опросов и голосований в чатах.

## 🎯 Функционал

### Управление опросами
- `CreatePoll` — создание опроса
- `GetPoll` — получение опроса с результатами
- `UpdatePoll` — обновление опроса
- `DeletePoll` — удаление опроса
- `GetChatPolls` — список опросов чата

### Голосование
- `Vote` — голосование (один или несколько вариантов)
- `Unvote` — отмена голоса
- `GetResults` — получение результатов с процентами
- `GetMyVote` — получение своего голоса

### Закрытие/открытие
- `ClosePoll` — закрытие опроса
- `ReopenPoll` — открытие опроса

### Внутренние методы
- `GetPollInternal` — получение опроса (для Message Service)

## 🗄️ База данных

### Таблицы

```sql
-- polls
CREATE TABLE polls (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL,
    created_by UUID NOT NULL,
    question TEXT NOT NULL,
    settings JSONB,
    status VARCHAR(20) DEFAULT 'active',
    total_votes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    closed_at TIMESTAMP
);

-- poll_options
CREATE TABLE poll_options (
    id UUID PRIMARY KEY,
    poll_id UUID REFERENCES polls(id),
    text VARCHAR(500) NOT NULL,
    votes_count INT DEFAULT 0
);

-- poll_votes
CREATE TABLE poll_votes (
    poll_id UUID REFERENCES polls(id),
    user_id UUID NOT NULL,
    option_ids TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(poll_id, user_id)
);