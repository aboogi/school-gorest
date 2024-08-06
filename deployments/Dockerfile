ARG repo=golang:1.22-bullseye
FROM $repo AS build
LABEL stage=builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY internal ./internal/
COPY sql ./sql/


RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o "schoolmat"


FROM alpine as prod
# DESCRIPTION 
LABEL stage=deployer

RUN set -ex && apk --update add libc6-compat

WORKDIR /app
COPY --from=build /app/schoolmat .

EXPOSE 9000
ENTRYPOINT [ "./schoolmat" ]
