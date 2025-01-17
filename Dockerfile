FROM golang:1.23.5 AS build

WORKDIR /src
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src .

RUN GOARCH=amd64 GOOS=linux GOEXPERIMENT=aliastypeparams go build -tags lambda.norpc -o main .
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /src/main ./main
ENTRYPOINT [ "./main" ]
