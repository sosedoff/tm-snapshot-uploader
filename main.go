package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	chunkSuffix  = ".chunk"
	snapshotName = "snapshot"
)

var cliOpts = struct {
	region string
	bucket string
	prefix string
	dir    string
	watch  bool
	debug  bool
}{}

func init() {
	flag.StringVar(&cliOpts.region, "region", "us-east-1", "AWS region")
	flag.StringVar(&cliOpts.bucket, "bucket", "", "AWS S3 bucket name")
	flag.StringVar(&cliOpts.prefix, "prefix", "", "Prefix name for uploads in the bucket")
	flag.StringVar(&cliOpts.dir, "dir", "", "Directory containing height snapshots")
	flag.BoolVar(&cliOpts.watch, "watch", false, "Run in watch mode (periodically scan dir)")
	flag.BoolVar(&cliOpts.debug, "debug", false, "Enable AWS debugging")
	flag.Parse()

	if cliOpts.dir == "" {
		log.Fatal("snapshots dir is not provided")
	}

	if cliOpts.bucket == "" {
		log.Fatal("upload bucket name is not provided")
	}
}

func main() {
	config := &aws.Config{
		Region: aws.String(cliOpts.region),
	}
	if cliOpts.debug {
		log.Println("AWS debug logging is enabled")
		config.WithLogLevel(aws.LogDebug)
	}

	sess := session.New(config)
	uploader := s3manager.NewUploader(sess)

	for {
		if err := findAndUploadSnapshots(uploader); err != nil {
			log.Fatal(err)
		}

		if !cliOpts.watch {
			return
		}

		time.Sleep(time.Minute) // TODO: make configurable
	}
}

func findAndUploadSnapshots(uploader *s3manager.Uploader) error {
	entries, err := os.ReadDir(cliOpts.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(cliOpts.dir, entry.Name())
		_, err := os.Stat(filepath.Join(fullPath, snapshotName))
		if err != nil {
			log.Printf("skipping %v, no snapshot file found\n", entry.Name())
			continue
		}

		log.Println("processing", entry.Name())
		if err := checkSnapshotDir(fullPath); err != nil {
			log.Fatalf("snapshot %v check failed: %v", entry.Name(), err)
		}

		prefix := cliOpts.prefix
		if prefix != "" {
			prefix = filepath.Join(prefix, entry.Name())
		}

		iter := NewSyncFolderIterator(fullPath, cliOpts.bucket, prefix)
		if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
			log.Fatalf("unexpected error has occurred: %v", err)
		}

		if err := iter.Err(); err != nil {
			log.Fatalf("unexpected error occurred during file walking: %v", err)
		}

		if err := os.RemoveAll(fullPath); err != nil {
			log.Fatalf("cant remove directory %v: %v", fullPath, err)
		}
	}

	return nil
}

func checkSnapshotDir(path string) error {
	snapshot := abci.Snapshot{}
	chunksFound := uint32(0)

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), chunkSuffix) {
			chunksFound++
		}
	}

	data, err := ioutil.ReadFile(filepath.Join(path, snapshotName))
	if err != nil {
		return err
	}

	if err := snapshot.Unmarshal(data); err != nil {
		return err
	}

	if snapshot.Chunks != chunksFound {
		return fmt.Errorf("expected %v chunks but found %v", snapshot.Chunks, chunksFound)
	}

	return nil
}
