kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
kill -9 $(lsof -t -i:9090 -sTCP:LISTEN)
kill -9 $(lsof -t -i:54246  -sTCP:LISTEN)
kill -9 $(lsof -t -i:26656 -sTCP:LISTEN)
sudo apt remove trustlesshub -y
make clean-files
sudo apt install ./trustlesshub_0.7.0-74-g1ee95b8_amd64.deb -y

./inittest.sh
