# Simulate
## Command
```
./lucy simulate -c=2
```

## Hook
```
#!/bin/sh

gofmt -l -s -w .

golint lucy/simulate

go vet lucy/simulate
```