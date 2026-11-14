# Tasks API

Простой HTTP API для создания и управления задачами. Данные хранятся в памяти процесса: после остановки или перезапуска сервера список задач очищается.

## Структура проекта

```text
.
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа: сборка хранилища и HTTP-сервера
├── internal/
│   ├── api/
│   │   ├── dto.go                 # DTO запросов и формат JSON-ошибок
│   │   ├── handlers.go            # Обработчики HTTP-эндпоинтов
│   │   ├── middleware.go          # Логирование HTTP-запросов
│   │   └── server.go              # Настройка маршрутов и запуск сервера
│   ├── models/
│   │   └── task.go                # Модель задачи и ошибки домена
│   └── storage/
│       ├── storage.go             # Интерфейс хранилища
│       └── memory/
│           └── memory.go          # In-memory хранилище задач
├── postman/
│   └── Tasks API.postman_collection.json # Коллекция запросов Postman
├── go.mod                         # Модуль и зависимости Go
├── go.sum                         # Контрольные суммы зависимостей
└── README.md                      # Документация проекта
```

## Postman

Коллекция запросов находится в [`postman/Tasks API.postman_collection.json`](./postman/Tasks%20API.postman_collection.json). Она использует формат Postman Collection v2.1 и содержит запросы для проверки состояния сервера, создания, просмотра списка, получения, обновления и удаления задач.

Чтобы начать работу:

1. Запустите API командой `go run ./cmd/server`.
2. В Postman выберите **Import** и укажите файл `postman/Tasks API.postman_collection.json`.
3. При необходимости измените переменную коллекции `baseUrl` (по умолчанию `http://localhost:9091`).
4. Выполните запрос **Create task**. После успешного создания ID сохраняется в переменную `taskId`.
5. Выполняйте **Get task**, **Update task** и **Delete task** для сохраненного ID. **List tasks** показывает все существующие задачи.

Переменные коллекции можно посмотреть в меню коллекции **Edit → Variables**. После перезапуска API задачи удаляются из памяти; создайте новую задачу перед запросами с `taskId`.

## Запуск

Требуется Go версии, совместимой с `go.mod` (указана `1.27.4`).

Из корня проекта запустите сервер:

```bash
go run ./cmd/server
```

Сервер слушает порт `9091`. Остановить его можно сочетанием `Ctrl+C`.

Проверка доступности:

```bash
curl -i http://localhost:9091/health
```

Пример ответа:

```http
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

## Эндпоинты

Все ответы с телом используют JSON. Для запросов `POST /tasks` и `PUT /tasks/{id}` передавайте `Content-Type: application/json`.

| Метод | Путь | Назначение | Успешный статус |
|---|---|---|---|
| `GET` | `/health` | Проверить состояние сервера | `200 OK` |
| `POST` | `/tasks` | Создать задачу | `201 Created` |
| `GET` | `/tasks` | Получить список задач | `200 OK` |
| `GET` | `/tasks/{id}` | Получить задачу по ID | `200 OK` |
| `PUT` | `/tasks/{id}` | Полностью обновить задачу | `200 OK` |
| `DELETE` | `/tasks/{id}` | Удалить задачу | `204 No Content` |

### Модель задачи

```json
{
  "id": 1,
  "title": "Подготовить отчет",
  "description": "Собрать результаты за неделю",
  "done": false,
  "created_at": "2026-11-14T00:00:00Z"
}
```

`id` назначается сервером, `done` при создании изначально равен `false`, `created_at` содержит время создания в UTC. `description` — необязательное дополнительное поле. Время в примере условное.

### Создать задачу

`POST /tasks` требует поле `title`, содержащее хотя бы один символ, отличный от пробельных. Поле `description` можно опустить.

```bash
curl -i -X POST http://localhost:9091/tasks \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Подготовить отчет",
    "description": "Собрать результаты за неделю"
  }'
```

Пример ответа (`201 Created`):

```json
{
  "id": 1,
  "title": "Подготовить отчет",
  "description": "Собрать результаты за неделю",
  "done": false,
  "created_at": "2026-11-14T00:00:00Z"
}
```

Ответ также содержит заголовок `Location: /tasks/{id}` с адресом созданной задачи.

### Получить список задач

```bash
curl -i http://localhost:9091/tasks
```

Пример ответа (`200 OK`):

```json
[
  {
    "id": 1,
    "title": "Подготовить отчет",
    "description": "Собрать результаты за неделю",
    "done": false,
    "created_at": "2026-11-14T00:00:00Z"
  }
]
```

Если задач нет, возвращается пустой массив `[]`.

### Получить задачу

```bash
curl -i http://localhost:9091/tasks/1
```

Пример ответа (`200 OK`):

```json
{
  "id": 1,
  "title": "Подготовить отчет",
  "description": "Собрать результаты за неделю",
  "done": false,
  "created_at": "2026-11-14T00:00:00Z"
}
```

### Обновить задачу

`PUT /tasks/{id}` заменяет значения полей задачи. Обязательны `title`, содержащий хотя бы один символ, отличный от пробельных, и поле `done` типа JSON boolean. `description` необязательно; если его опустить, значение будет заменено пустой строкой.

```bash
curl -i -X PUT http://localhost:9091/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Подготовить отчет",
    "description": "Отчет готов и отправлен",
    "done": true
  }'
```

Пример ответа (`200 OK`):

```json
{
  "id": 1,
  "title": "Подготовить отчет",
  "description": "Отчет готов и отправлен",
  "done": true,
  "created_at": "2026-11-14T00:00:00Z"
}
```

### Удалить задачу

```bash
curl -i -X DELETE http://localhost:9091/tasks/1
```

Успешный ответ: `204 No Content`, без тела.

## Ошибки

Все ошибки API, включая ошибки маршрутизации, возвращаются с `Content-Type: application/json` и единым объектом:

```json
{
  "error": "<текст ошибки>",
  "time": "<время ошибки в формате ISO8601>"
}
```

Поле `time` указывает время формирования ответа об ошибке и сериализуется в формате RFC3339 (ISO 8601). Для запросов `POST` и `PUT` тело должно содержать ровно одно JSON-значение; дополнительные значения или данные после него приводят к `400 Bad Request`. Ниже приведены ошибочные сценарии для ручной проверки.

### Проверка ошибок по эндпоинтам

`GET /health`: неподдерживаемый метод возвращает `405 Method Not Allowed`:

```bash
curl -i -X POST http://localhost:9091/health
```

`POST /tasks`: некорректный JSON и отсутствие обязательного `title` возвращают `400 Bad Request`:

```bash
curl -i -X POST http://localhost:9091/tasks \
  -H 'Content-Type: application/json' \
  -d '{'
```

```bash
curl -i -X POST http://localhost:9091/tasks \
  -H 'Content-Type: application/json' \
  -d '{"description":"Описание без заголовка"}'
```

В первом случае `error` содержит ошибку разбора JSON, во втором — `title is empty`. В обоих случаях тело имеет общий формат выше.

`GET /tasks`: неподдерживаемый метод возвращает `405 Method Not Allowed`:

```bash
curl -i -X DELETE http://localhost:9091/tasks
```

`GET /tasks/{id}`: нечисловой ID возвращает `400 Bad Request` с ошибкой `invalid task id`, а числовой ID, которого нет в хранилище, — `404 Not Found`:

```bash
curl -i http://localhost:9091/tasks/abc
curl -i http://localhost:9091/tasks/999999
```

`PUT /tasks/{id}`: некорректный JSON, пропущенный обязательный `title` или `done`, а также нечисловой ID возвращают `400 Bad Request`. Числовой ID отсутствующей задачи возвращает `404 Not Found`.

```bash
curl -i -X PUT http://localhost:9091/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{'
curl -i -X PUT http://localhost:9091/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"description":"Без заголовка","done":true}'
curl -i -X PUT http://localhost:9091/tasks/abc \
  -H 'Content-Type: application/json' \
  -d '{"title":"Задача","done":true}'
curl -i -X PUT http://localhost:9091/tasks/999999 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Задача","done":true}'
```

Для отсутствующего `title` значение `error` — `title cannot be empty`, для нечислового ID — `invalid task id`, для отсутствующей задачи — `task not found`.

`DELETE /tasks/{id}`: нечисловой ID возвращает `400 Bad Request`, отсутствующая задача — `404 Not Found`:

```bash
curl -i -X DELETE http://localhost:9091/tasks/abc
curl -i -X DELETE http://localhost:9091/tasks/999999
```

Запрос к неизвестному пути возвращает JSON `404 Not Found`; запрос к существующему пути с неподдерживаемым методом возвращает JSON `405 Method Not Allowed` и заголовок `Allow` со списком поддерживаемых методов этого пути. Например:

```bash
curl -i http://localhost:9091/unknown
curl -i -X PATCH http://localhost:9091/tasks/1
```

В первом случае значение `error` — `route not found`, во втором — `method not allowed`.
