kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
sudo apt remove trustlesshub -y
make clean
sudo apt install ./trustlesshub_0.7.0-50-g75a3e70_amd64.deb -y

