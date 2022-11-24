#!/bin/bash

# TODO: MAKE THIS AN ANSIBLE VARIABLE
max_archive_size=1073741824 # 1 GB

# Make sure the destination directory has been specified.
dest_dir=$1
[ -n "${dest_dir}" ] || exit
shift

# Make sure that the batch size (in line numbers) has been properly specified.
max_lines=$1
[ "${max_lines}" -gt 0 ] || exit
shift

# Make sure that the maximal log archive size (in bytes) has been properly specified.
max_archive_size=$1
[ "${max_archive_size}" -gt 0 ] || exit
shift

# Initialize.
cd "${dest_dir}" || exit
lines=0
n=0
current_file=$(printf "%05d-%s.log" $n "$(date +%Y-%m-%d-%H-%M-%S_%Z)")

# Copy standard input to standard output AND to a file.
# Every $max_lines lines, compress the data output so far and start writing to a new file.
while IFS='$\n' read -r line; do

  # Copy input to output and increment line counter.
  echo "$line"
  echo "$line" >> "${current_file}"
  lines=$((lines + 1))

  # When maximal number of lines has been reached
  if [ $lines -ge "$max_lines" ]; then

    # Compress the current output.
    tar czf "${current_file}.tar.gz" "${current_file}"
    rm "${current_file}"

    # Reset / increment counters and create a new output file.
    lines=0
    n=$((n + 1))
    current_file=$(printf "%05d-%s.log" $n "$(date +%Y-%m-%d-%H-%M-%S_%Z)")

    # Delete oldest logs until space limit is not exceeded any more.
    while [ "$(du -bs . | awk '{print $1}')" -gt $max_archive_size ]; do
      if ls *.tar.gz; then
        # If there are any files to delete, delete the first one.
        rm "$(ls -tr1 | head -n 1)"
      else
        # If there are no files to delete, stop the loop
        break
      fi
    done

  fi

done
