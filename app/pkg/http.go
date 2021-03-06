package pkg

import (
	"encoding/json"

	"github.com/better-go/pkg/log"
	"github.com/go-resty/resty/v2"
)

/*


api docs:
	- https://docs.binance.org/api-reference/dex-api/paths.html#token
	- https://docs.binance.org/api-swagger/index.html#api-Tokens-getTokens
	- bsc chainID:
		- https://docs.binance.org/smart-chain/developer/rpc.html
		- mainnet: 56
		- testnet: 97
lib:
	- https://github.com/go-resty/resty

*/

const (
	// bsc:
	bscChainIdMainNet  = "56"
	bscChainIdTestNet  = "97"
	bscTokenMainNetUrl = "https://dex.binance.org/api/v1/tokens?limit=1000000000"
	bscTokenTestNetUrl = "https://testnet-dex.binance.org/api/v1/tokens?limit=1000000000"

	// eth:
	ethChainIdMainNet  = "1"
	ethTokenMainNetUrl = "https://raw.githubusercontent.com/MetaMask/contract-metadata/master/contract-map.json"
)

///
func HttpGet(url string) ([]byte, error) {
	// Create a Resty Client
	client := resty.New()

	// http get:
	resp, err := client.R().
		EnableTrace().
		Get(url)

	//log.Infof("resp: %+v, err: %v", resp, err)
	//log.Infof("resp: %+v, err: %v", resp.Body(), err)
	//log.Infof("resp: %+v, err: %v", resp.Body(), err)

	return resp.Body(), err

}

func GetEthTokenMeta() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// get:
	resp, err := HttpGet(ethTokenMainNetUrl)
	if err != nil {
		return nil, err
	}

	// fmt:
	err = json.Unmarshal(resp, &result)
	return result, err
}

func GetBscTokenMeta() (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// add mainNet:
	mainNet, err := getBscTokenMeta(bscTokenMainNetUrl)
	if err != nil {
		return nil, err
	}

	// add mainNet:
	data[bscChainIdMainNet] = mainNet

	// add testNet:
	testNet, err := getBscTokenMeta(bscTokenTestNetUrl)
	if err != nil {
		return nil, err
	}

	// add testNet:
	data[bscChainIdTestNet] = testNet
	return data, nil
}

func getBscTokenMeta(url string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	// slice:
	body := make([]map[string]interface{}, 0)
	//var body []map[string]interface{}

	// get:
	resp, err := HttpGet(url)
	if err != nil {
		return nil, err
	}

	// fmt:
	err = json.Unmarshal(resp, &body)
	//log.Infof("convet to body: %+v, err: %v", body, err)
	if err != nil {
		return nil, err
	}

	// convert:
	for _, item := range body {
		// assert:
		key, _ := item["original_symbol"].(string)

		// add:
		result[key] = item
	}

	log.Infof("result sum: %+v", len(body))
	return result, nil
}
