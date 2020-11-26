## GSCOPY
Simple _COPY_ functional of [gsutils](https://cloud.google.com/storage/docs/gsutil) written in Go

### Build
``` go build -o gscopy .```

### Usage
Prerequisites: You have to set __GOOGLE_APPLICATION_CREDENTIALS__ for service account variable according [Google Documentation](https://cloud.google.com/docs/authentication/getting-started#setting_the_environment_variable)

#### Copy (recursive)
Single-thread:

```./gscopy gs://bucket-name/mydir /local/path```

Multi-thread(concurrently):

```./gscopy -m 10 gs://bucket-name/mydir /local/path```

(e.g. __-m 10__ or __-m=10__ - maximum 10 files copied simultaneously, default 1)
