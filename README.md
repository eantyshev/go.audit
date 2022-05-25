[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

# go.audit

Trivial audit log service, providing API to store and retrieve audit events

_id_ -- uuid4 strings

## Manual use

После `make run` REST API сервиса доступно на порту 8888
(можно также просматривать содержимое Редиса на порту 7379)
* Добавить новый баннер 1 и привязать его к слоту 4:
> curl -v -H "Content-Type: application/json" http://localhost:8888/add_banner -d '{"slot_id": 4, "banner_id": 1, "description": "descr 22"}'
* Привязать баннер 1 к слоту 5
> curl -v -H "Content-Type: application/json" http://localhost:8888/add_banner -d '{"slot_id": 5, "banner_id": 1}'
* Отвязать баннер от слота 4
> curl -v -H "Content-Type: application/json" http://localhost:8888/remove_banner -d '{"slot_id": 4, "banner_id": 1}'
* Выбрать баннер для показа группе 100 в слоте 5
> curl -v -H "Content-Type: application/json" http://localhost:8888/choose_banner -d '{"slot_id": 5, "group_id": 100}'
> 
> {"id": 1}
* Добавить клик баннеру 1
> curl -v -H "Content-Type: application/json" http://localhost:8888/click -d '{"slot_id": 5, "banner_id": 1, "group_id": 100}

## Unit tests

```make unittest```

Алгоритм UCB1 тестируется в юнит тестах

## Integration tests

Чтобы запустить интеграционные тесты, надо сделать `make test`
Запустится тестовое окружение и контейнер с тестами, код возврата команды соответствует успеху или неуспеху