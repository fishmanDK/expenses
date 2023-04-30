FROM postgres:alpine

# Установка необходимых пакетов
RUN apk add --no-cache bash
RUN apk update && apk add curl
# RUN apt-get update && apt-get install -y curl
# Копирование скрипта в контейнер
COPY script.sh /usr/local/bin/

# Установка прав на скрипт
RUN chmod +x /usr/local/bin/script.sh

# Запуск скрипта при старте контейнера
CMD ["script.sh"]