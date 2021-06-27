cd test-network

rm -rf organizations/ordererOrganizations
rm -rf organizations/peerOrganizations

cd organizations/fabric-ca
mkdir temp
mkdir temp/ordererOrg
mkdir temp/org1
mkdir temp/org2
cp ordererOrg/*.yaml temp/ordererOrg/
cp org1/*.yaml temp/org1/
cp org2/*.yaml temp/org2/

rm -rf ordererOrg/*
rm -rf org1/*
rm -rf org2/*

cp temp/ordererOrg/*.yaml ordererOrg
cp temp/org1/*.yaml org1
cp temp/org2/*.yaml org2

rm -rf temp

cd ../..

./network.sh down
./network.sh up createChannel -ca -s couchdb
./network.sh deployCC -ccn goods-ledger -ccp ../goods-ledger/chaincode-go -ccl go

cd ../goods-ledger

rm -rf chaincode-go/vendor
rm -rf application-javascript/wallet
rm -rf application-javascript/node_modules

cd application-javascript

npm i
npm audit fix -f
npm start