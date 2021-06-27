package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}


///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////         Added structures        ///////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

type Account struct {
	AccountToken               string `json:"AccountToken"`
	AccountType                string `json:"AccountType"`
	AccountName                string `json:"AccountName"`
	AccountUsername            string `json:"AccountUsername"`
	AccountEmail               string `json:"AccountEmail"`
	AccountPhoneNumber         string `json:"AccountPhoneNumber"`
	AccountPassword            string `json:"AccountPassword"`
	AccountOwnerManufacturerID string `json:"AccountOwnerManufacturerID"`
	DocType                    string `json:"DocType"`
}

type Product struct {
	ProductOwnerAccountID        string `json:"ProductOwnerAccountID"`
	ProductManufacturerID        string `json:"ProductManufacturerID"`
	ProductManufacturerName      string `json:"ProductManufacturerName"`
	ProductFactoryID             string `json:"ProductFactoryID"`
	ProductID                    string `json:"ProductID"`
	ProductName                  string `json:"ProductName"`
	ProductType                  string `json:"ProductType"`
	ProductBatch                 string `json:"ProductBatch"`
	ProductSerialinBatch         string `json:"ProductSerialinBatch"`
	ProductManufacturingLocation string `json:"ProductManufacturingLocation"`
	ProductManufacturingDate     string `json:"ProductManufacturingDate"`
	ProductExpiryDate            string `json:"ProductExpiryDate"`
	DocType                      string `json:"DocType"`
}

type Manufacturer struct {
	ManufacturerAccountID      string `json:"ManufacturerAccountID"`
	ManufacturerName           string `json:"ManufacturerName"`
	ManufacturerTradeLicenceID string `json:"ManufacturerTradeLicenceID"`
	ManufacturerLocation       string `json:"ManufacturerLocation"`
	ManufacturerFoundingDate   string `json:"ManufacturerFoundingDate"`
	DocType                    string `json:"DocType"`
}

type Factory struct {
	FactoryManufacturerID string `json:"FactoryManufacturerID"`
	FactoryID             string `json:"FactoryID"`
	FactoryName           string `json:"FactoryName"`
	FactoryLocation       string `json:"FactoryLocation"`
	DocType               string `json:"DocType"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////         Added functions        ///////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////

func (s *SmartContract) RegisterAccount(ctx contractapi.TransactionContextInterface,
	accountKey string, accountToken string, accountType string, accountName string, accountUsername string,
	accountEmail string, accountPassword string, accountOwnerManufacturerID string, docType string) error {

	account := Account {
		AccountToken:               accountToken,
		AccountType:                accountType,
		AccountName:                accountName,
		AccountUsername:            accountUsername,
		AccountEmail:               accountEmail,
		AccountPassword:            accountPassword,
		AccountOwnerManufacturerID: accountOwnerManufacturerID,
		DocType:                    docType,
	}

	accountAsBytes, err := json.Marshal(account)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountKey, accountAsBytes)
}

func (s *SmartContract) AddManufacturer(ctx contractapi.TransactionContextInterface,
	manufacturerAccountID string, manufacturerKey string, manufacturerName string, manufacturerTradeLicenceID string,
	manufacturerLocation string, manufacturerFoundingDate string, docType string) error {

	manufacturer := Manufacturer {
		ManufacturerAccountID:      manufacturerAccountID,
		ManufacturerName:           manufacturerName,
		ManufacturerTradeLicenceID: manufacturerTradeLicenceID,
		ManufacturerLocation:       manufacturerLocation,
		ManufacturerFoundingDate:   manufacturerFoundingDate,
		DocType:                    docType,
	}

	manufacturerAsBytes, err := json.Marshal(manufacturer)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(manufacturerKey, manufacturerAsBytes)
}

func (s *SmartContract) AddFactory(ctx contractapi.TransactionContextInterface,
	factoryKey string, factoryManufacturerID string, factoryID string, factoryName string, factoryLocation string,
	docType string) error {

	factory := Factory {
		FactoryManufacturerID: factoryManufacturerID,
		FactoryID:             factoryID,
		FactoryName:           factoryName,
		FactoryLocation:       factoryLocation,
		DocType:               docType,
	}

	factoryAsBytes, err := json.Marshal(factory)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(factoryKey, factoryAsBytes)
}

func (s *SmartContract) AddProduct(ctx contractapi.TransactionContextInterface,
	productKey string, productOwnerAccountID string, productManufacturerID string, productManufacturerName string, productFactoryID string,
	productID string, productName string, productType string, productBatch string, productSerialinBatch string,
	productManufacturingLocation string, productManufacturingDate string, productExpiryDate string, docType string) error {

	product := Product {
		ProductOwnerAccountID:        productOwnerAccountID,
		ProductManufacturerID:        productManufacturerID,
		ProductManufacturerName:      productManufacturerName,
		ProductFactoryID:             productFactoryID,
		ProductID:                    productID,
		ProductName:                  productName,
		ProductType:                  productType,
		ProductBatch:                 productBatch,
		ProductSerialinBatch:         productSerialinBatch,
		ProductManufacturingLocation: productManufacturingLocation,
		ProductManufacturingDate:     productManufacturingDate,
		ProductExpiryDate:            productExpiryDate,
		DocType:                      docType,
	}

	productAsBytes, err := json.Marshal(product)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(productKey, productAsBytes)
}

func (s *SmartContract) UpdateProductOwner(ctx contractapi.TransactionContextInterface,
	productKey string, productOwnerAccountID string) error {

	productAsBytes, err := ctx.GetStub().GetState(productKey)

	if err != nil {
		return err
	}

	product := Product{}

	json.Unmarshal(productAsBytes, &product)
	product.ProductOwnerAccountID = productOwnerAccountID

	productAsBytes, err = json.Marshal(product)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(productKey, productAsBytes)
}

func (s *SmartContract) UpdateAccountOwnerManufacturerID(ctx contractapi.TransactionContextInterface,
	accountKey string, accountOwnerManufacturerID string) error {

	accountAsBytes, err := ctx.GetStub().GetState(accountKey)

	if err != nil {
		return err
	}

	account := Account{}

	json.Unmarshal(accountAsBytes, &account)
	account.AccountOwnerManufacturerID = accountOwnerManufacturerID

	accountAsBytes, err = json.Marshal(account)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountKey, accountAsBytes)
}

func (s *SmartContract) UpdateAccountToken(ctx contractapi.TransactionContextInterface,
	accountKey string, accountToken string) error {

	accountAsBytes, err := ctx.GetStub().GetState(accountKey)
	
	if err != nil {
		return err
	}
	
	account := Account{}

	json.Unmarshal(accountAsBytes, &account)
	account.AccountToken = accountToken

	accountAsBytes, err = json.Marshal(account)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountKey, accountAsBytes)
}

func (s *SmartContract) UpdateAccount(ctx contractapi.TransactionContextInterface,
	accountKey string, accountToken string, accountName string, accountEmail string, accountPhoneNumber string) error {

	accountAsBytes, err := ctx.GetStub().GetState(accountKey)

	if err != nil {
		return err
	}

	account := Account{}

	json.Unmarshal(accountAsBytes, &account)
	account.AccountToken = accountToken
	account.AccountName = accountName
	account.AccountEmail = accountEmail
	account.AccountPhoneNumber = accountPhoneNumber

	accountAsBytes, err = json.Marshal(account)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(accountKey, accountAsBytes)
}

func (s *SmartContract) UpdateManufacturer(ctx contractapi.TransactionContextInterface,
	manufacturerKey string, manufacturerName string, manufacturerTradeLicenceID string, manufacturerLocation string,
	manufacturerFoundingDate string) error {

	manufacturerAsBytes, err := ctx.GetStub().GetState(manufacturerKey)

	if err != nil {
		return err
	}

	manufacturer := Manufacturer{}

	json.Unmarshal(manufacturerAsBytes, &manufacturer)
	manufacturer.ManufacturerName = manufacturerName
	manufacturer.ManufacturerTradeLicenceID = manufacturerTradeLicenceID
	manufacturer.ManufacturerLocation = manufacturerLocation
	manufacturer.ManufacturerFoundingDate = manufacturerFoundingDate

	manufacturerAsBytes, err = json.Marshal(manufacturer)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(manufacturerKey, manufacturerAsBytes)
}

func (s *SmartContract) UpdateFactory(ctx contractapi.TransactionContextInterface,
	factoryKey string, factoryManufacturerID string, factoryName string, factoryLocation string) error {

	factoryAsBytes, err := ctx.GetStub().GetState(factoryKey)

	if err != nil {
		return err
	}

	factory := Factory{}

	json.Unmarshal(factoryAsBytes, &factory)
	factory.FactoryManufacturerID = factoryManufacturerID
	factory.FactoryName = factoryName
	factory.FactoryLocation = factoryLocation

	factoryAsBytes, err = json.Marshal(factory)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(factoryKey, factoryAsBytes)
}

func (s *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface,
	productKey string, productOwnerAccountID string, productFactoryID string, productName string, productType string, productBatch string,
	productSerialinBatch string, productManufacturingLocation string, productManufacturingDate string, productExpiryDate string) error {

	productAsBytes, err := ctx.GetStub().GetState(productKey)

	if err != nil {
		return err
	}

	product := Product{}

	json.Unmarshal(productAsBytes, &product)
	product.ProductOwnerAccountID = productOwnerAccountID
	product.ProductFactoryID = productFactoryID
	product.ProductName = productName
	product.ProductType = productType
	product.ProductBatch = productBatch
	product.ProductSerialinBatch = productSerialinBatch
	product.ProductManufacturingLocation = productManufacturingLocation
	product.ProductManufacturingDate = productManufacturingDate
	product.ProductExpiryDate = productExpiryDate

	productAsBytes, err = json.Marshal(product)

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(productKey, productAsBytes)
}

func (s *SmartContract) QueryAccountbyToken(ctx contractapi.TransactionContextInterface,
	accountToken string) ([]*Account, error) {
	
	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"account",
				"AccountToken":"%s"
			}
		}`,
		accountToken,
	)
	
	return getAccountQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryAccountbyEmail(ctx contractapi.TransactionContextInterface,
	accountEmail string) ([]*Account, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"account",
				"AccountEmail":"%s"
			}
		}`,
		accountEmail,
	)
	
	return getAccountQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryAccountbyUsername(ctx contractapi.TransactionContextInterface,
	accountUsername string) ([]*Account, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"account",
				"AccountUsername":"%s"
			}
		}`,
		accountUsername,
	)
	
	return getAccountQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryManufacturerbyAccountID(ctx contractapi.TransactionContextInterface,
	manufacturerAccountID string) ([]*Manufacturer, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"manufacturer",
				"ManufacturerAccountID":"%s"
			}
		}`,
		manufacturerAccountID,
	)

	return getManufacturerQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryManufacturerbyTradeLicenceID(ctx contractapi.TransactionContextInterface,
	manufacturerTradeLicenceID string) ([]*Manufacturer, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"manufacturer",
				"ManufacturerTradeLicenceID":"%s"
			}
		}`,
		manufacturerTradeLicenceID,
	)

	return getManufacturerQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryFactorybyID(ctx contractapi.TransactionContextInterface,
	factoryID string) ([]*Factory, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"factory",
				"FactoryID":"%s"
			}
		}`,
		factoryID,
	)

	return getFactoryQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryFactorybyManufacturerID(ctx contractapi.TransactionContextInterface,
	factoryManufacturerID string) ([]*Factory, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"factory",
				"FactoryManufacturerID":"%s"
			}
		}`,
		factoryManufacturerID,
	)

	return getFactoryQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryProductbyID(ctx contractapi.TransactionContextInterface,
	productID string) ([]*Product, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"product",
				"ProductID":"%s"
			}
		}`,
		productID,
	)

	return getProductQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryProductbyCode(ctx contractapi.TransactionContextInterface,
	productCode string) ([]*Product, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"product",
				"_id":"%s"
			}
		}`,
		productCode,
	)

	return getProductQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryProductbyOwnerAccountID(ctx contractapi.TransactionContextInterface,
	productOwnerAccountID string) ([]*Product, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"product",
				"ProductOwnerAccountID":"%s"
			}
		}`,
		productOwnerAccountID,
	)

	return getProductQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryProductbyManufacturerID(ctx contractapi.TransactionContextInterface,
	productManufacturerID string) ([]*Product, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"product",
				"ProductManufacturerID":"%s"
			}
		}`,
		productManufacturerID,
	)

	return getProductQueryResultForQueryString(ctx, queryString)
}

func (s *SmartContract) QueryProductbyFactoryID(ctx contractapi.TransactionContextInterface,
	productFactoryID string) ([]*Product, error) {

	var queryString = fmt.Sprintf(
		`{
			"selector":{
				"DocType":"product",
				"ProductFactoryID":"%s"
			}
		}`,
		productFactoryID,
	)

	return getProductQueryResultForQueryString(ctx, queryString)
}

func getAccountQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Account, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructAccountQueryResponseFromIterator(resultsIterator)
}

func getProductQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructProductQueryResponseFromIterator(resultsIterator)
}

func getManufacturerQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Manufacturer, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructManufacturerQueryResponseFromIterator(resultsIterator)
}

func getFactoryQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Factory, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructFactoryQueryResponseFromIterator(resultsIterator)
}

func constructAccountQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Account, error) {
	var accounts []*Account
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var account Account
		err = json.Unmarshal(queryResult.Value, &account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func constructProductQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Product, error) {
	var products []*Product
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var product Product
		err = json.Unmarshal(queryResult.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func constructManufacturerQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Manufacturer, error) {
	var manufacturers []*Manufacturer
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var manufacturer Manufacturer
		err = json.Unmarshal(queryResult.Value, &manufacturer)
		if err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, &manufacturer)
	}

	return manufacturers, nil
}

func constructFactoryQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Factory, error) {
	var factories []*Factory
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var factory Factory
		err = json.Unmarshal(queryResult.Value, &factory)
		if err != nil {
			return nil, err
		}
		factories = append(factories, &factory)
	}

	return factories, nil
}