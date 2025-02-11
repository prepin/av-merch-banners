# Наблюдения
* В схеме указано, покупка предмета за монеты происходит через `GET /api/buy/{item}`. Это вызывает вопросы, потому, что GET запрос не должен изменять данные на сервере, а только возвращать их. Решение: Делаем реализацию через POST запрос. Для GET делаем фоллбэк на случай, если уже существует клиент, который почему-то ориентируется на GET, но помечаем его как deprecated.
* API указано с префиксом просто `/api`, без версионирования. Это повлечёт за собой некоторую боль в будущем, когда понадобится реализовать изменения, ломающие совместимость. Решение: монтируем API по префиксу `/api/` и префиксу `/api/v1/`

# Миграции при помощи goose
Проверить
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

goose postgres "postgres://prepin:@localhost:5432/av-merch-shop?sslmode=disable" -dir=schema/seed -no-versioning reset
```


# Полезные файлы
* [Человекочитаемая OpenAPI cхема](docs/redoc.static.html)
