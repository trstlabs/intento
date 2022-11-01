kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
sudo apt remove trst -y
make clean
sudo apt install ./trst_0.7.7-106-g66551b04_amd64.deb -y

