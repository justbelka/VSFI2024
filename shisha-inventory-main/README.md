# Сборка docker образов

### Сборка клиента

```bash
cd shisha-inventory

docker build -t shisha-client -f Dockerfile-client .
```

### Сборка сервера

```bash
cd shisha-inventory

docker build -t shisha-server -f Dockerfile-server .
```

После сборки правильно тегаем образы исходя из нашего registry (nexus, gitlab registry, etc) и пушим его туда

# Запуск зависимых сервисов

docker compose up -d

# Запуск бэка

cd server
go run main.go --s3-endpoint localhost:9000

# Запуск фронта

cd client
npm start

# Tests

go test ./tests/

# Lint

golangci-lint run

# Deploy with kubernetes

Инструкция по установке shisha-inventory в k8s/k3s

## Deploy thirdparty

Установка приложений, которые нужны для нашего приложения

Скачиваем helm репозиторий

```
helm repo add bitnami https://charts.bitnami.com/bitnami
```

#### Redis

Redis используется для баз данныхб а также для реализации кэш-хранилищ или брокеров сообщений.

https://artifacthub.io/packages/helm/bitnami/redis

```
helm install redis bitnami/redis -f k8s/thirdparty/redis.yml
```

#### Postgresql

Postgresql (ака постгря) - классическая и популярная СУБД типа SQL (т.е. реляционная и структурированная):

https://artifacthub.io/packages/helm/bitnami/postgres

```
helm install postgres bitnami/postgresql -f k8s/thirdparty/postgresql.yml

```

#### Minio

S3-хранилище для хранения файлов или картинок/фото (является NoSQL СУБД):

https://artifacthub.io/packages/helm/bitnami/minio

```
helm install minio bitnami/minio --version 14.6.16 -f k8s/thirdparty/minio.yml
```

#### Cert Manager

Приложение для управления сертификатами (требуется для работы брокера сообщений RedPanda)

```
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --set installCRDs=true \
  --namespace cert-manager \
  --create-namespace
```

#### RedPanda

Брокер сообщений, используемый для обмена данными между приложениями. Аналог Kafka но на С++, а не Java

```
helm repo add redpanda https://charts.redpanda.com
kubectl create configmap redpanda-io-config --from-file=k8s/thirdparty/io-config.yml
helm install redpanda redpanda/redpanda --version 5.8.11 -f k8s/thirdparty/redpanda.yml
```

## Deploy Shisha

Развёртывание frontend и backend

### Deploy Backend

```
helm install server ./k8s/helm/server
```

### Deploy Client

```
helm install client ./k8s/helm/client
```

# Про CICD

Добавляем чарт:

```
helm repo add gitlab-runner https://nexus.vsfi.ru/repository/helm-rw/  --username helm-user --password helm-user
```

Создаем неймспейс и секрет gitlab-runner:

```
k create ns gitlab-runner

kubectl create secret docker-registry image-pull-secret \
  --namespace gitlab-runner \
  --docker-server="https://registry.vsfi.ru" \
  --docker-username="docker-user" \
  --docker-password="docker-user"
```

Создаем секрет для скачивания образов через docker in docker

```
kubectl create configmap docker-client-config --namespace gitlab-runner --from-file ./config.json
```

Делоим раннера (не забываем сменить runnerToken на свой)

```
helm upgrade -i --namespace gitlab-runner gitlab-runner -f runner-values.yaml gitlab-runner/gitlab-runner
```

Добавляем политики для раннера, чтобы можно было создавать поды в других неймспейсах:

```
kubectl create clusterrolebinding gitlab-runner-binding --clusterrole=cluster-admin --serviceaccount=gitlab-runner:default
```

Добавляем креды для логина в registry-rw.vsfi.ru:

```
Settings -> CI/CD -> Variables
Добавляем 2 переменные:
CI_REGISTRY_USER docker-user
CI_REGISTRY_PASSWORD docker-user

Снимаем галочку protected
```
