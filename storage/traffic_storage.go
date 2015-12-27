package storage

import (
	"log"
	"fmt"
	"time"

	"encoding/json"

	bolt "github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
	"io"
)

const (
	STORAGE_INMEMORY = "In memory"
	STORAGE_ONDISK = "On disk"

	BOLT_BUCKET= "SmartProxy"
)


type TrafficStorage struct {
	nature 		string // STORAGE_INMEMORY, STORAGE_ONDISK
	db 			*bolt.DB // database
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
	// Open the datafile in current directory, is created if it doesn't exist.
	dbFile := "capture.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Printf("[FATAL] cannot create database file: %s\n", dbFile)
		log.Fatal(err)
	}

	// Let's create the bucket if does not pre-exists
	db.Update(func(tx *bolt.Tx) error {
		//log.Printf("[STORAGE] creating bucket %s\n", BOLT_BUCKET)
		_, err := tx.CreateBucket([]byte(BOLT_BUCKET))
		if err != nil {
			//log.Printf("[STORAGE] cannot create bucket %s\n", BOLT_BUCKET)
			return fmt.Errorf("create bucket: %s", err)
		}
		log.Printf("[INFO] Created bucket %s to persist traffic captures\n", BOLT_BUCKET)
		return nil
	})

	return &TrafficStorage{STORAGE_INMEMORY, db}
}


func (storage *TrafficStorage) close() {
	log.Printf("[INFO] Closing storage database")
	storage.db.Close()
}


func (storage *TrafficStorage) CreateTrace() *TrafficTrace {
	trace := new(TrafficTrace)
	trace.ID = uuid.NewV4().String()

	log.Printf("[DEBUG] STORAGE Created new trace with id: %s\n", trace.ID)

    return trace
}

func (storage *TrafficStorage) StoreTrace(trace *TrafficTrace) {
	log.Printf("[DEBUG] STORAGE Storing trace ", trace.ID)

	storage.db.Update(func(tx *bolt.Tx) error {
		encoded, err1 := json.Marshal(trace)
		if err1 != nil {
			log.Printf("[WARNING] STORAGE Cannot encode trace with id: %s\n", trace.ID)
			return err1
		}

		b := tx.Bucket([]byte(BOLT_BUCKET))
		err2 := b.Put([]byte(trace.ID), encoded)
		if err2 != nil {
			log.Printf("[WARNING] STORAGE Error while storing capture trace with id: %s\n", trace.ID)
		}
		return err2
	})
}

func (storage *TrafficStorage) GetTraces() int {
	_ = "breakpoint"
	log.Printf("[DEBUG] STORAGE fetching all traces\n")

	count := 0
	storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLT_BUCKET))
		b.ForEach(func(k, v []byte) error {
			log.Printf("[STORAGE] key=%s, value=%s\n", k, v)
			count++
			return nil
		})
		return nil
	})

	return count
}



func (storage *TrafficStorage) DisplayLastTraces(w io.Writer, max int) int {
	log.Printf("[STORAGE] fetching last traces, %d max\n", max)

	count := 0

	storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLT_BUCKET))
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			fmt.Fprintf(w, "<p>Captured : id=%s, value=%s</p>\n", k, v)
			if count >= max {
				break;
			}
		}
		return nil
	})

	return count
}









