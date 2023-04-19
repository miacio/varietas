# varietas
Experimental business logic variants

## dynamic
Registering an object based on the registration method allows the object to call it through a string

see dynamic.dynamic_test.go

## web
Based on the quick routing interface registration encapsulated by gin, this method can effectively distinguish between routing interfaces and registered routes

It ensures the practicality of the original gin and only adds a routing group operation to it

see web.web_test.go -> TestWeb001

### web unique capabilities
High availability processing for large file uploads

Ability to slice and upload large files

see web.web_test.go -> TestChunkFileUploadServer and TestChunkFileUploadClient

## dbtp
Based database param type to golang struct package

## email
add email tool

because of [email](https://github.com/jordan-wright/email) package

see email.email_test.go

## log
add log package

because of uber.zap package

extremely simplified and fast build logs

``` go
import "github.com/miacio/varietas/log"

func main() {
    logParam := log.LoggerParam{
        Path:       "./log", // you log write folder path
        MaxSize:    256,
        MaxBackups: 10,
        MaxAge:     7,
        Compress:   false,
    }
    
    log := logParam.Default()

    log.Infoln("init success")
}
```

## util
util package encapsulated some commonly used basic methods

add slice stream logic methods

## mfs
mfs package a workflow pattern method factory developed based on facet oriented thinking and communicated through context management