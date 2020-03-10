package consts

const (
	//EDP Components
	Jenkins        = "jenkins"
	Openshift      = "openshift"
	Kubernetes     = "kubernetes"
	DockerRegistry = "docker-registry"
	Gerrit         = "gerrit"

	EdpCICDPostfix = "-edp-cicd"

	Application = "application"
	Autotest    = "autotests"
	Library     = "library"

	CodebasePlural       = "codebases"
	CodebaseBranchPlural = "codebasebranches"
	StagePlural          = "stages"
	CDPipelinePlural     = "cdpipelines"
	CodebaseKind         = "Codebase"

	LanguageJava        = "Java"
	FrameworkSpringBoot = "springboot"
)

var DefaultBuildNumber = "0"
