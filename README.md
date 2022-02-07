# Playing with consul

## Инструкция

Запуск сервера

```bash
docker run \
    -d \
    -p 8500:8500 \
    -p 8600:8600/udp \
    --name=badger \
    consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```

Получить ноды "датацентра"

```bash
docker exec badger consul members
```

Запуск управляющего клиента

```bash
docker run \
   --name=fox \
   consul agent -node=client-1 -join=172.17.0.2
```

Пример конфигурации сервиса в клиентской ноде. Конфиг лежит в /consul/config/counting.json

```json
{
    "service": {
        "name": "counting",
        "tags": ["go"],
        "port": 9001
     }
}
```

Применяется так:

```bash
docker exec fox consul reload
```

После этого можно запустить тестовый сервис:

```bash
docker run \
   -p 9001:9001 \
   -d \
   --name=weasel \
   hashicorp/counting-service:0.0.2
```
