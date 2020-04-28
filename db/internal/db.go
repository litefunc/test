package internal

import (
	"cloud/database/postgresql"
	"cloud/lib/logger"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

const query = `SELECT * FROM cloud.device_stat_vod`

func conn(n int) {
	for i := 1; i <= n; i++ {

		db, err := postgresql.Connect(dbConfig())
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debug(i, db)
	}
}

func read(db *sql.DB, i int) {
	rows, err := db.Query(query)
	if err != nil {
		logger.Error(err)
		return
	}
	defer rows.Close()
	// fmt.Println(i)
}

func goread(db *sql.DB, i int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := db.Query(query)
		if err != nil {
			logger.Error(err)
			return
		}
		defer rows.Close()
		// fmt.Println(i)
	}()
}

func chread(db *sql.DB, in, out chan int) {

	go func() {
		for i := range in {
			func() {

				rows, err := db.Query(query)
				if err != nil {
					logger.Error(err)
					return
				}
				defer rows.Close()
				out <- i
			}()

		}
	}()
}

func Read() {

	db, err := sql.Open("postgres", dbConfig())
	if err != nil {
		logger.Error("db connection", err)
		return
	}

	if err := db.Ping(); err != nil {
		logger.Error(err)
		return
	}

	n := 200

	t1 := time.Now()
	for i := 1; i <= n; i++ {
		read(db, i)
	}
	t2 := time.Now()
	d1 := t2.Sub(t1)

	fmt.Println("---------------------------------------")

	wg := &sync.WaitGroup{}

	t1 = time.Now()
	for i := 1; i <= n; i++ {
		goread(db, i, wg)
	}
	wg.Wait()
	t2 = time.Now()
	d2 := t2.Sub(t1)

	fmt.Println("---------------------------------------")
	in := make(chan int, 50)
	out := make(chan int)
	for i := 1; i <= 2; i++ {
		chread(db, in, out)
	}

	t1 = time.Now()
	go func() {
		for i := 1; i <= n; i++ {
			in <- i
		}
		defer close(in)
	}()
	for i := 1; i <= n; i++ {
		// fmt.Println(<-out)
		<-out
	}
	t2 = time.Now()
	d3 := t2.Sub(t1)

	fmt.Println(d1, d2, d3)

	// conn(50)

	var wc chan int
	<-wc
}
