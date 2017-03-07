#!/bin/bash

data_dir=$1
# mkdir -p "${data_dir}"

cd "${data_dir}"

# touch some random files
touch file.txt
touch blah.txt
touch fuf.elf
touch .ds_store


mkdir test_data/

/vagrant/setup_test.sh /vagrant/names.txt

