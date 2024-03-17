DEPLOYMENT_NAME=synergychat-web
IMAGE_NAME=bootdotdev/synergychat-web:latest

all:
	kubectl create deployment ${DEPLOYMENT_NAME} --image=${IMAGE_NAME}

