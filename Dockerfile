FROM scratch

# provide invalid default so that if APP_NAME
# isn't provided, the build will fail
ARG APP_NAME=INVALID_APP_NAME

COPY ./bin/$APP_NAME /app

ENTRYPOINT [ "/app" ]