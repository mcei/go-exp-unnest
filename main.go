package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// docker run --name bulk-postgres -e POSTGRES_PASSWORD=bulk -p 5432:5432 postgres
//
// docker start bulk-postgres

const (
	dbURL    = "postgres://postgres:bulk@0.0.0.0:5432/postgres?sslmode=disable"
	usersNum = 25_000
	//usersNum = 67_000
	//usersNum = 150_000
)

func main() {
	conn := setup()
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	values := makeTestUsers(usersNum)

	// обычная вставка, используя неоптимальное формирование строки
	RegularInsert(conn, values)

	// обычная вставка, используя strings.Builder
	RegularInsertSB(conn, values)

	// вставка, используя unnest
	UnnestInsert(conn, values)
}

func setup() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
	 id int,
	 name int
	);
	TRUNCATE TABLE users;
	`

	if _, err = conn.Exec(context.Background(), query); err != nil {
		panic(err)
	}

	return conn
}

// return a 2d array of int pairs
func makeTestUsers(max int) [][]int {
	arr := make([][]int, max+1)
	for i := range arr {
		arr[i] = []int{i, i + 1}
	}

	// [[0,1],[1,2], ..., [25000, 25001]]
	return arr
}

// converts a 2d array to query numbers and values
// input: [[1,2],[3,4]]
// output: ($1,$2),($3,$4) and [1,2,3,4]
func valuesToRows(values [][]int) (string, []interface{}) {
	start := time.Now()
	defer func() { fmt.Println("valuesToRows:", time.Since(start)) }()

	var rows []interface{}
	query := ""
	for i, s := range values {
		rows = append(rows, s[0], s[1])

		numFields := 2
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]

	return query, rows

}

func valuesToRowsStringBuilder(values [][]int) (string, []interface{}) {
	start := time.Now()
	defer func() { fmt.Println("valuesToRowsStringBuilder:", time.Since(start)) }()

	var rows []interface{}
	var query strings.Builder

	valuesNum := len(values)
	for i, s := range values {
		rows = append(rows, s[0], s[1])

		numFields := 2
		n := i * numFields

		query.WriteString(`(`)
		for j := 0; j < numFields; j++ {
			query.WriteString(`$`)
			query.WriteString(strconv.Itoa(n + j + 1))
			if j < numFields-1 {
				query.WriteString(`,`) // Drop last `,`
			}
			//query.WriteString(`,`)
		}
		query.WriteString(`)`)
		if i < valuesNum-1 {
			query.WriteString(`,`) // Drop last `,`
		}
	}

	return query.String(), rows
}

func RegularInsert(conn *pgx.Conn, values [][]int) {
	start := time.Now()
	defer func() { fmt.Println("RegularInsert", time.Since(start)) }()

	query := `
  INSERT INTO users
    (id, name)
    VALUES %s;`

	queryParams, params := valuesToRows(values)
	query = fmt.Sprintf(query, queryParams)

	if _, err := conn.Exec(context.Background(), query, params...); err != nil {
		fmt.Println(err)
	}
}

func RegularInsertSB(conn *pgx.Conn, values [][]int) {
	start := time.Now()
	defer func() { fmt.Println("RegularInsertSB", time.Since(start)) }()

	query := `
  INSERT INTO users
    (id, name)
    VALUES %s;`

	queryParams, params := valuesToRowsStringBuilder(values)
	query = fmt.Sprintf(query, queryParams)

	if _, err := conn.Exec(context.Background(), query, params...); err != nil {
		fmt.Println(err)
	}
}

func UnnestInsert(conn *pgx.Conn, values [][]int) {
	start := time.Now()
	defer func() { fmt.Println("UnnestInsert", time.Since(start)) }()

	var ids []int
	var names []int
	for _, v := range values {
		ids = append(ids, v[0])
		names = append(names, v[1])
	}

	query := `
  INSERT INTO users
    (id, name)
    (
      select * from unnest($1::int[], $2::int[])
    )`

	if _, err := conn.Exec(context.Background(), query, ids, names); err != nil {
		fmt.Println(err)
	}
}
