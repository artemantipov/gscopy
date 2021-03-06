## GSCOPY
Simple _COPY_ functional of [gsutils](https://cloud.google.com/storage/docs/gsutil) written in Go

### Build
``` go mod download && go build -o gscopy .```

### Usage
___Prerequisites___: You have to set variable `GOOGLE_APPLICATION_CREDENTIALS` for service account according [Google Documentation](https://cloud.google.com/docs/authentication/getting-started#setting_the_environment_variable)

#### Help
```./gscopy -h```

Output:
```bash
Usage: gscopy [-m] gs://bucketname/remote/dir /local/dir
  -m int
        number of concurrent copy tasks (default 1)
```


#### Copy (recursive)
Single-thread:

```./gscopy gs://bucket-name/mydir /local/path```

#### Multi-thread (concurrently):

```./gscopy -m 10 gs://bucket-name/mydir /local/path```

(e.g. `-m 10` or `-m=10` for maximum 10 files copied simultaneously, default 1)
