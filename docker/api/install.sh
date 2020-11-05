if [ ${USE} = "DEBUG" ]; then 
    CompileDaemon -directory=/go/src/src/ -log-prefix=false -command="dlv debug --headless --listen=:2345 --api-version=2 --log"
else
    CompileDaemon -directory=/go/src/src -log-prefix=false -build="go build -a -installsuffix cgo -o /go/bin/hello" -command="/go/bin/hello"
fi