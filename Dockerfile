ARG ARCH=

FROM ${ARCH}\/golang:alpine as build

WORKDIR /src

RUN sed -i -e 's/https/http/' /etc/apk/repositories \
  && apk update

ENV USER=speedtest
ENV UID=1414

RUN adduser \    
    --disabled-password --gecos "" \    
    --home "/" --shell "/sbin/nologin" \    
    --no-create-home --uid "${UID}" "${USER}"

COPY ./go.* ./
RUN go mod download

COPY ./ ./

RUN export tmp=${ARCH} \
    export ARCH=$(uname -m) && \
    if [ $ARCH == 'armv7l' ]; then export ARCH=arm; fi && \
    wget -O /src/scripts/speedtest.tgz https://bintray.com/ookla/download/download_file?file_path=ookla-speedtest-1.0.0-${ARCH}-linux.tgz && \
    tar zxvf /src/scripts/speedtest.tgz -C /src/scripts && \
    rm -f /src/scripts/speedtest.md /src/scripts/speedtest.5 /src/scripts/speedtest.tgz \
    export ARCH=${tmp}
    
RUN mkdir -p /src/assets

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/server ./cmd/speedtest-api/

FROM ${ARCH}\/alpine AS deploy

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /src/scripts/ /scripts/
COPY --from=build /src/assets /assets/
RUN chown -R speedtest:speedtest /assets

USER speedtest:speedtest

COPY --from=build /out/server ./server
EXPOSE 9001

ENTRYPOINT  [ "/server" ]