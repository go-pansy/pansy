### Generate a hex fake ID based on a real ID

### InstallInstallation
```sh
go get -u github.com/go-pansy/pansy
```

### Usage
```
    // gen
    faker      := pansy.NewFaker()
    fakerId := faker.GenId(id, true)

    // recover
    id = faker.RecoverId(fakerId)
```
