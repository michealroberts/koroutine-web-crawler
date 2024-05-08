#/*****************************************************************************************************************/

#// @author         Michael Roberts

#/*****************************************************************************************************************/

# The name of the base container:
CONTAINER_NAME="koroutine/web-crawler"

#/*****************************************************************************************************************/

# Build and run the development Docker container, with a running shell:
.PHONY: dev
dev:
	@echo "Building..."
	@docker build --target development --tag $(CONTAINER_NAME):dev .
	@echo "Configuring environment..."
	@docker run --volume .:/usr/src/app --tty --interactive $(CONTAINER_NAME):dev /bin/zsh

#/*****************************************************************************************************************/