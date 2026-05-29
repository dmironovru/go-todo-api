# 🚀 Todo Manager — Backend (Go + PostgreSQL)

REST API для менеджера задач с ограничением на 10 задач.

## Демо
[https://dmitrymironov.ru/lab/](https://dmitrymironov.ru/lab/)

## Технологии
- Go 1.25+
- PostgreSQL 16
- pgx (драйвер)
- systemd

## API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/tasks` | Получить все задачи |
| POST | `/api/tasks` | Создать задачу (макс 10) |
| PUT | `/api/tasks/{id}` | Обновить задачу (title/completed) |
| DELETE | `/api/tasks/{id}` | Удалить задачу |
| DELETE | `/api/tasks` | Очистить все задачи |

## Запуск локально

### Требования
- PostgreSQL 16+
- Go 1.25+

### Настройка БД
```sql
CREATE DATABASE todo_db;
CREATE USER todo_user WITH PASSWORD 'your_strong_password';
GRANT ALL PRIVILEGES ON DATABASE todo_db TO todo_user;
