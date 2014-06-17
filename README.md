### Background

A few years ago I wrote a very simple [lua](http://www.lua.org/)
[tool](https://github.com/drio/drd.bio.toolbox/blob/master/lua/submit.lua) to
help me interact with [HPC clusters](http://en.wikipedia.org/wiki/HPCC) by
abstracting away the sending jobs details. I decided to rewrite that code
for various reasons.

First, I wanted to make it generic so code is generated for different clusters.
Also, I suspect this tool can capture more functionality required by the
community. By community I mean any user or programming interacting with
clusters, specifically in the domain of [bioinformatics](http://biostars.org).

Another, more academic reason is seeing how the code evolves as we make it
more generic and add new funtionality.

I could have continued using [lua](http://www.lua.org/) for this new version
but I decided to use [go](http://golang.org) instead. I think go extensive
[library](http://golang.org/pkg/) should facilitate the implementation of new
features. Also, I wanted fast loading times. Other programming languages may
take seconds to load, specially in a cluster environment where files are stored
in network storage.

### Install

Installation requires a working Go build [environment](http://golang.org/doc/install.html).
I will offer precompiled binaries in the near future though.

Example for a i386 linux box (without go distribution installed):

```sh
$ wget http://golang.org/dl/go1.2.2.linux-386.tar.gz
$ tar zxvf go1.2.2.linux-386.tar.gz
$ mkdir $HOME/gocode
$ export GOROOT=$HOME/go
$ export GOPATH=$HOME/gocode
$ export PATH=$PATH:$GOROOT/bin
$ export PATH=$PATH:$GOPATH/bin
$ go get github.com/drio/hpcgo
$ go install github.com/drio/hpcgo
$ hpcgo
```

If you have go distribution installed:

```sh
$ go get github.com/drio/hpcgo
$ go install github.com/drio/hpcgo
$ hpcgo
```



### What can you do with hpcgo?

```sh
  # Iterate over a bunch of files and hgcgo a jobs using each file
  $ ls *.fastq | xargs -i hgcgo -s fmi.{} -m 8G -c 8 "sga index --no-reverse -d 5000000 -t 8 {}"

  # Similar to before but using shell's for
  $ for i in `ls ../input/reads.*.fastq`; do F=`basename $i .fastq`; hgcgo -s pp.$F "sga preprocess -o $F.pp.fastq --pe-mode 2 $i"; done

  # hgcgo a two jobs, the second one has to run after the first one completes
  $ hgcgo -s one "touch ./one.txt" | bash > /tmp/deps.txt ; hgcgo -s two -f /tmp/deps.txt  "sleep 2;touch ./two.txt" | bash

  # Same as before but now we specify the jobid in the command line instead in a file
  $ hgcgo -s filter -d 3678650.sug-moab -m 20G -c 6 "sga fm-merge -m 65 -t 6 final.filter.pass.fa"

  # And my favourite one.
  # The second command reads the jobid from the standard input and uses it as dep
  $ (hgcgo -s one "sleep 15; touch ./one.txt" | bash) | hgcgo -s two -f -  "sleep 2;touch ./two.txt" | bash
```





