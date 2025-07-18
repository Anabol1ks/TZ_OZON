# TZ_OZON — Мок-сервис курсов валют (ЦБ РФ)

Этот сервис эмулирует поведение [ЦБ РФ API](https://www.cbr.ru/scripts/XML_daily.asp?date_req=02/03/2002) по возврату курсов валют на указанную дату. 

---

## 🚀 Возможности

- Эмуляция успешного ответа от ЦБ РФ (код 200)
- Возможность симулировать ошибку (код 500)
- Поддержка Swagger-документации
- Генерация уникального XML-ответа для каждой даты
- Сохранение XML в базу данных PostgreSQL
- Поддержка proto-схемы (gRPC-интерфейс — в заделе на будущее)
- Makefile для сборки, генерации protobuf и запуска e2e тестов
- Юнит-тесты для основного хендлера

---

## 🔧 Технологии

- Go + Gin
- PostgreSQL + GORM
- Docker + Docker Compose
- Swagger (swaggo/gin-swagger)
- Protobuf
- Zap logger
- Testify для юнит-тестов

---


## ⚙️ Запуск

### 1. Клонируй репозиторий

```bash
git clone https://github.com/Anabol1ks/TZ_OZON.git
cd TZ_OZON
```

### 2. Создай .env файл

Скопируй `.env.example` в `.env` и при необходимости измени параметры:

```bash
copy .env.example .env  # Windows
```

### 3. Запуск через Docker Compose

```bash
make docker-up
# или
docker-compose -f docker-compose.yml up -d
```

### 4. Локальный запуск (без Docker)

```bash
make run
# или
go run cmd/main.go
```

### 5. Swagger доступен по адресу:

http://localhost:8080/swagger/index.html

---

## 📌 Примеры использования

### Успешный запрос

```http
GET /scripts/XML_daily.asp?date_req=17/07/2025
```

Ответ: XML-документ с курсами валют

### Возможная ошибка

Иногда при запросе новой даты, сервис может с вероятностью 50% вернуть 500 — для эмуляции нестабильного поведения стороннего API.

---


## 📁 Makefile команды

```bash
make run           # запуск приложения локально
make docker-up     # запуск через docker-compose
make generate_proto # генерация protobuf-файлов
make test          # запуск тестов (только handler)
make e2e           # запуск e2e тестов (параллельно, требуется запущенный сервер)
```

---


## 🧪 Тесты

### Юнит-тесты

Юнит-тесты для `GET /scripts/XML_daily.asp` доступны и покрывают:

* Успешный сценарий (200)
* Отсутствие параметра `date_req` (400)
* Некорректный формат даты (400)

```bash
make test
# или
go test internal/handler/exchange_test.go
```

### E2E тесты

E2E тесты имитируют параллельные запросы к сервису с разными датами и проверяют оба сценария: успешный ответ (200) и ошибку (500).

> **Важно:** Для запуска e2e тестов сервер должен быть запущен (`make run` или через Docker Compose), иначе тесты не смогут подключиться к API.

Запуск e2e тестов:

```bash
make e2e
# или
go test ./test/e2e/ -parallel 20 -count=1
```

В логах сервиса вы увидите множество параллельных запросов с разными датами.

---

## 🛰 Структура проекта

```
.
├── cmd/                  # main.go
├── internal/
│   ├── config/           # конфигурация из .env
│   ├── db/               # подключение и миграции PostgreSQL
│   ├── handler/          # обработчики Gin
│   ├── logger/           # инициализация zap логгера
│   ├── models/           # GORM-модели
│   └── router/           # маршруты + swagger
├── proto/                # proto-файлы (gRPC)
├── tz_ozon/              # сгенерированные protobuf-файлы
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## 📚 Пример gRPC интерфейса (ExchangeService)

> В задел на будущее, предусмотрен proto-файл `proto/exchange.proto`

```protobuf
service ExchangeService {
  rpc GetExchangeRate (ExchangeRateRequest) returns (ExchangeRateResponse);
}
```

---


## 🧾 .env файл (пример)

Пример содержимого `.env.example`:

```env
DB_HOST=ozon-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=ozon12341
DB_NAME=ozon-db
DB_SSLMODE=disable
```

---

## ✅ Выполнены требования задания

* [x] Уникальные данные под каждый тест
* [x] Поддержка кодов 200 и 500
* [x] Без использования заголовков
* [x] Makefile
* [x] Protobuf файл
* [x] Swagger
* [x] БД
* [x] Go

---

## 📎 Автор
[Anabol1ks](https://github.com/Anabol1ks)

