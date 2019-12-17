package main

/* This script counts kmers from a provided file in sequence data it finds on stdin. */

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Args struct {
	Limit int
	Kmers string
	Size  int
}

var args = Args{}

func init() {
	log.SetFlags(0)
	flag.StringVar(&args.Kmers, "kmers", "", "file listing kmers to look for")
	flag.IntVar(&args.Limit, "limit", 0, "limit the number of lines of stdin to consider (default = 0 = unlimited)")

	flag.Usage = func() {
		log.Println("usage: kmercount [options] < readsfile")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if args.Kmers == "" {
		log.Println("Missing -kmers parameter")
		flag.Usage()
		os.Exit(1)
	}

	kmers := make(map[string]int)
	kmerSize := 0

	fp, err := os.Open(args.Kmers)
	if err != nil {
		log.Fatal("Failed to open kmers file")
	}
	kmerScan := bufio.NewScanner(fp)
	for kmerScan.Scan() {
		kmer := kmerScan.Text()
		if len(kmer) == 0 {
			log.Fatal("Empty line in kmers file")
		}
		if kmerSize == 0 {
			kmerSize = len(kmer)
		}
		if len(kmer) != kmerSize {
			log.Fatalf("Expecting kmer to be length %d but %s is length %d\n", kmerSize, kmer, len(kmer))
		}
		kmers[kmer] = 0
	}

	log.Println("Found", len(kmers), "kmers")

	seqScan := bufio.NewScanner(os.Stdin)
	lineNum := 0
	b := ""
	offset := 0
	for seqScan.Scan() {
		if args.Limit > 0 && lineNum > args.Limit {
			break
		}
		b = b + strings.ToUpper(seqScan.Text())
		lineNum++
		for offset+kmerSize <= len(b) {
			kmer := string(b[offset:(offset + kmerSize)])
			_, present := kmers[kmer]
			if present {
				kmers[kmer] += 1
			}
			offset++
		}
		b = string(b[offset:len(b)])
		offset = 0
	}

	for kmer, count := range kmers {
		fmt.Printf("%s\t%d\n", kmer, count)
	}
}
