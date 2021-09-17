# Bootpay Go 플러그인

Bootpay Go 라이브러리는 Go 언어로 작성된 어플리케이션, 프레임워크 등에서 사용가능합니다.

## 설치하기 
 

```curl
go get github.com/bootpay/backend-go
``` 

## Getting Started

```go
package main

import (
	"fmt"
	"github.com/bootpay/backend-go"
)

func main() {
	bc := bootpay.Client{}.New("5b8f6a4d396fa665fdc2b5ea", "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw=", nil, "")

	fmt.Println("--------------- GetVerify() Start ---------------")
	token, err := bc.GetToken()
	fmt.Println("token : " + token.Data.Token)

	receiptId := "610c96352386840036db8bef"
  verify, err := bc.Verify(receiptId)
	if err != nil {
		fmt.Println("get token error: " + err.Error())
	}
	fmt.Println("--------------- GetVerify() End ---------------")
}
```

## Documentation

[부트페이 개발매뉴얼](https://app.gitbook.com/@bootpay)을 참조해주세요

## 기술문의

[부트페이 홈페이지](https://www.bootpay.co.kr) 우측 하단 채팅을 통해 기술문의 주세요!

## License

The gem is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
