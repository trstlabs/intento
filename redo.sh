kill -9 $(lsof -t -i:26657 -sTCP:LISTEN)
kill -9 $(lsof -t -i:1317 -sTCP:LISTEN)
kill -9 $(lsof -t -i:6060 -sTCP:LISTEN)
sudo apt remove trustlesshub -y
make clean
sudo apt install ./trustlesshub_0.7.0-48-ge9b2e77_amd64.deb -y

