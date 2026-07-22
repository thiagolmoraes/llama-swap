FROM docker:27-cli

RUN apk add --no-cache curl docker-cli-compose

COPY build/llama-swap-linux-amd64 /usr/local/bin/llama-swap
RUN chmod +x /usr/local/bin/llama-swap

ENTRYPOINT ["/usr/local/bin/llama-swap", "-config", "/config/config.yaml", "-listen", ":8080"]
