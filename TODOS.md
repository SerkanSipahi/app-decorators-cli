#### todos
- appdec.go
-- check weather an new version app-decorator is available. When "yes" ask for upgrade
-- build: go build *.go

- util/watch/watch.go
-- move watch.go to own repo: http://github.com/serkansipahi/watcher
-- check weather on runtime new directory with file will created
-- allow to pass something like this ./collapsible -r --ignore=node_modules --type=modify
-- example:
--- cmdStr := "./collapsible -r -ignore=node_modules"
--- cmdSplited := Split(cmdStr, " ")
--- cmd := exec.Command(cmdSplited...)
- when "a" changed, then b, sometimes intellij trigger "a" and "b" instead of just "b"
- http://stackoverflow.com/questions/10383498/how-does-go-update-third-party-packages
