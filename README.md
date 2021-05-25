# Application Top Category Positions 

## Здание

Реализовать микросервис Application Top Category Positions – получение данных о позициях приложения в топе по категориям за определенный день.

### Базовый функционал: 

Сбор и сохранение prepared данных для endpoint-а.
Endpoint для получения позиций в топ чарте маркета по категориям.

### Дополнительный функционал: (по желанию)

Добавить ограничение по количеству запросов на endpoint с одного ip-адреса: 5 запросов в минуту.
Добавить логирование запросов на endpoint. 
Использовать grpc + protobuf для запросов.
Добавить тестирование.

### Технические требования:

Ограничений по использованию библиотек – нет.
Формат ответа endpoint – json/protobuf.
Для хранения prepared данных можете использовать любую DBMS.
Результат - ссылка на репозиторий с кодом.

...

## Архитектура приложения

В приложении четыре сервиса:Data Processor, HTTP Endpoint, GRPC Endpoint и Data Service.

<img src="https://user-images.githubusercontent.com/41635300/119490271-52986e00-bd65-11eb-92c2-de2359537924.png" width="75%"></img> 

### - Data Processor
Запрашивает данные со стороннего API, обрабатывает, полностью подготавливает к использованию и загружает в хранилище. Данне хранятся в виде готового JSON.

###  - HTTP Endpoint
Обеспечивает функционал доступа к данным с помощью http запроса вида
```
http://<ваш домен>/appTopCategory?date=2021-04-01
```
В ответе на запрос передается JSON с данными на запрашиваемую дату.
```json
{
   "status_code": 200,
   "message": "OK",
   "data": {
      "134": 73,
      "2": 26,
      "23": 5
   }
}
```
В случае если запрашиваемые данные не удалось получить -  передаются соответствующие ошибки.

### - GRPC Endpoint
Работает так же как HTTP Endpoint.

Тип запроса - `proto.GetPositionsRequest`. В запросе указывается дата в формате `2021-05-12`.

Ответ содержит структуру с запрашиваемой датой, статусом, сообщением и данными аналогично ответу http. Тип ответа - `proto.GetPositionsResponse`. 

Пример ответа:
```
Date:    2021-05-12
Status:  200
Message: OK
{
   "status_code": 200,
   "message": "OK",
   "data": {
      "134": 73,
      "2": 26,
      "23": 5
   }
}
```
Сервис поддерживает как одиночные (unary) таки потоковые (streaming) запросы.

### - Data Service
Обеспечивает хранение и многопользовательский доступ к данным.

Данные хранятся с использованием встроенной СУБД LedisDB. Обмен данными с клиентами происходит  с помощью протокола GRPC.

В рамках приложения сервис не несет функциональной нагрузки, и может быть заменен на любую другую реализацию системы хранения и доступа к данным с помощью замены соответствующего модуля доступа к данным на стороне клиентов.

## Использование

---
**NOTE!**

Понятно, что все это прототип.

---


Скачайте репозиторий локально
```bash
$ git clone https://github.com/mrksmt/AppticaTask.git
```




Запустите сервис доступ к данным
```bash
$ make dataservice
```


Запустите обработчик данных
```bash
$ make dataprocessor
```


Запустите HTTP Endpoint
```bash
$ make httpendpoint
```
```bash
$ curl http://192.168.31.25:8081/appTopCategory?date=2021-05-17
{"status_code":200,"message":"OK","data":{"134":59,"2":23,"23":5}}
```

Запустите GRPC Endpoint
```bash
$ make grpcendpoint
```
Попробуйте запустить тестового grpc клиента в режимах одиночного и потокового запроса и получить ответ от эндпоинта:
```bash
$ make client_unary 
cd cmd/grpcclient && GRPC_HOST=localhost:8082 REQUEST_TYPE=unary DATES="2021-05-12" go run main.go
Date:    2021-05-12
Status:  200
Message: OK
{
   "status_code": 200,
   "message": "OK",
   "data": {
      "134": 73,
      "2": 26,
      "23": 5
   }
}
```
``` bash
$ make client_streaming 
cd cmd/grpcclient && GRPC_HOST=localhost:8082 REQUEST_TYPE=streaming DATES="2006-01-02 2021-04-12 2021-05-12 2021-05-22 55555" go run main.go
Date:    2006-01-02
Status:  404
Message: Data not found

Date:    2021-04-12
Status:  404
Message: Data not found

Date:    2021-05-12
Status:  200
Message: OK
{
   "status_code": 200,
   "message": "OK",
   "data": {
      "134": 73,
      "2": 26,
      "23": 5
   }
}

Date:    2021-05-22
Status:  200
Message: OK
{
   "status_code": 200,
   "message": "OK",
   "data": {
      "134": 91,
      "2": 34,
      "23": 6
   }
}

Date:    55555
Status:  400
Message: Bad request err: parsing time "55555" as "2006-01-02": cannot parse "5" as "-"
```