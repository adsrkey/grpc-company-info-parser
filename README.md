# GRPC COMPANY INFO PARSER

#### Description:
#### Service for parsing company information by taxpayer identification number

## How to run with docker-compose:

```
    make up_build
```
---

#### You can regenerate proto: 

```
    make gen
```

#### You can generate swagger documentation based on proto:

```
    make doc
```

### GET
``` 
    http://{ADDR}:8080/api/v1/company/inn/{INN}/info
``` 