.PHONY: create-keypair

PWD = $(shell pwd)
ACCOUNTPATH = $(PWD)/account

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCOUNTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCOUNTPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCOUNTPATH)/rsa_public_$(ENV).pem