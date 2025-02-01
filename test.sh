echo "making 10 requests"
url="http://localhost:80"
for i in {1..10}; do
    echo "Sending request $i at $(date +"%T")"
    curl -s "$url" &
done
wait
echo "All requests finished"
