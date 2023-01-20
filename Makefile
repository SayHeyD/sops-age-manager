build:
	@echo "Building application... 🔄"
	@go build -o bin/sam ./cmd/sops-age-manager
	@echo "Finished building ✅"

clean:
	@echo "Cleaning build products... 🔄"
	@rm -f ./bin/sam*
	@echo "Cleaning done ✅"