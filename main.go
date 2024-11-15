package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

type KeyValue struct {
	Key   string
	Value string
}

type Memtable struct {
	Data map[string]string
}

func NewMemtable() *Memtable {
	return &Memtable{Data: make(map[string]string)}
}

func (m *Memtable) Put(key, value string) {
	m.Data[key] = value
}

// Serialize the Memtable to an SSTable file
func (m *Memtable) FlushToDisk(filename string) error {
	// Sort the keys
	keys := make([]string, 0, len(m.Data))
	for key := range m.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Write the sorted key-value pairs to disk
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)
	for _, key := range keys {
		err := encoder.Encode(KeyValue{Key: key, Value: m.Data[key]})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func SearchSSTable(filename, key string) (string, bool) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening SSTable:", err)
		return "", false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)

	for {
		var kv KeyValue
		err := decoder.Decode(&kv)
		if err != nil {
			break
		}
		if kv.Key == key {
			return kv.Value, true
		}
	}

	return "", false
}

func main() {
	memtable := NewMemtable()

	memtable.Put("apple", "fruit")
	memtable.Put("carrot", "vegetable")
	memtable.Put("banana", "fruit")

	sstableFile := "sstable.gob"

	err := memtable.FlushToDisk(sstableFile)
	if err != nil {
		fmt.Println("Error flushing to SSTable:", err)
		return
	}
	fmt.Println("Flushed Memtable to SSTable:", sstableFile)

	queryKey := "carrot"

	value, found := SearchSSTable(sstableFile, queryKey)
	if found {
		fmt.Printf("Found key '%s' with value '%s' in SSTable\n", queryKey, value)
	} else {
		fmt.Printf("Key '%s' not found in SSTable\n", queryKey)
	}
}
