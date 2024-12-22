# FunCalc

## Шаг 1: Запуск Сервиса

Запуск:

```sh
go get .
go run main.go 
```

Это запустит веб-сервис на порту 8080 по умолчанию.

## Шаг 2: Отправка Запроса с помощью `curl`

Запрос должен быть отправлен на URL `/api/v1/calculate` с JSON-данными в теле запроса.

### 1 Успешное вычисление

Для вычисления выражения "2 + 2 * 3", используйте следующую команду:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 + 2 * 3"}' http://localhost:8080/api/v1/calculate
```

Сервер вернет JSON-ответ с результатом:

```json
{"result":8}
```

### 2 Некорректное выражение

Если вы отправите некорректное выражение, например, содержащее буквы:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"expression": "2 + a"}' http://localhost:8080/api/v1/calculate
```

Сервер вернет ошибку с кодом 422 (Unprocessable Entity):

```json
{"error":"Expression is not valid"}
```

### 3 Деление на ноль

При попытке деления на ноль:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"expression": "1 / 0"}' http://localhost:8080/api/v1/calculate
```

Сервер вернет ошибку с кодом 500 (Internal Server Error):

```json
{"error":"division by zero"}
```

Основные моменты:
`-X POST`: Указывает на использование метода POST.
`-H "Content-Type: application/json"`: Указывает тип содержимого запроса как JSON.
`-d '{"expression": "..."}'`: Передает JSON-данные в теле запроса. Замените "..." на ваше арифметическое выражение.
`http://localhost:8080/api/v1/calculate`: URL-адрес вашего веб-сервиса.

По вопросам тг: Artem1Belov
