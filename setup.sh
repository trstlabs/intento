
//reinstalling
make deb
sudo apt remove trustpriceprotocol -y
sudo apt install ./trustpriceprotocol_0.6.1.0-14-gd1f4289_amd64.deb
./init.sh

tppd init-bootstrap
PUBLIC_KEY=$(tppd parse attestation_cert.der 2> /dev/null | cut -c 3-)
echo $PUBLIC_KEY