/* access_record.go - recording access count and last access time */
/*
modification history
--------------------
2015/09/02, by Guang Yao, create
*/
/*
DESCRIPTION
*/

package access_record

import (
	"sync"
	"time"
)

type Record struct {
	Count    int64
	LastTime string
}

type AccessRecord struct {
	lock    sync.RWMutex       // protecting the records
	records map[string]*Record // record key => Record
}

// generate a new AccessRecord
func NewAccessRecord() *AccessRecord {
	r := new(AccessRecord)

	r.records = make(map[string]*Record)

	return r
}

// inc a record and update last access time
func (r *AccessRecord) Inc(key string, t time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// fill the key
	if _, exist := r.records[key]; !exist {
		r.records[key] = &Record{0, ""}
	}

	// update record
	record, _ := r.records[key]
	record.Count += 1
	record.LastTime = t.Format("2006/01/02 15:04:05")
}

// get a record
func (r *AccessRecord) Get(key string) *Record {
	r.lock.RLock()
	defer r.lock.RUnlock()

	record, exist := r.records[key]
	if !exist {
		return nil
	} else {
		return &Record{record.Count, record.LastTime}
	}
}

// get all the records
func (r *AccessRecord) GetAll() map[string]*Record {
	ret := make(map[string]*Record)

	r.lock.RLock()
	defer r.lock.RUnlock()

	// fill the table
	for key, record := range r.records {
		ret[key] = &Record{record.Count, record.LastTime}
	}

	return ret
}
