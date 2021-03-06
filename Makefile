it-storage:
	GANACHE_HOST="http://127.0.0.1:7545" \
	CONTRACT_ADDR=0x885b4C56eD4BB68f4CBFA10e379D7B4B1f6DCabd \
	ACCOUNT_PRIVATE_KEY=f778ae7e8a2859d36a77fe3ec0a66d4c58e6cfdc70508154cde53c88ed23b5bf \
	GANHACE_NETWORK_ID=5777 \
		go run ./cmd/storage/...

it-coin:
	GANACHE_HOST="http://127.0.0.1:7545" \
	CONTRACT_ADDR=0xacF34e0d0EFB95759955daC09d2E9170E48405f1 \
	GANHACE_NETWORK_ID=5777 \
		go run ./cmd/coin/...

it-coin-listen:
	GANACHE_HOST="ws://127.0.0.1:7545" \
	CONTRACT_ADDR=0xacF34e0d0EFB95759955daC09d2E9170E48405f1 \
	ACCOUNT_PRIVATE_KEY=0d16bf84ca3e99ad2e831c96d56ef8a87c049bf789a9cd327f433208f90d05d3 \
	GANHACE_NETWORK_ID=5777 \
		go run ./cmd/coin/... -listen

it-coin-transfer:
	GANACHE_HOST="http://127.0.0.1:7545" \
	CONTRACT_ADDR=0xacF34e0d0EFB95759955daC09d2E9170E48405f1 \
	ACCOUNT_PRIVATE_KEY=0d16bf84ca3e99ad2e831c96d56ef8a87c049bf789a9cd327f433208f90d05d3 \
	GANHACE_NETWORK_ID=5777 \
		go run ./cmd/coin/... -transfer -to=0x7A849379659E2B850d223112cdb474f2F9C2CC0a -amount=5

it-coin-mint:
	GANACHE_HOST="http://127.0.0.1:7545" \
	CONTRACT_ADDR=0xacF34e0d0EFB95759955daC09d2E9170E48405f1 \
	ACCOUNT_PRIVATE_KEY=1d621886c847c05e83949546d1908b7bdda0ae1bae1ce2df226d6c577addc3f2 \
	GANHACE_NETWORK_ID=5777 \
		go run ./cmd/coin/... -mint -to=0x9CCdD5CA42cE0448FF13E2EBdb592600F3CF78C5 -amount=10