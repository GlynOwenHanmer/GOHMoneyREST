FROM scratch

ARG APP_NAME

COPY ./bin/$APP_NAME /app

ENTRYPOINT [ "/app" ]