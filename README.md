# tap-cli

## Requirements

### Binary
There are no requirements for binary app.

### Compilation
* git (for pulling repo only) 
* go >= 1.6

## Compilation
* clone this repo
* in dir of repo just cloned, invoke: `make build_anywhere(_linux/_osx/_win32)` if you don't want to bother with setting golang workspace.
* alternatively when you have golang workspace and GOPATH set use `make build` *(will be much faster on 2nd and next builds)*
* binaries are available in ./application/

`make build_anywhere` will compile binaries for all platforms, `make build_anywhere_linux` for Linux, etc.


## Usage
```
./tap
NAME:
   TAP CLI - client for managing TAP

USAGE:
   tap-cli [global options] command [command options] [arguments...]

VERSION:
   0.8.0

COMMANDS:
     login                    login to TAP. If you don't provide password you'll be promped for it.
     info                     prints info about current api and user
     offering                 offering context commands
     create-service, cs       create instance of service
     delete-service, ds       delete instance of service
     service-start            start service
     service-stop             stop all service instances
     service-restart          restart service
     expose-service, expose   expose service ports
     bindings                 list bindings
     bind-instance, bind      bind instance to another
     unbind-instance, unbind  unbind instance from another
     push                     create application from archive provided or from compressed current directory by default,
                              manifest should be in current working directory
     applications, apps       list applications
     application, a           application instance details
     services, svcs           list all service instances
     service, s               service instance details
     scale, sc                scale application
     start                    start application with single instance
     stop                     stop all application instances
     restart                  restart application
     logs, log                get logs for all containers in instance
     credentials, creds       get credentials for all containers in service instance
     delete, d                delete application
     invite                   invite new user to TAP
     reinvite                 resend invitation for user
     users                    list platform users
     invitations, invs        list pending invitations
     delete-invitation, di    delete invitation
     delete-user, du          delete user from TAP
     chpasswd                 change password of currently logged user
     help, h                  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbosity value, -v value  logger verbosity [CRITICAL,ERROR,WARNING,NOTICE,INFO,DEBUG] (default: "CRITICAL")
   --help, -h                   show help
   --version, -V                print the version
```

## Offerings

### Context info
```
./tap offering
NAME:
   TAP CLI offering - offering context commands

USAGE:
   TAP CLI offering command [command options] [arguments...]

COMMANDS:
     info    show information about specific offering
     list    list available offerings
     create  create new offering
     delete  delete offering

OPTIONS:
   --verbosity value, -v value  logger verbosity [CRITICAL,ERROR,WARNING,NOTICE,INFO,DEBUG] (default: "CRITICAL")
   --help, -h                   show help

```

## Examples

### Authentication flow
```
./tap login api.exampledomain.com admin password
Authenticating...
Authentication succeeded

#If you omit address you will be logged to previously set target

./tap login admin password

./tap target
+-------------------------+----------+
|           API           | USERNAME |
+-------------------------+----------+
| api.exampledomain.com   | admin    |
+-------------------------+----------+
```

### Application preparation *Python*

#### Prepare list of dependencies used in requirements.txt
Can be done manually, or when using virtualenv, dumped using:
```
pip freeze > requirements.txt
```
#### Store python dependencies in a folder:
```
mkdir vendor
sudo pip install -r requirements.txt --download vendor
```
#### Prepare run.sh script which will install dependencies and start an application:

```
#!/usr/bin/env bash

pip install --no-index --find-links=./vendor -r requirements.txt
python ./src/__init__.py
```

#### Create an archive containing all files loosely:
```
tar czvf python-application.tar.gz ./*
```
#### Create manifest.json file in current directory describing created application:

```
{
    "type":"PYTHON2.7",
    "name":"my-python-app",
    "instances":1
}

```
#### Push application 

```
./tap push python-application.tar.gz
```


### Application preparation *Java*

Build jar and prepare all dependencies 

#### Prepare run.sh script which will start an application:

```
#!/usr/bin/env bash

exec java -jar java-application-0.1.0.jar 
```

#### Create an archive containing all files loosely:
```
tar czvf java-application.tar.gz ./*
```


#### Write manifest.json file describing created application:
```
{
    "type":"JAVA",
    "name":"my-java-app",
    "instances":1
}

```


#### Push application 

```
./tap push java-application.tar.gz
```
