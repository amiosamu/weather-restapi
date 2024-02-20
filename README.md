# weather-restapi

Тестовое задание на позицию Intern Go разработчика в компанию Kwaaka.


Используемые технологии

- mongo-driver/mongo (для хранилища данных)
- docker и docker-compose (для запуска сервиса)
- gin-gonic/gin (веб фреймворк)
- golang/mock, testify (для тестирования)

## Запуск

- Заполнить .env по .env_example
- Запустить сервис можно с помощью команды `sudo make compose-up`



## Тестирование

Для запуска тестов необходимо выполнить команду `make test`


## Примеры запросов

- [Получение погоды](#get-weather)
- [Обновление погоды](#update-weather)


### Получение информации о погоде по городу <a name="get-weather"></a>

Пример запроса

```
curl -X GET "http://localhost:8080/api/weather/Astana"

```

Пример ответа

```
{
    "id": "65d4894105e1600b7dee4535",
    "city": "Astana",
    "temperature": 253.12,
    "last_updated": "2024-02-20T11:13:05.300Z",
    "code": 200
}

```


### Обновление инфромации о погоде по городу <a name="update-weather"></a>

Пример запроса

```
curl -X PUT "http://localhost:8080/api/weather/" \
     -H "Content-Type: application/json" \
     -d '{
        "city": "Astana"
     }'

```

Примет ответа

```
{
    "id": "65d48ff155963693dc6d14f1",
    "temperature": 253.12,
    "code": 200
}
```
