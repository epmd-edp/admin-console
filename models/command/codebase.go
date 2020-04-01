package command

type CreateCodebase struct {
	Name                string      `json:"name" valid:"Required;Match(/^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	Strategy            string      `json:"strategy"`
	Lang                string      `json:"lang" valid:"Required"`
	Framework           *string     `json:"framework,omitempty"`
	BuildTool           string      `json:"buildTool" valid:"Required"`
	TestReportFramework *string     `json:"testReportFramework"`
	MultiModule         bool        `json:"multiModule,omitempty"`
	Type                string      `json:"type,omitempty" valid:"Required"`
	Repository          *Repository `json:"repository,omitempty"`
	Route               *Route      `json:"route,omitempty"`
	Database            *Database   `json:"database,omitempty"`
	Vcs                 *Vcs        `json:"vcs,omitempty"`
	Description         *string     `json:"description,omitempty"`
	Username            string      `json:"username"`
	GitServer           string      `json:"gitServer"`
	Versioning          Versioning  `json:"versioning"`
	GitUrlPath          *string     `json:"gitUrlPath"`
	JenkinsSlave        string      `json:"jenkinsSlave"`
	JobProvisioning     string      `json:"jobProvisioning"`
	DeploymentScript    string      `json:"deploymentScript"`
}

type Repository struct {
	Url      string `json:"url,omitempty" valid:"Required;Match(/(?:^git|^ssh|^https?|^git@[-\\w.]+):(\\/\\/)?(.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$/)"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vcs struct {
	Login    string `json:"login,omitempty" valid:"Required"`
	Password string `json:"password,omitempty" valid:"Required"`
}

type Route struct {
	Site string `json:"site,omitempty" valid:"Match(/^$|^[a-z][a-z0-9-]*[a-z0-9]$/)"`
	Path string `json:"path,omitempty" valid:"Match(/^\\/.*$/)"`
}

type Database struct {
	Kind     string `json:"kind,omitempty" valid:"Required"`
	Version  string `json:"version,omitempty" valid:"Required"`
	Capacity string `json:"capacity,omitempty" valid:"Required"`
	Storage  string `json:"storage,omitempty" valid:"Required"`
}

type DeleteCodebaseCommand struct {
	Name string `json:"name"`
}

type Versioning struct {
	Type      string  `json:"type" valid:"Required"`
	StartFrom *string `json:"startFrom,omitempty"`
}
