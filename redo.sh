sudo apt remove trustless-hub-node -y
make clean-files
sudo apt install ./trustless-hub-node_0.7.5-15-g669de44_amd64.deb -y

./inittest.sh
