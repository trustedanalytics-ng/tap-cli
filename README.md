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
./tap
NAME:
   TAP CLI - client for managing TAP

USAGE:
   tap-cli [global options] command [command options] [arguments...]

VERSION:
   0.8.0

COMMANDS:
     login                    login to TAP
     target, t                print actual credentials
     catalog, o               list available offerings
     create-offering, co      create new offering
     delete-offering, do      delete offering
     create-service, cs       create instance of service
     delete-service, ds       delete instance of service
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

## Examples

### Authentication flow
```
./tap login api.exampledomain.com admin password
Authenticating...
Authentication succeeded

./tap target
+-------------------------+----------+
|           API           | USERNAME |
+-------------------------+----------+
| api.exampledomain.com   | admin    |
+-------------------------+----------+
```

### Application preparation *Python*

#### Prepare list of dependencies used in requirements.txt
It can be done manually or when using virtualenv dumped using:
```
pip freeze > requirements.txt
```
#### Store python dependencies in a foler:
```
mkdir vendor
sudo pip install -r requirements.txt --download vendor
```
#### Prepare run.sh script which will install dependencies and start an application:

```
pip install --no-index --find-links=./vendor -r requirements.txt
python ./src/__init__.py
```

#### Create an archive containing all files loosely:
```
tar cxvf python-application.tar.gz ./*
```
#### Create manifest.json file in current directory describing created application:

```
{
    "type":"PYTHON",
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
exec java -jar java-application-0.1.0.jar 
```

#### Create an archive containing all files loosely:
```
tar cxvf java-application.tar.gz ./*
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
