


provider: provider.build provider.test provider.install

provider.build:
	cd terraform/check && go build

provider.install:
	cd terraform/check && go install

provider.test:
	cd terraform/check && TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

pulumi: pulumi.schema pulumi.bridge pulumi.install

pulumi.schema:
	cd pulumi/check/provider/cmd/pulumi-tfgen-check && go run main.go schema -o ../pulumi-resource-check

pulumi.bridge:
	cd pulumi/check/provider/cmd/pulumi-resource-check && go build

pulumi.install:
	cd pulumi/check/provider/cmd/pulumi-resource-check && go install

pulumi.test:
	cd pulumi/check/provider/cmd/pulumi-resource-check/test && pulumi up --stack dev