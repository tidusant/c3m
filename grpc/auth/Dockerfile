FROM alpine
# Add Maintainer Info
LABEL maintainer="Duy Ha <duyhph@gmail.com>"
# Set the Current Working Directory inside the container
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

WORKDIR /app
# Copy exec file and config
COPY main ./

# Build the Go app
#RUN go build -o main .

# Expose 8901 port to the outside world
EXPOSE 8901
# Run the executable
CMD ["./main"]