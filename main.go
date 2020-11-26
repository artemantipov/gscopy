package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	var bucketArg, localPathArg, bucketName, bucketPath string

	concurFlag := flag.Int("m", 1, "number of concurrent copy tasks")
	flag.Parse()
	afterFlagArgs := flag.Args()
	bucketArg = afterFlagArgs[0]
	localPathArg = afterFlagArgs[1]
	bucketSlice := strings.SplitN(strings.TrimPrefix(bucketArg, "gs://"), "/", 2)
	bucketName = bucketSlice[0]
	bucketPath = bucketSlice[1]

	fmt.Printf("Bucket: %v Path: %v\n", bucketName, bucketPath)

	listing, err := listFiles(bucketName, bucketPath)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	if isFlagSet("m") {
		log.Printf("Multi-thread copy: %v", *concurFlag)
	} else {
		log.Println("Single-thread copy")
	}

	copyFile := func( f string) {
		path := localPathArg + "/" + f
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Printf("Error creating dirs %v : %v", filepath.Dir(path), err)
		}
		var writer bytes.Buffer
		w := bufio.NewWriter(&writer)
		fileContent, err := downloadFile(w, bucketName, f)
		if err != nil {
			log.Printf("Error downloading file %v : %v", f, err)
		}
		err = ioutil.WriteFile(path, fileContent, 0644)
		if err != nil {
			log.Printf("Error saving file %v: %v", path, err)
		}
		log.Printf("Copied: %v",path)

	}

	// Create WaitGroup and Semaphore channel to control routine amount and completion
	var wg sync.WaitGroup
	sem := make(chan int, *concurFlag)
	for _,f := range listing {
		sem <- 1
		wg.Add(1)
		go func(f string) {
			copyFile(f)
			<-sem
			wg.Done()
		}(f)
	}
	wg.Wait()
}

// isFlagPassed check if flag exists in command line args
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// listFiles lists objects within specified bucket.
func listFiles(bucket, path string) (listing []string, err error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := &storage.Query{Prefix: path}
	it := client.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects: %v", bucket, err)
		}
		listing = append(listing, attrs.Name)
	}
	return listing,nil
}

// downloadFile downloads an object.
func downloadFile(w io.Writer, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	fmt.Fprintf(w, "Blob %v downloaded.\n", object)
	return data, nil
}
