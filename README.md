# Мікросервіси Products і Notifications

Два мікросервіси на Golang для управління продуктами з асинхронними повідомленнями.

## 🏗️ Архітектура

- **Products Service** - REST API для управління товарами
- **Notifications Service** - слухає та логує події з Products
- **PostgreSQL** - база даних
- **RabbitMQ** - брокер повідомлень
- **Prometheus** - збір метрик

## 📋 Вимоги

- Docker
- Docker Compose

## 🚀 Швидкий старт

```bash
# Клонувати репозиторій
git clone <repository-url>
cd <project-folder>

# Запустити всі сервіси
docker-compose up --build
```

Сервіси будуть доступні:
- Products API: http://localhost:8080
- Prometheus: http://localhost:9090
- RabbitMQ Management: http://localhost:15672 (guest/guest)

## 📡 API Ендпоінти

### Створити продукт
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ноутбук",
    "description": "Ігровий ноутбук",
    "price": 35000
  }'
```

### Отримати список продуктів
```bash
curl "http://localhost:8080/api/products?page=1&page_size=10"
```

### Видалити продукт
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

### Переглянути метрики
```bash
curl http://localhost:8080/metrics
```

## 📊 Метрики Prometheus

- `products_created_total` - кількість створених товарів
- `products_deleted_total` - кількість видалених товарів

## 📝 Логи Notifications

Переглянути логи сервісу повідомлень:
```bash
docker-compose logs -f notifications
```

Приклад виводу:
```
✅ Product CREATED: ID=1, Name=Ноутбук
🗑️  Product DELETED: ID=1
```

## 🧪 Тестування

```bash
# Створити тестову БД
createdb products_test

# Запустити unit тести
cd products
go test ./internal/repository -v
```

## 📦 Структура проекту

```
.
├── docker-compose.yml
├── prometheus.yml
├── products/
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── broker/
│   │   └── metrics/
│   ├── migrations/
│   └── Dockerfile
└── notifications/
    ├── cmd/main.go
    ├── internal/
    │   └── consumer/
    └── Dockerfile
```

## 🛠️ Технології

- **Go 1.21**
- **Chi** - HTTP router
- **PostgreSQL 15** - база даних
- **RabbitMQ** - черги повідомлень
- **Prometheus** - моніторинг
- **golang-migrate** - міграції БД
- **lib/pq** - PostgreSQL драйвер

## 🔧 Конфігурація

Змінні середовища можна налаштувати в `docker-compose.yml`:

```yaml
environment:
  DATABASE_URL: "postgres://user:pass@postgres:5432/db?sslmode=disable"
  RABBITMQ_URL: "amqp://guest:guest@rabbitmq:5672/"
  PORT: "8080"
```

## 🛑 Зупинка сервісів

```bash
docker-compose down

# Видалити volumes (БД буде очищена)
docker-compose down -v
```

## 📖 Додаткова інформація

### Пагінація

- `page` - номер сторінки (за замовчуванням: 1)
- `page_size` - кількість елементів на сторінці (за замовчуванням: 10, max: 100)

### Валідація

- `name` - обов'язкове поле
- `price` - повинна бути більше 0

### Відповідь при помилці

```json
{
  "error": "опис помилки"
}
```
