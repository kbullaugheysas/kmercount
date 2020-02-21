## kmercount

This utility is for counting the existence of known kmers (from a list) in
sequencing data. This allows one to compute an upperbound for the amount of
sequence contained in a run that is from some particular experimental origin,
such as barcodes, vectors, primers, the t7 promoter, linkers, etc. that one
might want to quantify.

### Usage

You can see the usage as follows:

    kmercount -help

Which produces output like this:

    usage: kmercount [options] < readsfile
      -all int
            count all kmers of given size (default = 0, only count kmers in -kmers file)
      -fasta string
            Accept input from this fasta file
      -kmers string
            file listing kmers to look for
      -limit int
            limit the number of lines of stdin to consider (default = 0 = unlimited)
      -unroll
            Treat lines as one long contiguous sequence

### Examples

The program can take a list of sequences on stdin. By default each line is
treated separately and kmers can't span line breaks:

    (echo "AABCAAABCAA"; echo "ABCAAAAA") | kmercount -kmers <(echo "ABC"; echo "AAA")

The program produces tabular output, and for the above command we'd get this:

    kmer    ABC     3
    kmer    AAA     4
    stat    queries 2
    stat    comparisons     15
    stat    sum     7

But if we add the -unroll parameter with the same inputs as above, like this:

    (echo "AABCAAABCAA"; echo "ABCAAAAA") | kmercount -kmers <(echo "ABC"; echo "AAA") -unroll

We get an extra count of "AAA" that spands the line break:

    kmer    ABC     3
    kmer    AAA     5
    stat    queries 2
    stat    comparisons     17
    stat    sum     8

It's also possible to count all kmers of a certain size using the `-all` parameter:

    (echo "AABCAAABCAA"; echo "ABCAAAAA") | kmercount -all 3

For which we observe the following output:

    kmer    AAB     2
    kmer    ABC     3
    kmer    BCA     3
    kmer    CAA     3
    kmer    AAA     4
    stat    queries 5
    stat    comparisons     15
    stat    sum     15




