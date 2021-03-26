#!/bin/bash

TOPIC="test"

echo "Publishing incrementing numbers to topic $TOPIC until stopped..."

i=0

while true
do
    ((i=i+1))
    gcloud pubsub topics publish $TOPIC --message="PID: $$, i: $i"
done
