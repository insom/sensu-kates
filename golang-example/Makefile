all:
	docker build -t golang-example .
	docker tag golang-example gcr.io/insom-161401/test/golang-example
	gcloud docker -- push gcr.io/insom-161401/test/golang-example
.PHONY: all
