package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	threads  = flag.Int("threads", 1, "Number of threads")
	limit    = flag.String("limit", "", "Download throughput limit. Example: 100k - 100 Kb/s; 1000 - 1000 b/s")
	urlsFile = flag.String("file", "", "Path to file wich contains list of files to download")
	outdir   = flag.String("outdir", "downloaded", "Path to directory where the downloaded files will be placed")
	wg       sync.WaitGroup
)

type Downloader struct {
	jobs   *chan string
	limit  int64
	outdir string
}

func NewDownloader(ch *chan string, limit, outdir string) *Downloader {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}
	multiplier := int64(1)
	switch limit[len(limit)-1] {
	case 'k', 'K':
		multiplier = int64(1024)
	case 'm', 'M':
		multiplier = int64(1024 * 1024)
	}
	limitBytes, err := strconv.ParseInt(limit[:len(limit)-1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	limitBytes = limitBytes * multiplier
	dnl := Downloader{jobs: ch, limit: limitBytes, outdir: path.Clean(outdir)}
	return &dnl
}

func (d *Downloader) Run() {
	for job := range *d.jobs {
		strArr := strings.Split(job, " ")
		url, filename := strArr[0], strArr[1]
		outfile := path.Join(d.outdir, filename)
		log.Printf("[ %s -> %s ] - downloading...", url, outfile)
		err := d.Get(url, outfile)
		if err != nil {
			log.Printf("[ %s -> %s ] - done!", url, outfile)
		}
	}
	wg.Done()
}

func (d *Downloader) Get(url, outfile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	of, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer of.Close()
	for range time.Tick(time.Second) {
		_, err := io.CopyN(of, resp.Body, d.limit)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *limit == "" || *urlsFile == "" {
		flag.Usage()
		log.Fatal("Missing flags!")
	}
	jobs := make(chan string)
	downl := NewDownloader(&jobs, *limit, *outdir)
	for i := 1; i <= *threads; i++ {
		wg.Add(1)
		go downl.Run()
	}
	f, err := os.Open(*urlsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		jobs <- scanner.Text()
	}
	close(jobs)
	wg.Wait()
}
