### Background

Submitting jobs is a common task all my software and analysis pipelines have to
do when working in a [HPC cluster](http://en.wikipedia.org/wiki/HPCC).
Abstracting away the HPC details facilitates the porting of pieces of software
from different HPC implementations. This tool implements that.

We tell the software, via passing arguments what command we want to run, with
what resources (memory, number of cores, queue to use ...)  and what
dependencies we have. We also specify what hpc implementation we want the tool
to use. With that information, submit will generate the necessary command we
should run in order to submit the job to the cluster. The examples should help
us understand this:

```sh
submit.go submit -s name [-b <backend> -d <dep> -m <mem> -c <num_cores>] <cmd>
```

### Examples:

```sh
  # Iterate over a bunch of files and submit a jobs using each file
  $ ls *.fastq | xargs -i submit -s fmi.{} -m 8G -c 8 "sga index --no-reverse -d 5000000 -t 8 {}"

  # Similar to before but using shell's for
  $ for i in `ls ../input/reads.*.fastq`; do F=`basename $i .fastq`; submit -s pp.$F "sga preprocess -o $F.pp.fastq --pe-mode 2 $i"; done

  # Submit a two jobs, the second one has to run after the first one completes
  $ submit -s one "touch ./one.txt" | bash > /tmp/deps.txt ; submit -s two -f /tmp/deps.txt  "sleep 2;touch ./two.txt" | bash

  # Same as before but now we specify the jobid in the command line instead in a file
  $ submit -s filter -d 3678650.sug-moab -m 20G -c 6 "sga fm-merge -m 65 -t 6 final.filter.pass.fa"

  # And my favourite one.
  # The second command reads the jobid from the standard input and uses it as dep
  $ (submit -s one "sleep 15; touch ./one.txt" | bash) | submit -s two -f -  "sleep 2;touch ./two.txt" | bash
```





