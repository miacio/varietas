# varietas
Experimental business logic variants

## fls
File operation module of operating system paradigm

## dynamic
Registering an object based on the registration method allows the object to call it through a string

see dynamic.dynamic_test.go

## web
Based on the quick routing interface registration encapsulated by gin, this method can effectively distinguish between routing interfaces and registered routes

It ensures the practicality of the original gin and only adds a routing group operation to it

see web.web_test.go -> TestWeb001

## email
add email tool

because of [email](https://github.com/jordan-wright/email) package

see email.email_test.go

### web unique capabilities
High availability processing for large file uploads

Ability to slice and upload large files

see web.web_test.go -> TestChunkFileUploadServer and TestChunkFileUploadClient

## dbtp
Based database param type to golang struct package
