Запуск

```
make run_in_mem //для запуска in-memory решения
make run_in_db //для запуска с использованием postgresql
```

Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый

Запрос
```
curl --header "Content-Type: application/json" \
--request POST \
--data '{"long_link": "https://gitlab.com/"}' \
  http://localhost:8080/add/
```
Ответ: ```{"short_link":"RiWMsVFzEX"}```

В случае ошибки выдает код ошибки и сообщение об ошибки ```{"message":"no such short link"}```

Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL

Запрос
```
curl --header "Content-Type: application/json" \
--request GET \
  http://localhost:8080/link/{SHORT_LINK}
```
где вместо ```{SHORT_LINK}``` нужно подставить сокращенный URL

Ответ: ```{"long_link":"https://gitlab.com"}```

В случае ошибки выдает код ошибки и сообщение об ошибки ```{"message":"no such short link"}```
