GLIDE_LOCKFILE   = glide.lock
NODE_MODULES_DIR = src/ur/admin/frontend/node_modules
FRONTEND_APP_DIR = src/ur/admin/frontend

start:
	@docker-compose up

stop:
	@docker-compose down

restart:
	@docker-compose restart

run-admin: $(GLIDE_LOCKFILE) $(NODE_MODULES_DIR)
	@go install ur/admin/cmd/admin
	@cd $(FRONTEND_APP_DIR); webpack --watch &
	@PORT=5000 bin/admin --dev

$(GLIDE_LOCKFILE):
	@glide install

$(NODE_MODULES_DIR):
	cd $(FRONTEND_APP_DIR); npm install

update-go-deps:
	@glide update

clean:
	@-rm $(GLIDE_LOCKFILE)
	@-rm -rf $(NODE_MODULES_DIR)

install:
	@docker-compose build
