
test:
	go test -v -run TestOriginPatterns -timeout 30s

publish:
	git add . 
	git commit
	git push origin master --ipv4
	@echo "Name of tag you want this commit to have: ";\
	read TAG;
	git tag -a $$(TAG) 
	git push origin --tags 
	