# Application Top Category Positions 

## Архитектура приложения

В приложении четыре сервиса:Data Processor, HTTP Endpoint, GRPC Endpoint и Data Service.

### Data Processor
Запрашивает данные со стороннего API, обрабатывает, полностью подготавливает к использованию и загружает в хранилище. Данне хранятся в виде готового JSON.

###  HTTP Endpoint
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

###  GRPC Endpoint
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

###  Data Service
Обеспечивает хранение и многопользовательский доступ к данным.

Данные хранятся с использованием встроенной СУБД LedisDB. Обмен данными с клиентами происходит  с помощью протокола GRPC.

В рамках приложения сервис не несет функциональной нагрузки, и может быть заменен на любую другую реализацию системы хранения и доступа к данным с помощью замены соответствующего модуля доступа к данным на стороне клиентов.







