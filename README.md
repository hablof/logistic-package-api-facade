Репозиторий является частью проекта https://github.com/hablof/logistic-package-api

## logistic-package-api-facade

Сервис выполняет простейшую функцию: читает сообчения в кафка-топиках   
  - "omp-package-events"
  - "omp-tgbot-commands"
  - "omp-tgbot-cache-events"

и выводит сообщения в `stdout` в читаемом виде.

### Топик "omp-package-events"

Сообщения в топике "omp-package-events" сериализованы в protobuf, описанном в субмодуле `pkg/kafka-proto` основного [репозитория](https://github.com/hablof/logistic-package-api). Сообщения десериализуются в объект и сериализуются с помощью `json.MarshalIndent` и выводятся в `stdout`.

### Топики "omp-tgbot-commands" и "omp-tgbot-cache-events"

Сообщения в топиках "omp-tgbot-commands" и "omp-tgbot-cache-events" сериализованы в json. Сообщения преобразуются в Indent-вид и выводятся в `stdout`.

## Docker

Описан докерфайл для образа фасада.

## Makefile

В мейкфайле описаны команды для локального запуска приложения, сборки докер-образа, и его запуска.