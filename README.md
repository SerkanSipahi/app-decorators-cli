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
  --minify         will minify code                                     [default: false]
  --debug          will show debug messages                             [default: false]
  --format         define component format (amd|cjs|umd|esm)            [default: "default"]
  --no-mangle      no mangle                                            [default: false]
```

#### Babel

To customize Babel, you have two options:

1. You override [`.babelrc`](https://babeljs.io/docs/usage/babelrc/) file in your project's root directory.

#### Systemjs

To customize Systemjs, override ```jspm.config.js``` file which set settings that will change Systemjs config.