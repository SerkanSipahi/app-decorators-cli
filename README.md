# app-decorators-cli (Beta)

> Start building a [`app-decorators`](https://github.com/SerkanSipahi/app-decorators) webcomponent

### Installation (use existing binaries)
```sh
# get app-decorators via github
git clone https://github.com/SerkanSipahi/app-decorators-cli.git

cd app-decorators-cli

# for osx
mv ./bin/osx/appdec /usr/local/bin/appdec

# for linux
mv ./bin/linux/appdec /usr/local/bin/appdec

# windows bin is in ./bin/win/appdec.exe located
```

### Installation (self building)

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
appdec create --name=mymodule
appdec run --name=mymodule --watch --server

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
appdec create --name=mymodule

===================
### run command ###
===================

# compile all file in project src directory
appdec run --name=mymodule

# --watch: compile file on any file change
appdec run --name=mymodule --watch

# --server: start server on port 3000
appdec run --name=mymodule --server

# --production: build a bundle file
appdec run --name=mymodule --production

# --format: set module format (work only with --production)
# available formats: default|amd|cjs|umd|esm
appdec run --name=mymodule --production --format=cjs

# --minify: miniy,reduce the file
appdec run --name=mymodule --minify --production --format=cjs

# --no-mangle: mangle
appdec run --name=mymodule --no-mangle=true --production --format=cjs --minify

# --debug: mangle
appdec run --name=mymodule --debug

======================
### delete command ###
======================

# delete existing app
appdec delete --name=mymodule
```

#### Babel

To customize Babel:

You override [`.babelrc`](https://babeljs.io/docs/usage/babelrc/) file in your project's root directory.

#### Systemjs

To customize Systemjs, override ```jspm.config.js``` file which set settings that will change Systemjs config.
