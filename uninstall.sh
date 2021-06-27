cd fabric-samples

./stop.sh

cd ..

mkdir .temp
attrib +h .temp
cp -r fabric-samples/goods-ledger .temp
cp -r fabric-samples/4host-swarm .temp
cp fabric-samples/run.sh .temp
cp fabric-samples/stop.sh .temp
cp fabric-samples/.gitignore .temp
rm -rf fabric-samples

mkdir fabric-samples

cp -r .temp/goods-ledger fabric-samples
cp -r .temp/4host-swarm fabric-samples
cp .temp/run.sh fabric-samples
cp .temp/stop.sh fabric-samples
cp .temp/.gitignore fabric-samples
rm -rf .temp