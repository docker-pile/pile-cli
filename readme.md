## Brew Instillation
```
brew cask add docker-pile/pile-cli && \
brew install pile
```

### development
install dependencies
```
go mod init pile-cli
go get github.com/spf13/cobra

./build.sh
```


### init directory
```
$HOME/pile/
$HOME/pile/pile.config.yaml
$HOME/pile/pile.network.yaml
```
`pile install postgres` will result in creating this `docker compose` file.
```
$HOME/pile/postgres/compose.yaml
```

### Cli Interface
some of these are a work in progress. i put them here for my references during development.
```
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
```

## pile.env
now:
```
APPS:
  - open-webui
  - postgres
  - wordpress

```
future:
```
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
```

install:
  mongo

    
