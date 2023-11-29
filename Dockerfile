FROM golang:1.21 as base
WORKDIR /app

FROM base as development

FROM base as build
COPY go.mod go.sum ./
COPY main.go .
RUN go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023 as test
RUN curl -Lo /usr/bin/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie && chmod +x /usr/bin/aws-lambda-rie
COPY --from=build /app/main ./main
ENTRYPOINT [ "/usr/bin/aws-lambda-rie", "./main" ]

FROM public.ecr.aws/lambda/provided:al2023 as production
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]
