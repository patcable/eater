TARGETS = eater

eater: 
	GOARCH=amd64 GOOS=linux go build

eater.zip: eater
	zip eater.zip eater foodlog.gohtml style.css

deploy: eater.zip
	aws --profile eater-deploy --no-cli-pager lambda update-function-code --function-name eater --zip-file fileb://eater.zip

clean:
	rm -f ${TARGETS} ${TARGETS}.zip