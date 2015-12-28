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
	nature 		string 			// STORAGE_INMEMORY, STORAGE_ONDISK
	db 			*bolt.DB 		// database
	cursor		*bolt.Cursor 	// filled if a cursor is actually opened, nil otherwise
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

func OnDiskTrafficStorage () *TrafficStorage {
	// Open the datafile in current directory, is created if it doesn't exist.
	dbFile := "capture.db"
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Printf("[FATAL] STORAGE cannot create database file: %s\n", dbFile)
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
		log.Printf("[INFO] STORAGE Created bucket %s to persist traffic captures\n", BOLT_BUCKET)
		return nil
	})

	return &TrafficStorage{STORAGE_ONDISK, db, nil}
}


func (storage *TrafficStorage) close() {
	log.Printf("[INFO] STORAGE Closing database")
	storage.db.Close()
}


func (storage *TrafficStorage) CreateTrace() *TrafficTrace {
	trace := new(TrafficTrace)
	trace.ID = uuid.NewV1().String()

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
			log.Printf("[DEBUG] STORAGE key=%s, value=%s\n", k, v)
			count++
			return nil
		})
		return nil
	})

	return count
}



func (storage *TrafficStorage) DisplayLatestTraces(w io.Writer, max int) int {
	log.Printf("[DEBUG] STORAGE fetching last traces, %d max\n", max)

	count := 0
	storage.cursor = nil

	storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BOLT_BUCKET))
		c := b.Cursor()
		for k, _ := c.Last(); k != nil; k, _ = c.Prev() {
			fmt.Fprintf(w, "<p>Captured : %s</p>\n", k)
			count++
			if count >= max {
				break;
			}
		}
		return nil
	})

	if count == 0 {
		// No traffic captured yet
		fmt.Fprintf(w, "<p>No traffic so far</p>\n")
	}

	return count
}


func (storage *TrafficStorage) DisplayLastTrace(w io.Writer) {
	log.Printf("[DEBUG] STORAGE fetching last trace\n")

	var trace TrafficTrace

	storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BOLT_BUCKET))
		cursor := bucket.Cursor()
		key, bytes := cursor.Last()

		if key == nil {
			log.Printf("[DEBUG] STORAGE last trace not found, no capture presumed\n")
			fmt.Fprintf(w, "<p>No traffic capture found</p>")
			return nil
		}

		storage.cursor = cursor
		err := json.Unmarshal(bytes, &trace)
		if err != nil {
			log.Printf("[DEBUG] STORAGE json decode failed for bytes %s\n", bytes)
			fmt.Fprintf(w, "<p>Cannot read traffic data</p>")
			return nil
		}

		return nil
	})

	displayTrace(w, &trace)
	return
}


func (storage *TrafficStorage) DisplayFirstTrace(w io.Writer) {
	log.Printf("[DEBUG] STORAGE fetching first trace\n")

	var trace TrafficTrace

	storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BOLT_BUCKET))
		cursor := bucket.Cursor()
		key, bytes := cursor.First()

		if key == nil {
			log.Printf("[DEBUG] STORAGE first trace not found, no capture presumed\n")
			fmt.Fprintf(w, "<p>No traffic capture found</p>")
			return nil
		}

		storage.cursor = cursor
		err := json.Unmarshal(bytes, &trace)
		if err != nil {
			log.Printf("[DEBUG] STORAGE json decode failed for bytes %s\n", bytes)
			fmt.Fprintf(w, "<p>Cannot read traffic data</p>")
			return nil
		}

		return nil
	})

	displayTrace(w, &trace)
	return
}


func (storage *TrafficStorage) DisplayPrevTrace(w io.Writer) {
	log.Printf("[DEBUG] STORAGE fetching first trace\n")

	if (storage.cursor == nil) {
		log.Printf("[DEBUG] STORAGE cursor not initialized, going to first trace\n")
		storage.DisplayFirstTrace(w)
		return
	}

	key, bytes := storage.cursor.Prev()
	if key == nil {
		log.Printf("[DEBUG] STORAGE first trace not found, no capture presumed\n")
		fmt.Fprintf(w, "<p>No traffic capture found</p>")
		return
	}

	var trace TrafficTrace
	err := json.Unmarshal(bytes, &trace)
	if err != nil {
		log.Printf("[DEBUG] STORAGE json decode failed for bytes %s\n", bytes)
		fmt.Fprintf(w, "<p>Cannot read traffic data</p>")
		return
	}

	displayTrace(w, &trace)
	return
}


func (storage *TrafficStorage) DisplayNextTrace(w io.Writer) {
	log.Printf("[DEBUG] STORAGE fetching next trace\n")

	if (storage.cursor == nil) {
		log.Printf("[DEBUG] STORAGE cursor not initialized, going to first trace\n")
		storage.DisplayFirstTrace(w)
		return
	}

	key, bytes := storage.cursor.Next()
	if key == nil {
		log.Printf("[DEBUG] STORAGE first trace not found, no capture presumed\n")
		fmt.Fprintf(w, "<p>No traffic capture found</p>")
		return
	}

	var trace TrafficTrace
	err := json.Unmarshal(bytes, &trace)
	if err != nil {
		log.Printf("[DEBUG] STORAGE json decode failed for bytes %s\n", bytes)
		fmt.Fprintf(w, "<p>Cannot read traffic data</p>")
		return
	}

	displayTrace(w, &trace)
	return
}


func displayTrace(w io.Writer, trace *TrafficTrace) {
	fmt.Fprintf(w, "<p>ID : %s</p>", trace.ID)
	fmt.Fprintf(w, "<p>Method : %s</p>", trace.HttpMethod)
	fmt.Fprintf(w, "<p>URI : %s</p>", trace.URI)
	fmt.Fprintf(w, "<p>Status : %v</p>", trace.HttpStatus)
	start := time.Time(trace.Start)
	end := time.Time(trace.End)
	fmt.Fprintf(w, "<p>Duration : %v</p>", end.Sub(start))
	fmt.Fprintf(w, "<p>Started at : %v</p>", start)
	fmt.Fprintf(w, "<p>Completed at : %v</p>", end)
	fmt.Fprintf(w, "<p>Outgoing size : %v bytes</p>", trace.Length)

	fmt.Fprintf(w, "<p>Incoming : %s</p>", string(*trace.Ingress.Bytes))
	fmt.Fprintf(w, "<p>Outgoing : %s</p>", string(*trace.Egress.Bytes))
}














