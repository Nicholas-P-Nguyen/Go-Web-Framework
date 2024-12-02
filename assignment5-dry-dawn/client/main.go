package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"cos316.princeton.edu/assignment5/conn_pool"
)

var wg sync.WaitGroup

func acquireConnection(p *conn_pool.Pool) conn_pool.Conn {
	defer wg.Done()
	conn := p.Open()

	return conn
}

func getMovieTitleYear(year int, p *conn_pool.Pool) {
	defer wg.Done()
	fmt.Println("get year", year)
	conn := p.Open()

	fmt.Println(year)

	query := "select title from Movies where year = " + strconv.Itoa(year)
	var title string

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		rows.Scan(&title)
		fmt.Println(title)
	}
	p.Close(conn)

}

func holdForASec(p *conn_pool.Pool) {
	defer wg.Done()
	conn := p.Open()

	time.Sleep(time.Second)

	p.Close(conn)
}

func main() {
	pool, err := conn_pool.NewPool(Connector())
	if err != nil {
		fmt.Println(err)
		return
	}

	pool.SetMaxConnections(5)
	conn := pool.Open()
	var title, genre string

	rows, err := conn.Query("select title, genres from Movies where year = 1933;")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		rows.Scan(&title, &genre)
		fmt.Println(title, genre)
	}

	pool.Close(conn)

	startYear := 1933
	span := 10

	for i := 0; i < span; i++ {
		wg.Add(1)
		go getMovieTitleYear(startYear+i, pool)
	}

	for i := 0; i < span; i++ {
		wg.Add(1)
		go holdForASec(pool)
	}

	wg.Wait()
}
