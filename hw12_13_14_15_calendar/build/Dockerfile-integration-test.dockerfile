# Собираем в гошке
FROM golang:1.16.2 as build

ENV WAIT_VERSION 2.7.2
ENV BIN_FILE /opt/calendar/tests
ENV WAIT_FILE /opt/calendar/wait
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go test --tags integration_test -c \
        -o ${BIN_FILE} test/integration/*

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait ${WAIT_FILE}
RUN chmod +x ${WAIT_FILE}

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar integration test"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/calendar/tests"
ENV WAIT_FILE /opt/calendar/wait
COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY --from=build ${WAIT_FILE} ${WAIT_FILE}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config-docker.yaml ${CONFIG_FILE}

CMD ${WAIT_FILE} && ${BIN_FILE} -test.v