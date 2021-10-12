kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
sudo apt remove trustpriceprotocol -y
make clean
sudo apt install ./trustpriceprotocol_0.7.0-4-g35aa614_amd64.deb -y
./init.sh
