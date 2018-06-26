FROM qasico/go:latest

# Apps Dir
RUN mkdir -p "$GOPATH/src/git.qasico.com/cuxs/cuxs"
ADD . "$GOPATH/src/git.qasico.com/cuxs/cuxs"
WORKDIR "$GOPATH/src/git.qasico.com/cuxs/cuxs"

# Installing Dependency
RUN glide update

# Entry Point
ENTRYPOINT "$GOPATH/src/git.qasico.com/cuxs/cuxs"

CMD ["sleep","3600"]
