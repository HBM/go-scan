package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/browser"
)

func NewDB() *DB {
	return &DB{
		table: make(map[string]Announcement),
	}
}

type DB struct {
	mutex sync.Mutex
	table map[string]Announcement
}

func (db *DB) Set(ip string, ann Announcement) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.table[ip] = ann
}

func (db *DB) Delete(ip string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	delete(db.table, ip)
}

func (db *DB) Snapshot() map[string]Announcement {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.table
}

var (
	timers = map[string]*time.Timer{}
)

const (
	udp     = "udp"
	address = "239.255.77.76:31416"
)

type Announcement struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	APIVersion string    `json:"apiVersion"`
	Device     Device    `json:"device"`
	Services   []Service `json:"services"`
	Expiration int       `json:"expiration"`
}

type Service struct {
	Type string `json:"type"`
	Port int    `json:"port"`
}

type Device struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Label           string `json:"label"`
	FamilyType      string `json:"familyType"`
	FirmwareVersion string `json:"firmwareVersion"`
	IsRouter        bool   `json:"isRouter"`
}

func main() {

	flag.Parse()

	db := NewDB()

	udpAddr, err := net.ResolveUDPAddr(udp, address)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP(udp, nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// add file server
	
	fmt.Println("HBM scan: serving files from local file system")
	http.Handle("/", http.FileServer(http.Dir("public")))


	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db.Snapshot())
	})

	// start server
	go func() {
		fmt.Println("HBM scan: web server running at http://localhost:8085. ctrl+c to stop it.")
		log.Fatal(http.ListenAndServe(":8085", nil))
	}()

	// open default browser
	browser.Stdout = nil
	if err := browser.OpenURL("http://localhost:8085"); err != nil {
		panic(err)
	}

	// start listening for announcements
	for {
		// create new buffer
		b := make([]byte, 1500)

		// read message from udp into buffer
		n, src, err := conn.ReadFromUDP(b)
		if err != nil {
			panic(err)
		}

		// convert raw json bytes to struct
		var ann Announcement
		if err := json.Unmarshal(b[:n], &ann); err != nil {
			panic(err)
		}

		// add announcement to db
		ip := src.IP.String()
		db.Set(ip, ann)

		// check for existing timer
		timer, ok := timers[ip]
		if ok {
			delete(timers, ip)
			timer.Stop()
		}

		// start new timer for device
		timer = time.AfterFunc(time.Second*time.Duration(ann.Params.Expiration), func() {
			delete(timers, ip)
			db.Delete(ip)
		})

		// store timer in timers db
		timers[ip] = timer
	}
}
