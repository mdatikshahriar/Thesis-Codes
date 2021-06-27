# GoodsLedger-Server
A server based on Hyperledger Fabric which will help GoodsLedger app to detect counterfeit product.


### Prerequisite

**OS**: Any Linux distro

**Languages**: Javascript, Go

**Other**: cURL, Docker, Docker-compose, node, npm


### How to run the server

Follow below steps to run the server




#### Development environment setup

1. Clone the repository in a local computer running any Linux distro.

2. Then setup the the computer for Hyperledger Fabric. Follow these steps:

>i. Install cURL - a command line tool to access web protocols in command line:
  
```
sudo apt-get install curl
```

>ii. Next, install Docker and Docker Compose using the following steps.
    
>>a. Remove the older versions of Docker (if any) using the following commands:
    
```
sudo apt-get remove docker docker-engine docker.io containerd run
```

>>It is okay if the above command reports that none of these packages is installed
    
>>b. Update apt: 
      
```
sudo apt-get update
```
      
>>c. Install packages to allow apt to use a repository over HTTPS:
   
```
sudo apt-get install apt-transport-https ca-certificates gnupg-agent software-properties-common
```

>>d. Add Docker’s official GPG key:
    
```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

>>e. Use the following command to set up the stable repository.
      
```
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

>>f. Update the apt package again:
      
```
sudo apt-get update
```

>>g. Install docker:
      
```
sudo apt-get install docker-ce docker-ce-cli containerd.io
```
      
>>h. Test docker installation:
    
```
docker run hello-world
```

>>If it shows a message like "Hello from Docker" -> Then docker is fully installed
       
>>i. If you cannot run the above command and it shows the permission denied message, then you might need to use the following commands to ensure that the user can run the docker commands without being sudo:
       
>>>i.
        
```
sudo groupadd docker
```

>>>(it might report that the docker group is already added, it is fine)
         
>>>ii.
         
```
sudo usermod -aG docker $USER / sudo gpasswd -a $USER docker
```

>>>(adding current user to the docker group, any command will do)
         
>>>iii. 
         
```
newgrp docker
```
        
>>>(to activate the change)
         
>>>iv.
         
```
docker run hello-world
```

>>>If it shows a message like "Hello from Docker" -> Then docker is fully installed


3. Next, install docker-compose using the following steps. Docker Compose makes it easier for users to streamline the docker containers management processes, including starting up, shutting down, and setting up intra-container linking and volumes.

>a. Download the docker-componse binary file from this location:

```
sudo curl -L https://github.com/docker/compose/releases/download/1.21.2/dockercompose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
```
    
>>b. Make it executable:
```
sudo chmod +x /usr/local/bin/docker-compose
```

>>c. Check the installation:

```
docker-compose –version
```

>>If it shows an output with a version, the docker-compose is properly installed.

4. Next, install go programming language using the following steps:

>a. Go to: `https://golang.org/dl/`. Then click for the Linux tar file. Once clicked, a popup will appear, select Save.
  
>b. The file will be downloaded in ~/Downloads. cd into ~/Downloads in the command prompt and then issue the following command:
  
```
sudo tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz
```

>c. Examine if the following structure has been created: /usr/local/go. If it does not match it like this, see what is directory structure it has created. You need locate the directory called "go" under which api, bin and other directories need to reside. Note   the directory path containing the "go" directory.

>d. Add the "go" directory in the PATH variable. To do this, use the following commands:
 
>>i.

```
cd
```

>>ii.
```
nano .profile
```
    
>>iii. Add this line at the end of your profile:

```
export PATH=$PATH:/PATH_TO_GO/bin
```

>>*(Remember to modify the value of PATH_TO_GO accordingly.)*

>>iv. Save by using: `ctrl+o` and then `ctrl+x`

>>v.

```
source .profile
```

>>vi. Check if the proper go path is printed in the console:

```
echo $PATH
```

>e. Add GOPATH variable into your path:
      
>>i.
  
```
cd
```

>>ii.
 
```
mkdir -p go/src
```

>>iii.

```
mkdir -p go/bin
```

>>iv.

```
cd
```

>>v.

```
nano .bashrc
```

>>vi. Add:

```
export GOPATH=$HOME/go
```
        
>>vii. Add:

```
export PATH=$PATH:$GOPATH/bin
```

>>viii. `ctrl+o` and then `ctrl+x` to save and exit

>>ix.

```
source .bashrc
```

>>x. Check if the proper GOPATH is printed in the console:
  
```
echo $GOPATH
```

>>xi. Check if bin under GOPATH is in the path:

```
echo $PATH
```


5. Fabric uses a specific version of Node. Next, install this specific version using the following commands:
  
>a.
  
```
sudo npm cache clean -f
```

>b.
 
```
sudo npm install -g n
```

>c.
  
```
sudo n 8.9
```

>d. Check node version:
 
```
node --version
```


6. Now add the location of Fabric installation (~/GoodsLedger-Server/fabric-samples) into the path in `.bashrc`:

>a.
  
```
cd
```

>b.
    
```
nano .bashrc
```

>c. Add this line at the end of `.bashrc` file:

```
export PATH=$PATH:/PATH_TO_SERVER/GoodsLedger-Server/fabric-samples/bin
```

>*(Remember to modify the value of PATH_TO_SERVER accordingly.)*
  
>d. `crtl+o` and then `crtl+x` to save and exit
  
>e.

```
source .bashrc
```
  
>f. Check if the fabric is in your path:
  
```
echo $PATH
```


    
#### Chaincode development, deployment & interaction

1. At first, cd into the fabric directory:

```
cd GoodsLedger-Server/fabric-samples/fabcar
```
  
2. While in the fabcar directory, issue the following command:

```
./startFabric.sh go
```

This will create the required Fabric network. If there is no error, then the network is created successfully.

3. Next, issue the following command:

```
cd javascript
```

4. Then create a new file inside `/javascript` folder named `.env`. Inside the file type the following and save it:

```
TOKEN_SECRET = YOUR_PREFERRED_TOKEN
```

*(Remember to modify the value of YOUR_PREFERRED_TOKEN accordingly.)*
  
5. Next, install the npm modules:

```
npm install
```
  
6. The following two commands will interact with the CA. To understand and observe what is happening it is often useful to see the logs generated by CA docker container. Use this command in a separate terminal:

```
docker logs -f ca.example.com
```
  
7. Next, we need to enroll an Admin for the network who can issue identities to different users. For this, use the following command:

```
node enrollAdmin.js
```

Check the logs in the CA container.
   
8. To utilize the chaincode we users who are registered with the CA, which will be carried out by the previously enrolled admin. Issue this command:

```
node registerUser.js
```

Again, check the logs in the CA container.
   
9. Finally run the server by issuing following command:

```
npm start
```


###### If you followed everything above, hopefully the server will run successfully in your local computer. If you want to see a simple web interface of the server just go to `localhost:3000`.
