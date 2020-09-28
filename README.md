# cadenceLearning
My studying of cadence workflow engine. 

- Follow instructions on https://cadenceworkflow.io/ 
- Spin up docker compose. 
- UI: http://localhost:8088/ 
- Create a new domain, I'll use "test-domain" in following examples.


## Running cadence command
It's useful to create an alias:
```
$ alias cadence="docker run --rm ubercadence/cli:master --address host.docker.internal:7933"
```

## Commands:
```
$ cadence --domain test-domain workflow start --tasklist cadence-learning-app --workflow_type CheckDrivingLicence --et 600 -i '"Jan"'

$ cadence --domain test-domain workflow signal -w 5a907a10-e51c-430a-9a64-3cfe1a48e6a5 --name age-confirmation -i '"confirmed"'

$ cadence --domain test-domain workflow signal -w bf3aaf6b-0d65-4b94-b1f7-84e1761d7e95 --name stop -i '"true"'

$ cadence --domain test-domain workflow start --tasklist cadence-learning-app --workflow_type FlipImage --et 60 -i '{"url": "https://media.npr.org/assets/img/2020/05/29/ibam_arttile_final_sq-b6691366c421f19f130bd0beef01c1597c57bdf3-s400-c85.png", "output_path": "/tmp/test01.png"}'
```

## Running go code
```
$ go mod tidy
$ go run .
```
