sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-2-g147d243_amd64.deb -y

./init.sh
