# Add Applications

Admin Console allows to create, clone, import an application and add it to the environment with its subsequent deployment in Gerrit and building of the Code Review and Build pipelines in Jenkins. 

To add an application, navigate to the **Applications** section on the left-side navigation bar and click the Create button.

Once clicked, the six-step menu will appear: 

* The Codebase Info Menu
* The Application Info Menu
* The Advanced Settings Menu
* The Version Control System Info Menu
* The Exposing Service Info Menu
* The Database Menu

_**NOTE**: The Version Control System Info menu is available in case this option is predefined._

## The Codebase Info Menu

![add-app1](../readme-resource/addapp1.png "add-app1")

1. In the **Codebase Integration Strategy** field, select the necessary option that is the configuration strategy for the replication with Gerrit:
    - Create – creates a project on the pattern in accordance with an application language, a build tool, and a framework.
    - Clone – clones the indicated repository into EPAM Delivery Platform. While cloning the existing repository, you have to fill in the additional fields as well.
    - Import - allows configuring a replication from the Git server. While importing the existing repository, you have to select the Git server and define the respective path to the repository.
               
    _**NOTE**: In order to use the import strategy, make sure to adjust it by following the [Adjust Import Strategy](../documentation/import-strategy.md) page._ 

2. In the **Git Repository URL** field, specify the link to the repository that is to be cloned.
3. Select the **Codebase Authentication** check box and fill in the requested fields:
    - Repository Login – enter your login data.
    - Repository password (or API Token) – enter your password or indicate the API Token.
    
    _**NOTE**: The Codebase Authentication check box should be selected just in case you clone the private repository. If you define the public one, there is no need to enter credentials._ 
4. Click the Proceed button to be switched to the next menu.

    ## The Application Info Menu

    ![add-app2](../readme-resource/addapp2.png "add-app2")

5. Type the name of the application in the **Application Name** field by entering at least two characters and by using the lower-case letters, numbers and inner dashes.

    _**INFO:** If the Import strategy is used, the Application Name field will not be displayed._
    
6. Select any of the supported application languages with its framework in the **Application Code Language/framework** field:

    - Java – selecting Java allows using Java 8 or Java 11.
    - JavaScript - selecting JavaScript allows using the React framework.
    - DotNet - selecting DotNet allows using the DotNet v.2.1 and DotNet v.3.1.
    - Go - selecting Go allows using the Beego and Operator SDK frameworks.
    - Python - selecting Python allows using the Python v.3.8.
    - Other - selecting Other allows extending the default code languages when creating a codebase with the clone/import strategy. To add another code language, inspect the [Add Other Code Language](add_other_code_language.md) section.

    _**NOTE**: The Create strategy does not allow to customize the default code language set._
    
7. Choose the necessary build tool in the Select Build Tool field:

    - Java - selecting Java allows using the Gradle or Maven tool.
    - JavaScript - selecting JavaScript allows using the NPM tool.
    - .Net - selecting .Net allows using the .Net tool.

    _**NOTE**: The Select Build Tool field disposes of the default tools and can be changed in accordance with the selected code language._
8. Select the **Multi-Module Project** check box that becomes available if the Java code language and the Maven build tool are selected. 

    _**NOTE**: If your project is a multi-modular, add a property to the project root POM-file:_

    `<deployable.module> for a Maven project.`

    `<DeployableModule> for a DotNet project.`

9. Click the Proceed button to be switched to the next menu.

    ## The Advanced Settings Menu

    ![add-app3](../readme-resource/add_app_250.png "add-app3")

10. Select CI pipeline provisioner that will be handling a codebase. For details, refer to the [Add Job Provision](https://github.com/epmd-edp/jenkins-operator/blob/master/documentation/add-job-provision.md#add-job-provision) instruction and become familiar with the main steps to add an additional job provisioner.
11. Select Jenkins slave that will be used to handle a codebase. For details, refer to the [Add Jenkins Slave](https://github.com/epmd-edp/jenkins-operator/blob/master/documentation/add-jenkins-slave.md#add-jenkins-slave) instruction and inspect the steps that should be done to add a new Jenkins slave.  
12. Select the necessary codebase versioning type:
     
     * **default** - the previous versioning logic that is realized in EDP Admin Console 2.2.0 and lower versions. Using the default versioning type, in order to specify the version of the current artifacts, images, and tags in the Version Control System, a developer should navigate to the corresponding file and change the version **manually**.
      
     * **edp** - the new versioning logic that is available in EDP Admin Console 2.3.0 and subsequent versions. Using the edp versioning type, a developer indicates the version number from which all the artifacts will be versioned and, as a result, **automatically** registered in the corresponding file (e.g. pom.xml). 
     
       When selecting the edp versioning type, the extra field will appear:
         
       ![add-app3_2](../readme-resource/addapp3_2.png "add-app3_2")
     
       a. Type the version number from which you want the artifacts to be versioned.
     
     _**NOTE**: The Start Version From field should be filled out in compliance with the semantic versioning rules, e.g. 1.2.3 or 10.10.10._
                      
13. In the **Select Deployment Script** field, specify one of the available options: helm-chart / openshift-template that are predefined in case it is OpenShift or EKS.  
14. In the **Select CI Tool** field, choose the necessary tool: Jenkins or GitLab CI, where Jenkins is the default tool and
the GitLab CI tool can be additionally adjusted. For details, please refer to the [Adjust GitLab CI Tool](../documentation/ci-tool.md) page.

    >_**NOTE**: The GitLab CI tool is available only with the Import strategy and makes the **Jira integration** feature unavailable._

15. Select the **Integrate with Jira Server** checkbox in case it is required to connect Jira tickets with the commits and have a respective label in the Fix Version field.
    >_**NOTE**: To adjust the Jira integration functionality, first apply the necessary changes described on the [Adjust Integration With Jira Server](../documentation/jira-server.md) page, and setup the [VCS Integration With Jira Server](../documentation/jira_vcs_integration.md). Pay attention that the Jira integration feature is not available when using the GitLab CI tool._ 
                                                                                                                                                                                 
    ![add-app3](../readme-resource/add_app3_ji2.png "add-app3_2")

16. As soon as the Jira server is set, select it in the **Select Jira Server** field.
17. Indicate the pattern using any character, which is followed on the project, to validate a commit message.
18. Indicate the pattern using any character, which is followed on the project, to find a Jira ticket number in a commit message.
19. Click the Proceed button to be switched to the next menu.

    ## The Version Control System Info Menu

    ![add-app4](../readme-resource/addapp_4.png "add-app4")
    
20. Enter the login credentials into the **VCS Login** field.
21. Enter the password into the **VCS Password (or API Token)** field OR add the API Token.
22. Click the Proceed button to be switched to the next menu.
    
    >_**NOTE**: The VCS Info step is skipped in case there is no need to integrate the version control for the application deployment. If the cloned application includes the VCS, this step should be completed as well._

    ## The Exposing Service Info Menu

    ![add-app5](../readme-resource/addapp_5.png "add-app5")

23. Select the **Need Route** check box to create a route component in the OpenShift project for the externally reachable host name. As a result, the added application will be accessible in a browser.
    
    Fill in the necessary fields:
    
    - Name – type the name by entering at least two characters and by using the lower-case letters, numbers and inner dashes. The mentioned name will be as a prefix for the host name.
    - Path – specify the path starting with the **/api** characters. The mentioned path will be at the end of the URL path.
    
24. Click the Proceed button to be switched to the final menu.

    ## The Database Menu

    ![add-app6](../readme-resource/addapp_6.png "add-app6")

25. Select the **Need Database** check box in case you need a database. Fill in the required fields:
    
    - Database – the PostgreSQL DB is available by default.
    - Version – the latest version (postgres:9.6) of the PostgreSQL DB is available by default.
    - Capacity – indicate the necessary size of the database and its unit of measurement (Mi – megabyte, Gi – gigabyte, Ti – terabyte). There is no limit for the database capacity.
    - Persistent storage – select one of the available storage methods: efs or gp2.
    
26. Click the Create button. Once clicked, the CONFIRMATION summary will appear displaying all the specified options and settings, click Continue to complete the application addition.
    
>_**NOTE**: After the complete adding of the application, please refer to the [Inspect Application](../documentation/inspect_application.md) page._

### Related Articles

* [Inspect Application](../documentation/inspect_application.md)
* [Delivery Dashboard Diagram](../documentation/d_d_diagram.md)
---
* [Add CD Pipelines](../documentation/add_CD_pipelines.md)
* [Adjust Integration With Jira Server](../documentation/jira-server.md)
* [Adjust VCS Integration With Jira Server](../documentation/jira_vcs_integration.md)