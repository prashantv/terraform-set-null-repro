PHONY: install
install:
	go build -o terraform-provider-tftest
	mkdir -p ~/.terraform.d/plugins/local/prashantv/tftest/0.1.0/darwin_amd64
	mv terraform-provider-tftest ~/.terraform.d/plugins/local/prashantv/tftest/0.1.0/darwin_amd64
	make clean-tf

.PHONY: clean
clean: clean-tf
	rm -rf ~/.terraform.d/plugins/local/prashantv/tftest

.PHONY: clean-tf
clean-tf:
	cd repro && rm -rf db .terraform .terraform.lock.hcl terraform.tfstate*

repro/.terraform.lock.hcl:
	cd repro && terraform init

.PHONY: apply
apply: repro/.terraform.lock.hcl
	cd repro && terraform apply
