# tap-cli

## Requirements

### Binary
There is no requirements for binary app.

### Compilation
* git (for pulling repo only) 
* go >= 1.6

## Compilation
* clone this repo
* in dir of just cloned repo invoke: `make build_anywhere(_linux/_osx/_win32)`
* binaries are available in ./application/

`make build_anywhere` will compile binaries for all platforms, `make build_anywhere_linux` for Linux etc.


## Usage
```
$ ./application/tap
NAME:
   TAP CLI - client for managing TAP

USAGE:
   tap [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     login                    login to TAP
     target                   print actual credentials
     catalog                  list available services
     create-offering, co      create new offering
     create-service, cs       create instance of service
     delete-service, ds       delete instance of service
     bind-instance, bind      bind instance to another
     unbind-instance, unbind  unbind instance from another
     push                     create application from archive, manifest should be in current working directory
     applications, apps       list applications
     application, a           application instance details
     services, svcs           list services
     service, s               service instance details
     scale, sc                scale application
     start                    start application with single instance
     stop                     stop all application instances
     logs, log                get logs for all containers in instance
     delete, d                delete application
     help, h                  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Examples

### Authentication flow
```
$ kubectl get svc
NAME                  CLUSTER-IP   EXTERNAL-IP   PORT(S)                              AGE
api-service       10.0.0.243   <nodes>       80/TCP                               1h
container-broker      10.0.0.87    <nodes>       80/TCP                               1h
kubernetes            10.0.0.1     <none>        443/TCP                              1h

$ ./application/tap login 10.0.0.243 admin password
Authenticating...
OK

$ ./application/tap target
+------------+----------+
|    API     | USERNAME |
+------------+----------+
| 10.0.0.243 | admin    |
+------------+----------+

```
