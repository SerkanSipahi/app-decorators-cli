# app-decorators-cli (Beta)

> Start building a app-decorators webcomponent

### Installation

```sh
# get app-decorators via github
git clone https://github.com/SerkanSipahi/app-decorators-cli.git

# make sure you have "go" installed!
# If not please visit following page https://golang.org/dl/

cd app-decorators-cli
go build *.go
mv appdec /usr/local/bin/appdec
```

### Quickstart
```sh
appdec create my-module
appdec run --name=my-module --watch --server

# then open localhost:3000
```

### Commands

`appdec create your-app-name`: create a new app

`appdec run`: start or build an app

`appdec delete`: start a dev server

`appdec help`: help

### CLI Options

```sh
$ appdec create

  --name           Directory and package name for the webcomponent.
  
$ appdec delete

  --name           Delete directory or package for the webcomponent.

$ appdec run

  --help           show all options (see below)
  
  --name           set name of component
  --watch          watch file                                           [default: false]
  --server         start server                                         [default: false]
  --production     set production environment                           [default: false]
  --format         define component format (amd|cjs|umd|esm)            [default: "default"]
  --minify         will minify code                                     [default: false]
  --debug          will show debug messages                             [default: false]
  --no-mangle      no mangle                                            [default: false]
  --test           (under construction)                                 [default: false]
```

### Advanced usage

This is just a recommendation, you can use the options in any combination:
```sh

======================
### create command ###
======================

# create a new app
appdec create --name=my-module

===================
### run command ###
===================

# compile all file in project src directory
appdec run --name=my-module

# --watch: compile file on any file change
appdec run --name=my-module --watch

# --server: start server on port 3000
appdec run --name=my-module --server --watch

# --production: build a bundle file
appdec run --name=my-module --production --watch --server

# --format: set module format (work only with --production)
appdec run --name=my-module --production --format=cjs --watch --server

# --minify: miniy,reduce the file
appdec run --name=my-module --minify --production --format=cjs --watch --server

# --no-mangle: mangle
appdec run --name=my-module --no-mangle=true --watch --production --format=cjs --server --minify

# --debug: mangle
appdec run --name=my-module --debug --watch --production --format=cjs --server --minify --no-mangle=true

======================
### delete command ###
======================

# delete existing app
appdec delete --name=my-module
```

#### Babel

To customize Babel, you have two options:

1. You override [`.babelrc`](https://babeljs.io/docs/usage/babelrc/) file in your project's root directory.

#### Systemjs

To customize Systemjs, override ```jspm.config.js``` file which set settings that will change Systemjs config.