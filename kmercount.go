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
	Limit  int
	Kmers  string
	All    int
	Fasta  string
	Size   int
	Unroll bool
}

var args = Args{}

var totalCount = 0

func init() {
	log.SetFlags(0)
	flag.StringVar(&args.Kmers, "kmers", "", "file listing kmers to look for")
	flag.IntVar(&args.All, "all", 0, "count all kmers of given size (default = 0, only count kmers in -kmers file)")
	flag.IntVar(&args.Limit, "limit", 0, "limit the number of lines of stdin to consider (default = 0 = unlimited)")
	flag.StringVar(&args.Fasta, "fasta", "", "Accept input from this fasta file")
	flag.BoolVar(&args.Unroll, "unroll", false, "Treat lines as one long contiguous sequence")

	flag.Usage = func() {
		log.Println("usage: kmercount [options] < readsfile")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if args.Kmers == "" && args.All == 0 {
		log.Println("Must supply -kmers argument or -all argument")
		flag.Usage()
		os.Exit(1)
	}
	if args.Kmers != "" && args.All > 0 {
		log.Println("can't supply both -kmers and -all arguments")
		flag.Usage()
		os.Exit(1)
	}

	kmers := make(map[string]int)
	kmerSize := 0
	if args.All > 0 {
		kmerSize = args.All
	} else {
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
		fp.Close()
	}

	lineNum := 0
	b := ""
	offset := 0
	if args.Fasta != "" {
		fp, err := os.Open(args.Fasta)
		if err != nil {
			log.Fatalf("Failed to open fasta file %s: %v", args.Fasta, err)
		}
		fastaScan := bufio.NewScanner(fp)
		for fastaScan.Scan() {
			if args.Limit > 0 && lineNum > args.Limit {
				break
			}
			line := fastaScan.Text()
			if strings.HasPrefix(line, ">") {
				// Skip fasta header lines, but reset the buffer
				log.Println(line)
				b = ""
				offset = 0
				continue
			}
			b = b + strings.ToUpper(line)
			lineNum++
			offset = countKmers(b, kmerSize, kmers)
			b = string(b[offset:len(b)])
		}
	} else {
		seqScan := bufio.NewScanner(os.Stdin)
		for seqScan.Scan() {
			if args.Limit > 0 && lineNum > args.Limit {
				break
			}
			b = b + strings.ToUpper(seqScan.Text())
			lineNum++
			offset = countKmers(b, kmerSize, kmers)
			if args.Unroll {
				b = string(b[offset:len(b)])
			} else {
				offset = 0
				b = ""
			}
		}
	}

	sum := 0
	for kmer, count := range kmers {
		fmt.Printf("kmer\t%s\t%d\n", kmer, count)
		sum += count
	}
	fmt.Printf("stat\tqueries\t%d\n", len(kmers))
	fmt.Printf("stat\tcomparisons\t%d\n", totalCount)
	fmt.Printf("stat\tsum\t%d\n", sum)
}

func countKmers(buf string, k int, kmers map[string]int) int {
	offset := 0
	for offset+k <= len(buf) {
		kmer := string(buf[offset:(offset + k)])
		_, present := kmers[kmer]
		if present {
			kmers[kmer] += 1
		} else {
			if args.All > 0 {
				kmers[kmer] = 1
			}
		}
		offset++
		totalCount++
	}
	return offset
}
