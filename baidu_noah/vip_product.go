/* vip_product.go - get product info for vip */
/*
modification history
--------------------
2016/06/14, by liuxiaowei07, create
*/
/*
DESCRIPTION
*/

package baidu_noah

import (
	"encoding/json"
	"fmt"
	"strings"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

const (
	TIME_OUT = 30
)

/*
Format of message:
{
    "success":true,
	"message":"ok",
	"data":
	    [
	        {
			    "productId":"200001313",
				"id":"373365",
				"ip":"103.235.46.141",
				"product":"BAIDU_OP_BFE",
				"type":"bgw",
				"username":"libingyi,chenxiyang,luojiao01,zhangmiao02,zhangjiyang01,taochunhua,zhangweiwei09,yangsijie,liuzhuo,yaoguang",
				"manager":"OP-MANAGER:luojiao01,wugongwei"
			}
		]
}
*/

type ProductByVipItem struct {
	Id        *string // id of vip
	Ip        *string // ipaddr of vip
	Type      *string // type of vip
	Product   *string // product noah path of vip
	ProductId *string // product noah id of vip
	Username  *string // op-engineer for vip
	Manager   *string // op-manager for vip
}

type ProductByVip struct {
	Success bool               // true is success
	Message string             // ok when success, error msg when failed
	Data    []ProductByVipItem // [] when failed
}

// check product info format for result of api
func checkProductByVipItem(productInfoList []ProductByVipItem) error {
	for _, productInfo := range productInfoList {
		if productInfo.Id == nil || productInfo.Ip == nil || productInfo.Type == nil ||
			productInfo.Product == nil || productInfo.ProductId == nil || productInfo.Username == nil ||
			productInfo.Manager == nil {

			return fmt.Errorf("wrong format")
		}
	}

	return nil
}

func getProductByVips(urlStr string) ([]ProductByVipItem, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := ProductByVip{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// check result
	if !result.Success {
		return nil, fmt.Errorf("noah err: %s", result.Message)
	}
	if result.Data == nil {
		return nil, fmt.Errorf("noah err: no data return")
	}
	if err = checkProductByVipItem(result.Data); err != nil {
		return nil, fmt.Errorf("noah err: %s", err.Error())
	}

	return result.Data, nil
}

/* GetProductByVips - get product by the vip
 *
 * Params:
 *      - vips: vip list
 *
 * Returns:
 *      - (product info for all vips, err)
 */
func GetProductByVips(vips []string) ([]ProductByVipItem, error) {
	// generate url
	url := fmt.Sprintf("http://goat.noah.baidu.com/goat/index.php?r=Host/Assets/GetAssetsInfo&condition=%s&realType=0&searchType=assets&showFields=id,name,product,type&type=ip",
		strings.Join(vips, ","))

	// request api for result
	return getProductByVips(url)
}
