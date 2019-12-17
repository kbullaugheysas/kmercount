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
      -kmers string
        	file listing kmers to look for
      -limit int
        	limit the number of lines of stdin to consider (default = 0 = unlimited)
