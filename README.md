# tm-snapshot-uploader

A tool to upload TM snapshots to S3

## Usage

```
Usage of ./build/tm-snapshot-uploader:
  -bucket string
    	AWS S3 bucket name
  -dir string
    	Directory containing height snapshots
  -prefix string
    	Prefix name for uploads in the bucket
  -region string
    	AWS region (default "us-east-1")
  -watch
    	Run in watch mode (periodically scan dir)
```
