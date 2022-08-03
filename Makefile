.PHONY: create-keypair
PWD = $(shell pwd)
ACCT_PATH = $(PWD)/api
create-keypair:
	@echo "Creating rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCT_PATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCT_PATH)/rsa_private_$(ENV).pem -pubout -out $(ACCT_PATH)/rsa_public_$(ENV).pem
