kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
sudo apt remove trustpriceprotocol -y
make clean
sudo apt install ./trustpriceprotocol_0.7.0-5-g9e2d314_amd64.deb -y
./init.sh
