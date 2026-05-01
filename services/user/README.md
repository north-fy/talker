# User Service

Сервис аутентификации и управления пользователями.

## 🎯 Функционал

### Регистрация и вход
- `Register` — регистрация нового пользователя
- `Login` — вход и получение JWT токена
- `GetMe` — получение информации о текущем пользователе
- `ValidateToken` — валидация JWT токена (для других сервисов)

## 🗄️ База данных

### Таблицы

```sql
-- users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);