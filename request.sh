#!/bin/bash

# Capture the start time in nanoseconds
start_time=$(date +%s%N)

t=$1
while [ $t -gt 0 ];
do
t=$((t-1))
tokenid=$(( RANDOM % 20000 + 1 ))
participantcodeid="$(( RANDOM % 200 + 1 ))"
segmentid=$(( RANDOM % 20 + 1 ))
dealerid="$(( RANDOM % 100000 + 1 ))"

# echo "tokenid: $tokenid"
# echo "participantcodeid: $participantcodeid"
# echo "segmentid: $segmentid"
# echo "dealerid: $dealerid"
curl "http://10.10.198.204:3000/getScripMaster?tokenid=$tokenid&marketsegmentid=$segmentid" -s > /dev/null &      
curl "http://10.10.198.204:3000/getParticipantMaster?participantcodeid=$participantcodeid&segmentid=$segmentid" -s > /dev/null &
curl "http://10.10.198.204:3000/getUserData?dealerid=$dealerid" -s > /dev/null&
wait
curl "http://10.10.198.204:3000/processOrderRequest?dealerid=$dealerid&tokenid=$tokenid&participantCodeId=$participantCodeId&segmentid=$segmentid" -s > /dev/null
done
end_time=$(date +%s%N)

# Calculate the time difference in microseconds
time_diff=$(( (end_time - start_time) / 1000 ))

# Print the time difference in microseconds
echo "Time difference: ${time_diff} microseconds"