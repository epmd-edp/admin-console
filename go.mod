module edp-admin-console

go 1.14

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20180801171038-322a19404e37

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/astaxie/beego v1.12.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/epmd-edp/cd-pipeline-operator/v2 v2.3.0-58.0.20200520111933-aefa684da6d1
	github.com/epmd-edp/codebase-operator/v2 v2.3.0-95.0.20200903073301-2eefd21202bb
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/lib/pq v1.0.0
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/openshift/api v3.9.0+incompatible
	github.com/openshift/client-go v3.9.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	go.uber.org/zap v1.14.1
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.10.0
	k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)
