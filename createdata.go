package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
)

func map2sortslice(input map[uint64]int) []struct {
	a uint64
	b int
} {
	var tmp = make([]struct {
		a uint64
		b int
	}, 0, 1024)

	for k, v := range input {
		tmp = append(tmp, struct {
			a uint64
			b int
		}{k, v})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].b > tmp[j].b
	})

	return tmp
}

func createdata() {

	NewSource := rand.NewSource(2019)
	r := rand.New(NewSource)
	zipf := rand.NewZipf(r, 1.1, 1.2, 10000)

	var store = make(map[uint64]int)

	buf := bytes.NewBuffer(nil)
	for index := 0; index < 10000; index++ {
		val := zipf.Uint64()
		store[val]++
		fmt.Fprintf(buf, "%d ", val)
	}
	ioutil.WriteFile("data.txt", buf.Bytes(), 0755)
	buf.Reset()
	ss := map2sortslice(store)
	for i, _ := range ss {
		fmt.Fprintln(buf, ss[i].a, ss[i].b)
	}
	f, _ := os.OpenFile("sort.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	defer f.Close()
	buf.WriteTo(f)
}
