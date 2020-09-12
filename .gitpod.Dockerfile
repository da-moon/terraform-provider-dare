FROM gitpod/workspace-full
USER gitpod
RUN go env -w GOPRIVATE=github.com/da-moon