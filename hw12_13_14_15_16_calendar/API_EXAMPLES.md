# Calendar API Examples

Примеры использования Calendar API.

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Health Check

Проверка работоспособности сервиса.

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "time": "2026-01-13T14:30:00Z"
}
```

---

### 2. Get All Events

Получить список всех событий.

```bash
curl http://localhost:8080/api/events
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly sync",
    "start_time": "2026-01-15T10:00:00Z",
    "end_time": "2026-01-15T11:00:00Z",
    "user_id": "user1",
    "notify_before": 900
  }
]
```

---

### 3. Create Event

Создать новое событие.

```bash
curl -X POST http://localhost:8080/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly sync meeting",
    "start_time": "2026-01-15T10:00:00Z",
    "end_time": "2026-01-15T11:00:00Z",
    "user_id": "user1",
    "notify_before": 900
  }'
```

**Request Body:**
```json
{
  "title": "Team Meeting",
  "description": "Weekly sync meeting",
  "start_time": "2026-01-15T10:00:00Z",
  "end_time": "2026-01-15T11:00:00Z",
  "user_id": "user1",
  "notify_before": 900
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly sync meeting",
  "start_time": "2026-01-15T10:00:00Z",
  "end_time": "2026-01-15T11:00:00Z",
  "user_id": "user1",
  "notify_before": 900,
  "created_at": "2026-01-13T14:30:00Z"
}
```

---

### 4. Get Event by ID

Получить событие по ID.

```bash
curl http://localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly sync meeting",
  "start_time": "2026-01-15T10:00:00Z",
  "end_time": "2026-01-15T11:00:00Z",
  "user_id": "user1",
  "notify_before": 900
}
```

---

### 5. Update Event

Обновить существующее событие.

```bash
curl -X PUT http://localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Team Meeting",
    "description": "Updated description",
    "start_time": "2026-01-15T14:00:00Z",
    "end_time": "2026-01-15T15:00:00Z",
    "user_id": "user1",
    "notify_before": 1800
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated Team Meeting",
  "description": "Updated description",
  "start_time": "2026-01-15T14:00:00Z",
  "end_time": "2026-01-15T15:00:00Z",
  "user_id": "user1",
  "notify_before": 1800
}
```

---

### 6. Delete Event

Удалить событие.

```bash
curl -X DELETE http://localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "message": "Event deleted successfully"
}
```

---

### 7. Search Events by Date

Найти события на конкретную дату.

```bash
curl "http://localhost:8080/api/events?date=2026-01-15"
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "start_time": "2026-01-15T10:00:00Z",
    "end_time": "2026-01-15T11:00:00Z"
  }
]
```

---

### 8. Search Events by Date Range

Найти события в диапазоне дат.

```bash
curl "http://localhost:8080/api/events?start=2026-01-14&end=2026-01-17"
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "start_time": "2026-01-15T10:00:00Z",
    "end_time": "2026-01-15T11:00:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "title": "Project Review",
    "start_time": "2026-01-16T14:00:00Z",
    "end_time": "2026-01-16T15:00:00Z"
  }
]
```

---

### 9. Search Events by User

Найти события конкретного пользователя.

```bash
curl "http://localhost:8080/api/events?user_id=user1"
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "user_id": "user1",
    "start_time": "2026-01-15T10:00:00Z"
  }
]
```

---

## Error Responses

### 400 Bad Request

Неверный формат запроса.

```json
{
  "error": "Invalid request body",
  "details": "title is required"
}
```

### 404 Not Found

Событие не найдено.

```json
{
  "error": "Event not found",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 500 Internal Server Error

Внутренняя ошибка сервера.

```json
{
  "error": "Internal server error",
  "message": "Database connection failed"
}
```

---

## Using with HTTPie

Более удобный способ тестирования API с помощью [HTTPie](https://httpie.io/):

```bash
# Get all events
http GET localhost:8080/api/events

# Create event
http POST localhost:8080/api/events \
  title="Meeting" \
  description="Team sync" \
  start_time="2026-01-15T10:00:00Z" \
  end_time="2026-01-15T11:00:00Z" \
  user_id="user1" \
  notify_before:=900

# Update event
http PUT localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440000 \
  title="Updated Meeting" \
  start_time="2026-01-15T14:00:00Z" \
  end_time="2026-01-15T15:00:00Z"

# Delete event
http DELETE localhost:8080/api/events/550e8400-e29b-41d4-a716-446655440000
```

---

## Using with Postman

1. Импортируйте OpenAPI спецификацию: `http://localhost:8080/openapi.yaml`
2. Postman автоматически создаст коллекцию запросов
3. Настройте переменную окружения `base_url = http://localhost:8080`

---

## Testing Script

Используйте готовый скрипт для тестирования всех endpoints:

```bash
bash test-api.sh
```

Скрипт автоматически:
- Проверит health check
- Создаст событие
- Обновит событие
- Получит список событий
- Удалит событие
