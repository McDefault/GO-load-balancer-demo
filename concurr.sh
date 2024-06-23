#!/bin/bash
for i in {1..1000000}
do
   echo "This is loop iteration number $i"
   ./client.exe &
done
