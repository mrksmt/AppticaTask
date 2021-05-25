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

<img src="https://user-images.githubusercontent.com/41635300/119490271-52986e00-bd65-11eb-92c2-de2359537924.png" width="50%"></img> 

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







