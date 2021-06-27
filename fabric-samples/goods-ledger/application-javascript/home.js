/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('../../test-application/javascript/CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('../../test-application/javascript/AppUtil.js');

const channelName = 'mychannel';
const chaincodeName = 'goods-ledger';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';
let contract = null;

// pre-requisites:
// - fabric-sample two organization test-network setup with two peers, ordering service,
//   and 2 certificate authorities
//         ===> from directory /fabric-samples/test-network
//         ./network.sh up createChannel -ca
// - Use any of the asset-transfer-basic chaincodes deployed on the channel "mychannel"
//   with the chaincode name of "basic". The following deploy command will package,
//   install, approve, and commit the javascript chaincode, all the actions it takes
//   to deploy a chaincode to a channel.
//         ===> from directory /fabric-samples/test-network
//         ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-javascript/ -ccl javascript
// - Be sure that node.js is installed
//         ===> from directory /fabric-samples/asset-transfer-basic/application-javascript
//         node -v
// - npm installed code dependencies
//         ===> from directory /fabric-samples/asset-transfer-basic/application-javascript
//         npm install
// - to run this test application
//         ===> from directory /fabric-samples/asset-transfer-basic/application-javascript
//         node app.js

// NOTE: If you see  kind an error like these:
/*
    2020-08-07T20:23:17.590Z - error: [DiscoveryService]: send[mychannel] - Channel:mychannel received discovery error:access denied
    ******** FAILED to run the application: Error: DiscoveryService: mychannel error: access denied

   OR

   Failed to register user : Error: fabric-ca request register failed with errors [[ { code: 20, message: 'Authentication failure' } ]]
   ******** FAILED to run the application: Error: Identity not found in wallet: appUser
*/
// Delete the /fabric-samples/asset-transfer-basic/application-javascript/wallet directory
// and retry this application.
//
// The certificate authority must have been restarted and the saved certificates for the
// admin and application user are not valid. Deleting the wallet store will force these to be reset
// with the new certificate authority.
//

/**
 *  A test application to show basic queries operations with any of the asset-transfer-basic chaincodes
 *   -- How to submit a transaction
 *   -- How to query and check the results
 *
 * To see the SDK workings, try setting the logging to show on the console before running
 *        export HFC_LOGGING='{"debug":"console"}'
 */
async function main() {
	try {
		// build an in memory object with the network configuration (also known as a connection profile)
		const ccp = buildCCPOrg1();

		// build an instance of the fabric ca services client based on
		// the information in the network configuration
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');

		// setup the wallet to hold the credentials of the application user
		const wallet = await buildWallet(Wallets, walletPath);

		// in a real application this would be done on an administrative flow, and only once
		await enrollAdmin(caClient, wallet, mspOrg1);

		// in a real application this would be done only when a new user was required to be added
		// and would be part of an administrative flow
		await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');

		// Create a new gateway instance for interacting with the fabric network.
		// In a real application this would be done as the backend server session is setup for
		// a user that has been verified.
		const gateway = new Gateway();

		try {
			// setup the gateway instance
			// The user will now be able to create connections to the fabric network and be able to
			// submit transactions and query. All transactions submitted by this gateway will be
			// signed by this user using the credentials stored in the wallet.
			await gateway.connect(ccp, {
				wallet,
				identity: org1UserId,
				discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
			});

			// Build a network instance based on the channel where the smart contract is deployed
			const network = await gateway.getNetwork(channelName);

			// Get the contract from the network.
			contract = network.getContract(chaincodeName);

		}  catch (error) {
			console.error(`******** FAILED to run the application: ${error}`);
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
	}
}

main();

//.......................................

const express = require('express');
const router = express.Router();
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');

//Home
router.use(express.static(path.join(__dirname, 'views')));

router.get('/', (req, res) => {
    res.render(path.join(__dirname, 'views', 'home.ejs'));
});

router.post('/loginAccount', async (req, res) => {
    const accountUsername = String(req.body.accountUsername);

    const usernameResult = await contract.evaluateTransaction('QueryAccountbyUsername', accountUsername);
    let usernameResultObject = usernameResult.toJSON().data;

    if(usernameResultObject[0] == null){
        return res.send('Username doesn\'t exist');
    }

    usernameResultObject = JSON.parse(usernameResult)[0];
    console.log(usernameResultObject);

    const validPassword = await bcrypt.compare(String(req.body.accountPassword), usernameResultObject.AccountPassword);
    if(!validPassword){
        return res.send('Invalid password.');
    }

    const secretKey = process.env.TOKEN_SECRET || "TOKEN_SECRET";
    const accountToken = jwt.sign({_id: usernameResultObject.AccountToken}, secretKey);

    await contract.submitTransaction('UpdateAccountToken', usernameResultObject.AccountToken, accountToken);

    res.send(JSON.stringify(usernameResultObject));
});

router.post('/registerAccount', async (req, res) => {

    const accountType = String(req.body.accountType);
    const accountName = String(req.body.accountName);
    const accountUsername = String(req.body.accountUsername);
    const accountEmail = String(req.body.accountEmail);
    const accountOwnerManufacturerID = String(req.body.accountOwnerManufacturerID);

    const salt = await bcrypt.genSalt(10);
    let hashedAccountPassword = await bcrypt.hash(req.body.accountPassword, salt);
    const hashedConfirmedAccountPassword = await bcrypt.hash(req.body.accountConfirmedPassword, salt);

    const usernameResult = await contract.evaluateTransaction('QueryAccountbyUsername', accountUsername);
    const usernameResultObject = usernameResult.toJSON().data;

    if(usernameResultObject[0] != null){
         return res.send('Username already exist.');
    }

    const emailResult = await contract.evaluateTransaction('QueryAccountbyEmail', accountEmail);
    const emailResultObject = emailResult.toJSON().data;

    if(emailResultObject[0] != null){
         return res.send('Email already exist.');
    }

    if(hashedAccountPassword !== hashedConfirmedAccountPassword){
        return res.send('Password didn\'t match.');
    }

    const accountKeyValue = accountType + accountEmail + accountUsername;
    const newSalt = await bcrypt.genSalt(10);
    let accountKey = await bcrypt.hash(accountKeyValue, newSalt);
    const docType = "account";

    const secretKey = process.env.TOKEN_SECRET || "TOKEN_SECRET";
    let accountToken = jwt.sign({_id: accountKey}, secretKey);

    accountKey = String(accountKey);
    accountToken = String(accountToken);
    hashedAccountPassword = String(hashedAccountPassword);

    await contract.submitTransaction('RegisterAccount', accountKey, accountToken, accountType, accountName, accountUsername, accountEmail, hashedAccountPassword, accountOwnerManufacturerID, docType);

    res.send(JSON.stringify({ accountKey, accountToken, accountType, accountName, accountUsername, accountEmail, hashedAccountPassword, accountOwnerManufacturerID, docType }));
});

router.post('/addManufacturer', async (req, res) => {
    const manufacturerAccountID = String(req.body.manufacturerAccountID);
    const manufacturerName = String(req.body.manufacturerName);    
    const manufacturerTradeLicenceID = String(req.body.manufacturerTradeLicenceID);
    const manufacturerLocation = String(req.body.manufacturerLocation);
    const manufacturerFoundingDate = String(req.body.manufacturerFoundingDate);
    const docType = "manufacturer";

    const tradeLicenceIDResult = await contract.evaluateTransaction('QueryManufacturerbyTradeLicenceID', manufacturerTradeLicenceID);
    const tradeLicenceIDResultObject = tradeLicenceIDResult.toJSON().data;

    if(tradeLicenceIDResultObject[0] != null){
         return res.send('This Trade Licence belongs to someone else.');
    }

    const manufacturerKeyValue = manufacturerAccountID + manufacturerName + manufacturerTradeLicenceID;
    const salt = await bcrypt.genSalt(10);
    const manufacturerKey = await bcrypt.hash(manufacturerKeyValue, salt);

    const accountKey = manufacturerAccountID;
    const accountOwnerManufacturerID = manufacturerKey;

    await contract.submitTransaction('UpdateAccountOwnerManufacturerID', accountKey, accountOwnerManufacturerID);

    await contract.submitTransaction('AddManufacturer', manufacturerKey, manufacturerAccountID, manufacturerName, manufacturerTradeLicenceID, manufacturerLocation, manufacturerFoundingDate, docType);

    res.send(JSON.stringify({ manufacturerKey, manufacturerAccountID, manufacturerName, manufacturerTradeLicenceID, manufacturerLocation, manufacturerFoundingDate, docType }));
});

router.post('/addFactory', async (req, res) => {
    const factoryManufacturerID = String(req.body.factoryManufacturerID);
    const factoryID = String(req.body.factoryID);
    const factoryName = String(req.body.factoryName);
    const factoryLocation = String(req.body.factoryLocation);
    const docType = "factory";

    const factoryKeyValue = factoryManufacturerID + factoryID + factoryName;
    const salt = await bcrypt.genSalt(10);
    const factoryKey = await bcrypt.hash(factoryKeyValue, salt);

    await contract.submitTransaction('AddFactory', factoryKey, factoryManufacturerID, factoryID, factoryName, factoryLocation, docType);

    res.send(JSON.stringify({ factoryKey, factoryManufacturerID, factoryID, factoryName, factoryLocation, docType }));
});

router.post('/addProduct', async (req, res) => {
    const productOwnerAccountID = String(req.body.productOwnerAccountID);
    const productManufacturerID = String(req.body.productManufacturerID);
    const productManufacturerName = String(req.body.productManufacturerName);
    const productFactoryID = String(req.body.productFactoryID);
    const productID = String(req.body.productID);
    const productName = String(req.body.productName);
    const productType = String(req.body.productType);
    const productBatch = String(req.body.productBatch);
    const productSerialinBatch = String(req.body.productSerialinBatch);
    const productManufacturingLocation = String(req.body.productManufacturingLocation);
    const productManufacturingDate = String(req.body.productManufacturingDate);
    const productExpiryDate = String(req.body.productExpiryDate);
    const docType = "product"

    const productKeyValue = productManufacturerID + productFactoryID + productBatch + productID + productSerialinBatch;
    const salt = await bcrypt.genSalt(10);
    const productKey = await bcrypt.hash(productKeyValue, salt);

    await contract.submitTransaction('AddProduct', productKey, productOwnerAccountID, productManufacturerID, productManufacturerName, productFactoryID, productID, productName, productType, productBatch, productSerialinBatch, productManufacturingLocation, productManufacturingDate, productExpiryDate, docType);

    res.send(JSON.stringify({ productKey, productOwnerAccountID, productManufacturerID, productManufacturerName, productFactoryID, productID, productName, productType, productBatch, productSerialinBatch, productManufacturingLocation, productManufacturingDate, productExpiryDate, docType }));
});

router.post('/updateProductOwner', async (req, res) => {
    const productOwnerAccountID = String(req.body.productOwnerAccountID);
    const productKey = String(req.body.productKey);

    await contract.submitTransaction('UpdateProductOwner', productKey, productOwnerAccountID);

    res.send(JSON.stringify({ productKey, productOwnerAccountID }));
});

router.post('/updateAccountToken', async (req, res) => {
    const accountKey = String(req.body.accountKey);
    const accountToken = String(req.body.accountToken);

    await contract.submitTransaction('UpdateAccountToken', accountKey, accountToken);

    res.send(JSON.stringify({ accountKey, accountToken }));
});

router.post('/updateAccount', async (req, res) => {
    const accountKey = String(req.body.accountKey);
    const accountToken = String(req.body.accountToken);
    const accountName = String(req.body.accountName);
    const accountEmail = String(req.body.accountEmail);
    const accountPhoneNumber = String(req.body.accountPhoneNumber);

    const emailResult = await contract.evaluateTransaction('QueryAccountbyEmail', accountEmail);
    const emailResultObject = emailResult.toJSON().data;

    if(emailResultObject[0] != null){
         return res.send('Email already exist.');
    }

    await contract.submitTransaction('UpdateAccount', accountKey, accountToken, accountName, accountEmail, accountPhoneNumber);

    res.send(JSON.stringify({ accountKey, accountToken, accountName, accountEmail, accountPhoneNumber }));
});

router.post('/updateManufacturer', async (req, res) => {
    const manufacturerKey = String(req.body.manufacturerKey);
    const manufacturerName = String(req.body.manufacturerName);
    const manufacturerTradeLicenceID = String(req.body.manufacturerTradeLicenceID);
    const manufacturerLocation = String(req.body.manufacturerLocation);
    const manufacturerFoundingDate = String(req.body.manufacturerFoundingDate);

    await contract.submitTransaction('UpdateManufacturer', manufacturerKey, manufacturerName, manufacturerTradeLicenceID, manufacturerLocation, manufacturerFoundingDate);

    res.send(JSON.stringify({ manufacturerKey, manufacturerName, manufacturerTradeLicenceID, manufacturerLocation, manufacturerFoundingDate }));
});

router.post('/updateFactory', async (req, res) => {
    const factoryKey = String(req.body.factoryKey);
    const factoryManufacturerID = String(req.body.factoryManufacturerID);
    const factoryName = String(req.body.factoryName);
    const factoryLocation = String(req.body.factoryLocation);

    await contract.submitTransaction('UpdateFactory', factoryKey, factoryManufacturerID, factoryName, factoryLocation);

    res.send(JSON.stringify({ factoryKey, factoryManufacturerID, factoryName, factoryLocation }));
});

router.post('/updateProduct', async (req, res) => {
    const productKey = String(req.body.productKey);
    const productOwnerAccountID = String(req.body.productOwnerAccountID);
    const productFactoryID = String(req.body.productFactoryID);
    const productName = String(req.body.productName);
    const productType = String(req.body.productType);
    const productBatch = String(req.body.productBatch);
    const productSerialinBatch = String(req.body.productSerialinBatch);
    const productManufacturingLocation = String(req.body.productManufacturingLocation);
    const productManufacturingDate = String(req.body.productManufacturingDate);
    const productExpiryDate = String(req.body.productExpiryDate);

    await contract.submitTransaction('UpdateProduct', productKey, productOwnerAccountID, productFactoryID, productName, productType, productBatch, productSerialinBatch, productManufacturingLocation, productManufacturingDate, productExpiryDate);

    res.send(JSON.stringify({ productKey, productOwnerAccountID, productFactoryID, productName, productType, productBatch, productSerialinBatch, productManufacturingLocation, productManufacturingDate, productExpiryDate }));
});

router.post('/queryAccountbyToken', async (req, res) => {
    const accountToken = String(req.body.accountToken);

    const result = await contract.evaluateTransaction('QueryAccountbyToken', accountToken);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject[0]));
});

router.post('/queryAccountbyEmail', async (req, res) => {
    const accountEmail = String(req.body.accountEmail);

    const result = await contract.evaluateTransaction('QueryAccountbyEmail', accountEmail);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject[0]));
});

router.post('/queryAccountbyUsername', async (req, res) => {
    const accountUsername = String(req.body.accountUsername);

    const result = await contract.evaluateTransaction('QueryAccountbyUsername', accountUsername);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject[0]));
});

router.post('/queryManufacturerbyAccountID', async (req, res) => {
    const manufacturerAccountID = String(req.body.manufacturerAccountID);

    const result = await contract.evaluateTransaction('QueryManufacturerbyAccountID', manufacturerAccountID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject[0]));
});

router.post('/queryManufacturerbyTradeLicenceID', async (req, res) => {
    const manufacturerTradeLicenceID = String(req.body.manufacturerTradeLicenceID);

    const result = await contract.evaluateTransaction('QueryManufacturerbyTradeLicenceID', manufacturerTradeLicenceID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject[0]));
});

router.post('/queryFactorybyManufacturerID', async (req, res) => {
    const factoryManufacturerID = String(req.body.factoryManufacturerID);

    const result = await contract.evaluateTransaction('QueryFactorybyManufacturerID', factoryManufacturerID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryFactorybyID', async (req, res) => {
    const factoryID = String(req.body.factoryID);

    const result = await contract.evaluateTransaction('QueryFactorybyID', factoryID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryProductbyID', async (req, res) => {
    const productID = String(req.body.productID);

    const result = await contract.evaluateTransaction('QueryProductbyID', productID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryProductbyCode', async (req, res) => {
    const productCode = String(req.body.productCode);

    const result = await contract.evaluateTransaction('QueryProductbyCode', productCode);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryProductbyOwnerAccountID', async (req, res) => {
    const productOwnerAccountID = String(req.body.productOwnerAccountID);

    const result = await contract.evaluateTransaction('QueryProductbyOwnerAccountID', productOwnerAccountID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryProductbyManufacturerID', async (req, res) => {
    const productManufacturerID = String(req.body.productManufacturerID);

    const result = await contract.evaluateTransaction('QueryProductbyManufacturerID', productManufacturerID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

router.post('/queryProductbyFactoryID', async (req, res) => {
    const productFactoryID = String(req.body.productFactoryID);

    const result = await contract.evaluateTransaction('QueryProductbyFactoryID', productFactoryID);
    const resultObject = JSON.parse(result);

    res.send(JSON.stringify(resultObject));
});

module.exports = router;