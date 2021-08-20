module github.com/libs4go/ethers

go 1.16

require (
	github.com/libs4go/cfd4go v0.0.1
	github.com/libs4go/crypto v0.0.0-20210720063913-ba28ec544c0f
	github.com/libs4go/encoding v0.0.0-20210720054946-fe0a4a6f4c7a
	github.com/libs4go/errors v0.0.3
	github.com/libs4go/fixed v0.0.4
	github.com/libs4go/jsonrpc v0.0.0-20210720025424-ccdb148bd313
	github.com/libs4go/scf4go v0.0.1
	github.com/libs4go/slf4go v0.0.4
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
)

replace github.com/libs4go/crypto => ../crypto

replace github.com/libs4go/jsonrpc => ../jsonrpc
