sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.7-87-gf2e5658_amd64.deb -y

./inittest.sh