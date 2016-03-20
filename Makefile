PROJECT ?= mqtt-proxy
EB_BUCKET ?= ninjablocks-sphere-docker

APP_NAME ?= mqtt-proxy
APP_ENV ?= mqtt-proxy-prod

SHA1 := $(shell git rev-parse --short HEAD | tr -d "\n")

DOCKERRUN_FILE := Dockerrun.aws.json
APP_FILE := ${SHA1}.zip

build:
	docker build -t "ninjablocks/${PROJECT}:${SHA1}" .

push:
	docker push "ninjablocks/${PROJECT}:${SHA1}"

services:
	docker run --name ninja-rabbit -p 5672:5672 -p 15672:15672 -d mikaelhg/docker-rabbitmq

local:
	docker run -t -i --rm --link ninja-rabbit:rabbit -e "DEBUG=true" \
		-p 6300:6300 -t "ninjablocks/${PROJECT}:${SHA1}"

deploy:
	sed "s/<TAG>/${SHA1}/" < Dockerrun.aws.json.template > ${DOCKERRUN_FILE}
	zip -r ${APP_FILE} ${DOCKERRUN_FILE} .ebextensions

	aws s3 cp ${APP_FILE} s3://${EB_BUCKET}/${APP_ENV}/${APP_FILE}

	aws elasticbeanstalk create-application-version --application-name ${APP_NAME} \
	   --version-label ${SHA1} --source-bundle S3Bucket=${EB_BUCKET},S3Key=${APP_ENV}/${APP_FILE}

	# # Update Elastic Beanstalk environment to new version
	aws elasticbeanstalk update-environment --environment-name ${APP_ENV} \
       --version-label ${SHA1}

clean:
	rm *.zip || true
	rm ${DOCKERRUN_FILE} || true
	rm -rf build || true

.PHONY: all build push local services deploy clean
