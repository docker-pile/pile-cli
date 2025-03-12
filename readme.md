## development
install dependencies
```
go mod init pile-cli
go get github.com/spf13/cobra

./build.sh
```


## init directory
$HOME/pile/pile.env
$HOME/pile/.tmux.conf
$HOME/pile/beaver/beaver.yaml
$HOME/pile/jwtio/jwtio.yaml

## Cli Interface
pile init
pile edit <service-name>
pile env
pile logs <service-name>
pile up
pile up --[g]roups data,app,utils
pile down --[g]roups app
pile down
pile install <pile-library-name>
pile list
pile help

pile status/ps
pile commands
pile ports
pile images


## pile.env
install: postgres mysql redis open-webui beaver

groups:
  data:
    postgres
    mysql
    redis
  apps:
    open-webui
  work:
    oneimaging
  utils:
    beaver

install:
  mongo

    