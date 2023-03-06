# cadenceLearning
My studying of cadence workflow engine. 

- Follow instructions on https://cadenceworkflow.io/
  https://cadenceworkflow.io/docs/get-started/installation/#install-docker
- Alternatively: docker compose and build cli using insttructions here:
  https://github.com/uber/cadence/blob/master/docker/README.md
- UI: http://localhost:8088/


## Running cadence command

It is handy to have a Bash script in directory which comes before system ones in your `$PATH`:
```
$ cat ~/bin/cadence-local
#!/bin/bash

docker run --network=host ubercadence/cli:master "$@"
```

# Registering a new domain

Register a new domain using cli. I'm using "test-domain" in following examples:
```
$ cadence-local --do test-domain domain register -rd 1
```

## Running go code
```

# This command may fail with weird output.
# I had to remove go.mod file and do `go mod init` before 
$ go mod tidy
$ go run .
```

## Launching workflows
Now when we have a domain registered and also have a running code in terminal, let's launch workflow: 
```
$ cadence-local --domain test-domain workflow start --tasklist cadence-learning-app --workflow_type FlipImage --et 60 -i '{"url": "https://media.npr.org/assets/img/2020/05/29/ibam_arttile_final_sq-b6691366c421f19f130bd0beef01c1597c57bdf3-s400-c85.png", "output_path": "/tmp/test01.png"}'
$ open /tmp/test01.png # you should see an image flipped upside down
```

Now let's run a bit more complicated workflow which requires signals to proceed: 
```
$ cadence-local --domain test-domain workflow start --tasklist cadence-learning-app --workflow_type CheckDrivingLicence --et 600 -i '"Jan"'
$ cadence-local --domain test-domain workflow signal -w 5a907a10-e51c-430a-9a64-3cfe1a48e6a5 --name age-confirmation -i '"confirmed"'
$ cadence-local --domain test-domain workflow signal -w bf3aaf6b-0d65-4b94-b1f7-84e1761d7e95 --name stop -i '"true"'
```

## Local version of Cadence client

Since we're about to work on cadence client, we need to replace github version with a local checkout.
Make sure you have following in your go.mod:
```
replace go.uber.org/cadence => /Users/dkrot/go/src/go.uber.org/cadence
```

Now you can play around and see your changes to client are immediately take effect.
