/*-----------
あとで関数に分ける．
ref] json  https://qiita.com/nayuneko/items/2ec20ba69804e8bf7ca3
2019/05/29 FUKUSHIMA Kaito
--------*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	
	"net/http"
	"bytes"

	"log"
)

/** JSONデコード用の構造体 **/
type Value struct{
	Light	float32	`json:"d1"`
	Vib	float32	`json:"d2"`
	DT	string	`json:"created"`
}

func main() {
	//http----------------
	url :="http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1"
	res,err := http.Get(url)
	if err != nil{
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(body)
	html := buf.String()
	fmt.Println(html)

	//JSON---------------
	bytes := []byte(html)
	// //JSONデコード
	var values []Value
	if err := json.Unmarshal(bytes, &values); err != nil {
		log.Fatal(err)
	}
	//デコードデータの表示
	for _, d:=range values{
		fmt.Printf("%f : %f\n",d.Light,d.Vib)
	}
}