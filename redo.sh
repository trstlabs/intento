sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ././trustless-hub-node_0.7.5-14-g9760714_amd64.deb -y

./inittest.sh
