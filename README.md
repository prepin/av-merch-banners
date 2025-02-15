# Наблюдения
* В схеме указано, покупка предмета за монеты происходит через `GET /api/buy/{item}`. Это вызывает вопросы, потому, что GET запрос не должен изменять данные на сервере, а только возвращать их. Решение: Делаем реализацию через POST запрос. Для GET делаем фоллбэк на случай, если уже существует клиент, который почему-то ориентируется на GET, но помечаем его как deprecated.
* Названия эндпойнтов тоже не совсем соответствуют REST (глаголы вместо существительных), но переделывать все ручки в рамках задачи не будет осмысленным, просто отметим этот факт.
* Названия полей в схеме должны быть в формате to_user, а не toUser, и линтеры на это тоже ругаются. Но оставим как есть, ибо ТЗ.
* API указано с префиксом просто `/api`, без версионирования. Это повлечёт за собой некоторую боль в будущем, когда понадобится реализовать изменения, ломающие совместимость. Решение: монтируем API по префиксу `/api/` и префиксу `/api/v1/`.

# Реализация
Сервис реализован на Go (фреймворк Gin). В качестве базы данных используется Postgres, в качестве кэша Redis.

Основное тестовое покрытие достигается за счёт E2E тестов, юнит-тесты используются для закрытия тех мест, куда E2E достать сложно (например ошибки БД).

Для транзакций перевода монеток между пользователями используется принцип двойной записи. На каждый перевод создаётся входящая и исходящая транзакции. Найти пару транзакций относящихся к одному переводу можно по ключу reference_id.

Для того, чтобы можно было начислять монеты, добавлен эндпойнт `POST /api/credit`, доступный только пользователям с ролью admin (определяется по полю `role` в таблице `users`).
Для начисления нужно указать добавляемую сумму в поле `amount`. Отрицательные суммы тоже можно указывать, но списан будет максимум тот баланс, который есть у пользователя.
```
POST /api/v1/credit
Authorization: Bearer your_admin_user_token

{
  "username": "employee",
  "amount": 1000
}
```
Админский пользователь  с логином `director` и паролем `password` создаётся при загрузке seed-данных.

# Нагрузочное тестирование
Проводилось с нагрузкой 1000 RPS в течение 60 секунд на каждый эндпойнт при помощи [oha](https://github.com/hatoo/oha). Тестовый стенд — докер, запущенный на Apple M1/16GB, без выставленных лимитов ресурсов на отдельные контейнеры.

Полные [отчёты о тестировании](docs/load_test_results) лежат в папке docs.

В SLI среднего времени ответа в 50ms укладываются все эндпойнты кроме /auth (там используется BCrypt и он упирается в процессор). Кэшировать его ответы по понятным причинам не стоит.

```
POST /sendCoin
  Success rate:	100.00%
  Total:	60.0010 secs
  Slowest:	3.1058 secs
  Fastest:	0.0011 secs
  Average:	0.0086 secs
  Requests/sec:	1000.0002
```
```
POST /buy
  Success rate:	100.00%
  Total:	60.0013 secs
  Slowest:	1.0138 secs
  Fastest:	0.0011 secs
  Average:	0.0030 secs
  Requests/sec:	999.9944
```
```
GET /info
  Success rate:	100.00%
  Total:	60.0019 secs
  Slowest:	1.3064 secs
  Fastest:	0.0003 secs
  Average:	0.0091 secs
  Requests/sec:	999.9679
```
```
POST /auth
  Success rate:	100.00%
  Total:	60.0157 secs
  Slowest:	54.6378 secs
  Fastest:	0.1885 secs
  Average:	27.4741 secs
  Requests/sec:	91.1762
```

# Инструкции
## Запуск приложения в контейнере и запуск тестов
* Для запуска приложения в докере можно воспользоваться командой `make build-restart`, которая пересоберёт контейнеры, запустит базу данных и контейнер с приложением, а также заставит отработать контейнер с миграциями и сидами.
```sh
make build-restart
```
* Запуск тестов `make test`, тест-контейнер с базой данных поднимется и опустится самостоятельно.
```
make test
```
* Просмотр покрытия `make coverage`, предварительно нужен хотя бы один прогон тестов, чтобы появился файл с покрытием.
```
make coverage
```
* Запуск линтера
```
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
make lint
```
## Запуск приложения локально (для разработки)
* Установить модули
```
go mod download
```
* Приложение поддерживает `.env` файл с переменными окружения. См. пример `.env.example`. Можно переименовать его «как есть» и добавить туда нужные параметры базы данных.
```
cp .env .env.example
```
* Для миграций используется goose, инструкции по запуску миграций — ниже по файлу.
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
* Запуск сервера
```
go run main.go
```

# Нагрузочное тестирование

Должны быть запущены контейнеры с приложением, базой и редисом, установлены jq и oha.

```
make load-test
```

# Замечания по реализации
* И транзакции и у заказа в базе данных есть своё поле `user_id`. В текущей логике это может привести к рассогласованию данных, когда вещь передана одному человеку, а деньги списаны со счёта другого. Но это даёт гибкость для будущих реализаций, например если нужно будет реализовать дарение вещи, а не просто передачу монеток.
* Для сущности `Info` (отчёт о монетах, инвентаре и транзакциях) прописаны теги для json и db. Вообще, им наверное не место в предметной области, но добавлять слои моделей и презентеров для такого небольшого проекта кажется оверинжинирингом. Будем считать, что для отчёта это окей.
* Эндпойнты со слешом на конце (`/auth/` и `/auth` например) реализованы как редиректы с первого на второй. При проверке нужно указать это клиенту (например `curl -L -X POST...`).

# Версии и требования
* Go 1.24
* Postgres 17
* Redis 7
* Для E2E тестов используются тест-контейнеры с Постгресом и Редисом, поэтому для запуска тестов нужно, чтобы на машине был Докер.



# Миграции при помощи goose
Проверить, текущий статус
```
goose postgres "postgres://prepin:@localhost:5432/av-merch-shop?sslmode=disable" -dir=schema/migrations status
```
Запустить
```
goose postgres "postgres://prepin:@localhost:5432/av-merch-shop?sslmode=disable" -dir=schema/migrations up
```

Сиды
```
goose postgres "postgres://prepin:@localhost:5432/av-merch-shop?sslmode=disable" -dir=schema/seed/ -no-versioning up
```
```
goose postgres "postgres://prepin:@localhost:5432/av-merch-shop?sslmode=disable" -dir=schema/seed -no-versioning reset
```
