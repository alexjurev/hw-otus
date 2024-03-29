FROM gomicro/goose

ADD migrations/*.sql /migrations/
ADD build/migrations-entrypoint.sh /migrations/entrypoint.sh
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /migrations/wait
RUN chmod +x /migrations/wait
RUN chmod +x /migrations/entrypoint.sh

ENTRYPOINT ["/migrations/entrypoint.sh"]