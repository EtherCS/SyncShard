.PHONY: build create_testnet run_testnet

build:
	bash ./scripts/go_build_executables.sh

run_test:
	bash ./scripts/run_test.sh
