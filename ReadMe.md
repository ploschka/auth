# Сервер аутентификации

## Развёртывание в production окружении

1. Установить переменные окружения либо в командной оболочке, либо в `deployment/.env`

    ```sh
    POSTGRES_DB=auth_server # Имя базы данных
    POSTGRES_PASSWORD=secret # Пароль полльзвателя базы данных
    POSTGRES_USER=app # Имя пользователя базы данных
    SIGN_KEY=c2VjcmV0X3NpZ25fa2V5 # Ключ подписи токенов в кодировке base64
    ENCRYPTION_KEY=bOVcHwoCIhSF5EM9gC15PAOY1KAm3i6h9lELYnh1BO4= # Ключ шифрования длиной 32 байта в кодировке base64
    ```

2. Развернуть приложение в docker в директории `deployment`

    ```sh
    docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
    ```

## Развёртывание в окружении разработчика

1. Установить зависимости

    ```sh
    go mod download && go mod verify
    ```

2. Установить переменные окружения в командной оболочке, либо в `deployment/.env`

    ```sh
    POSTGRES_DB=auth_server # Имя базы данных
    POSTGRES_PASSWORD=secret # Пароль полльзвателя базы данных
    POSTGRES_USER=app # Имя пользователя базы данных
    ```

3. Установить переменные окружения в командной оболочке

    ```sh
    - DATABASE_LINK=host=127.0.0.1 user=app password=secret dbname=auth port=5432 sslmode=disable # URL подключения к базе данных для приложения
    SIGN_KEY=c2VjcmV0X3NpZ25fa2V5 # Ключ подписи токенов в кодировке base64
    ENCRYPTION_KEY=bOVcHwoCIhSF5EM9gC15PAOY1KAm3i6h9lELYnh1BO4= # Ключ шифрования длиной 32 байта в кодировке base64
    ```

4. Развернуть базу данных в docker в директории `deployment`

    ```sh
    docker compose up -d
    ```

5. Запустить локальный сервер

    ```sh
    go run cmd/auth/auth.go
    ```
