sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.5-21-g566c800_amd64.deb -y

./init.sh
