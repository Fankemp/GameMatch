# GameMatch

**GameMatch** — это приложение для геймеров, которые ищут тимейтов для совместной игры.

Пользователи могут создавать карточки для разных игр, свайпать профили других игроков и находить людей с похожими интересами.

## Основной flow

1. Пользователь регистрируется.
2. Создаёт карточки для игр (например, для Valorant указывает ранг, роль).
3. Видит фид других игроков и свайпает **like/dislike**.
4. Если оба пользователя свайпнули друг друга (**bilateral match**) → они видят контакты друг друга.
5. Можно писать друг другу в игре (чат внутри приложения пока не нужен).

---

## Высоконагруженные части

- Фид игроков должен быстро генерироваться (много фильтров).
- Свайпы должны обрабатываться мгновенно.
- Уведомления отправляются асинхронно через воркеры (не блокируют основной API).

---

## Архитектура

- **Monolith API (Gin)** — обрабатывает основные действия: свайпы, профили, матчи.
- **Workers (горутины)** — слушают Redis Streams и отправляют push-уведомления.
- **PostgreSQL** — основная база данных для хранения пользователей, карточек и матчей.
- **Redis** — кеширует фид и используется как очередь событий.

---

## Технологии

- **Backend:** Go (Gin)
- **Database:** PostgreSQL
- **Cache / Queue:** Redis
- **Workers:** Go горутины для асинхронных задач
- **CI/CD:** (будет добавлено позже через GitHub Actions)

---

## Структура проекта (начальная)
Не финальный вариант, возможны улучшение


```text
gamematch/
├─ cmd/                          ← Entry points (где запускается код)
│  ├─ api/
│  │  └─ main.go                 ← Запуск REST API сервера
│  └─ worker/
│     └─ main.go                 ← Запуск async workers
│
├─ internal/                      ← Приватный код (не экспортируется)
│  ├─ config/
│  │  └─ config.go               ← Загрузка .env и переменных
│  │
│  ├─ handler/                   ← HTTP handlers (контроллеры)
│  │  ├─ auth.go                 ← /auth endpoints
│  │  ├─ card.go                 ← /cards endpoints
│  │  ├─ swipe.go                ← /swipes endpoints
│  │  ├─ feed.go                 ← /feed endpoints
│  │  ├─ match.go                ← /matches endpoints
│  │  └─ middleware.go            ← JWT auth, logging, CORS
│  │
│  ├─ service/                   ← Бизнес-логика
│  │  ├─ auth_service.go         ← Регистрация, логин
│  │  ├─ card_service.go         ← CRUD карточек
│  │  ├─ swipe_service.go        ← Обработка свайпов, матчинг
│  │  ├─ feed_service.go         ← Генерация фида
│  │  └─ notification_service.go ← Отправка уведомлений
│  │
│  ├─ repository/                ← Доступ к БД (queries)
│  │  ├─ user_repo.go            ← SELECT/INSERT users
│  │  ├─ card_repo.go            ← SELECT/INSERT cards
│  │  ├─ swipe_repo.go           ← SELECT/INSERT swipes
│  │  └─ match_repo.go           ← SELECT/INSERT matches
│  │
│  ├─ model/                     ← Структуры данных
│  │  ├─ user.go
│  │  ├─ card.go
│  │  ├─ swipe.go
│  │  └─ match.go
│  │
│  ├─ worker/                    ← Async workers
│  │  ├─ notification_worker.go  ← Worker для push-уведомлений
│  │  └─ worker.go               ← Interface и helper'ы
│  │
│  ├─ db/                        ← Database
│  │  ├─ postgres.go             ← PostgreSQL connection
│  │  └─ migrations/             ← SQL migration файлы
│  │     ├─ 001_init_schema.up.sql
│  │     └─ 001_init_schema.down.sql
│  │
│  └─ redis/                     ← Redis
│     └─ redis.go                ← Redis connection & methods
│
├─ pkg/                          ← Переиспользуемые пакеты
│  └─ logger/
│     └─ logger.go               ← Структурированное логирование
│
├─ .env                          ← Переменные окружения (локально)
├─ .env.example                  ← Пример .env
├─ docker-compose.yml            ← PostgreSQL, Redis в Docker
├─ Dockerfile                    ← Образ приложения
├─ go.mod                        ← Зависимости
├─ go.sum
└─ README.md                     ← Документация