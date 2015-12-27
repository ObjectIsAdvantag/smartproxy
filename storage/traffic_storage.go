package storage

import (
	"log"
	"fmt"
	"time"

	"encoding/json"

	bolt "github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
)

const (
	STORAGE_INMEMORY = "In memory"
	STORAGE_ONDISK = "On disk"

	BOLT_BUCKET= "SmartProxy"
)


type TrafficStorage struct {
	nature 		string // STORAGE_INMEMORY, STORAGE_ONDISK
	db 			*bolt.DB // database
	bucket  	*bolt.Bucket
}

type TrafficTrace struct {
	ID			string // unique identifier of the trace
	Start		time.Time
	End			time.Time
	HttpStatus 	int
	HttpMethod	string
	URI 		string
	Length 		int
	Ingress 	*TrafficIngress
	Egress 		*TrafficEgress
}

type TrafficIngress struct {
	Bytes 		*[]byte
    //Headers 	[]http.Header
	//Body 		string
}

type TrafficEgress struct {
	Bytes 		*[]byte
}

func VolatileTrafficStorage () *TrafficStorage {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbFile := "capture.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Printf("[STORAGE] cannot create database file: %s", dbFile)
		log.Fatal(err)
	}
	var bucket *bolt.Bucket
	db.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucket([]byte(BOLT_BUCKET))
		if err != nil {
			log.Printf("[STORAGE] cannot create bucket: %s", BOLT_BUCKET)
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return &TrafficStorage{STORAGE_INMEMORY, db, bucket}
}


func (storage *TrafficStorage) close() {
	log.Printf("[STORAGE] closing storage database")
	storage.db.Close()
}


func (storage *TrafficStorage) CreateTrace() *TrafficTrace {
	trace := new(TrafficTrace)
	trace.ID = uuid.NewV4().String()

	log.Printf("[STORAGE] created new trace with id: %s\n", trace.ID)

    return trace
}

func (storage *TrafficStorage) StoreTrace(trace *TrafficTrace) {
	log.Printf("[STORAGE] storing trace with id: %s\n", trace.ID)

	storage.db.Update(func(tx *bolt.Tx) error {
		encoded, err1 := json.Marshal(trace)
		if err1 != nil {
			return err1
		}
		err2 := storage.bucket.Put([]byte(trace.ID), encoded)
		return err2
	})
}






