package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

const projectId = "fatreporting"
const bucket = "ultrafish95"
const credentials = "credentials.json"

func main() {
	// get current working dir
	path, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}

	// open dir
	dir, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}

	defer dir.Close()

	// scan dir
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Panic(err)
	}

	jsonBuffer, err := bufferFiles(files)
	if err != nil {
		log.Panic(err)
	}

	if err := uploadToGCS(projectId, bucket, jsonBuffer); err != nil {
		log.Panic(err)
	}
}

func bufferFiles(files []os.FileInfo) (*bytes.Buffer, error) {
	var sessions []*Session
	var buffer = &bytes.Buffer{}
	var totalBytes = 0
	var filesToMove []*os.File

	// loop through files, and parse them
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "csv") {
			file, err := os.Open(file.Name())
			if err != nil {
				return nil, err
			}

			session, err := parseSession(file)
			if err != nil {
				return nil, err
			}

			err = file.Close()
			if err != nil {
				return nil, err
			}

			filesToMove = append(filesToMove, file)
			sessions = append(sessions, session)
		}
	}

	// marshall sessions to json
	for _, session := range sessions {
		jsonInBytes, err := json.Marshal(session)
		if err != nil {
			return nil, err
		}

		sessionBytes, err := buffer.WriteString(string(jsonInBytes) + "\r\n")
		if err != nil {
			return nil, err
		}

		totalBytes = totalBytes + sessionBytes
	}

	defer cleanupCSV(filesToMove)
	log.Printf("%d bytes buffered", totalBytes)
	return buffer, nil
}

func uploadToGCS(projectId string, bucket string, buffer *bytes.Buffer) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucket)
	if err := bkt.Create(ctx, projectId, nil); err != nil {
		// creates bucket if it doesn't exist
	}

	s := uuid.New().String()
	obj := bkt.Object(s)
	w := obj.NewWriter(ctx)

	if _, err := w.Write(buffer.Bytes()); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func cleanupCSV(files []*os.File) {
	if len(files) > 0 {
		// make directory if it doesn't exist
		if err := os.MkdirAll("processed", 0755); err != nil {
			log.Panic(err)
		}

		// copy files to processed dir
		for _, file := range files {
			if err := os.Rename(file.Name(), "processed/"+file.Name()); err != nil {
				log.Panic(err)
			}
		}
	}
}
