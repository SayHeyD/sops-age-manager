build:
	@echo "Building application... 🔄"
	@go build -o bin/sam . && echo "Finished building ✅" || (echo echo "Build failed ❌"; exit 1)

test:
	@echo "Running test... 🔄"
	@go test ./... -v -parallel 4 -cover && echo "Tests finished successfully ✅" || (echo echo "Tests finished with errors ❌"; exit 1)

clean:
	@echo "Cleaning build products... 🔄"
	@rm -f ./bin/sam* && echo "Cleaning done ✅" || (echo echo "Cleaning failed ❌"; exit 1)