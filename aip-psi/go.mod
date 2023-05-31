module github.com/greenboxal/aip/aip-psi

go 1.20

replace github.com/greenboxal/aip/aip-controller => ../aip-controller

replace github.com/greenboxal/aip/aip-wiki => ../aip-wiki

replace github.com/greenboxal/aip/aip-sdk => ../aip-sdk

replace github.com/greenboxal/aip/aip-forddb => ../aip-forddb

replace github.com/greenboxal/aip/aip-psi => ../aip-psi

replace github.com/greenboxal/aip/aip-langchain => ../aip-langchain

require (
	github.com/gomarkdown/markdown v0.0.0-20230322041520-c84983bdbf2a
	github.com/ipld/go-ipld-prime v0.20.0
)

require (
	github.com/ipfs/go-cid v0.4.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.8.1 // indirect
	github.com/multiformats/go-multihash v0.2.2 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)
