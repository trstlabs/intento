sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.5-13-g3d93a90_amd64.deb -y

./inittest.sh
