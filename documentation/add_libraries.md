# Add Libraries

Admin Console helps to create, clone or import a library and add it to the environment with its subsequent deployment in Gerrit and building of the Code Review and Build pipelines in Jenkins. 

Navigate to the **Libraries** section on the left-side navigation bar and click the Create button.

Once clicked, the four-step menu will appear:

* The Codebase Info Menu
* The Library Info Menu
* The Advanced Settings Menu
* The Version Control System Info Menu

_**NOTE**: The Version Control System Info menu is available in case this option is predefined._

## The Codebase Info Menu

![addlib1](../readme-resource/addlib1.png "addlib1")

1. In the **Codebase Integration Strategy** field, select the necessary option that is the configuration strategy for the replication with Gerrit:
    - Create – creates a project on the pattern in accordance with a code language, a build tool, and a framework.
    - Clone – clones the indicated repository into EPAM Delivery Platform.
    >_**NOTE**: While cloning the existing repository, you have to fill in the additional fields as well._
    
    - Import - allows configuring a replication from the Git server. While importing the existing repository, you have to select the Git server and define the respective path to the repository.
    >_**NOTE**: In order to use the import strategy, make sure to adjust it by following the [Adjust Import Strategy](../documentation/import-strategy.md) page._ 
    
2. In the **Git Repository URL** field, specify the link to the repository that is to be cloned.
3. Select the **Codebase Authentication** check box and fill in the requested fields:
    - Repository Login – enter your login data.
    - Repository password (or API Token) – enter your password or indicate the API Token.
4. Click the Proceed button to be switched to the next menu.

    ## The Library Info Menu

    ![addlib2](../readme-resource/add_lib_2.png "addlib2")

5. Type the name of the library in the **Library Name** field by entering at least two characters and by using the lower-case letters, numbers and inner dashes.

    _**INFO**: If the Import strategy is used, the Library Name field will not be displayed._

6. Select any of the supported code languages in the **Library Code Language** block:

    - Java – selecting Java allows specify Java 8 or Java 11, and further usage of the Gradle or Maven tool.
    - JavaScript - selecting JavaScript allows using the NPM tool.
    - DotNet - selecting DotNet allows using the DotNet v.2.1 and DotNet v.3.1.
    - Groovy-pipeline - selecting Groovy-pipeline allows having the ability to customize a stages logic. For details, please refer to the [Customize CD Pipeline](https://github.com/epmd-edp/admin-console/blob/master/documentation/cicd_customization/customize-deploy-pipeline.md#customize-cd-pipeline) page.
    - Python - selecting Python allows using the Python v.3.8.
    - Other - selecting Other allows extending the default code languages when creating a codebase with the clone/import strategy. To add another code language, inspect the ([Add Other Code Language](add_other_code_language.md) page.

    _**NOTE**: The Create strategy does not allow to customize the default code language set._

7. The **Select Build Tool** field disposes of the default tools and can be changed in accordance with the selected code language.

8. Click the Proceed button to be switched to the next menu.

    ## The Advanced Settings Menu
    
    ![addlib3](../readme-resource/add_lib_250.png "addlib3")

9. Select the CI pipeline provisioner that will be used to handle a codebase. For details, refer to the [Add Job Provision](https://github.com/epmd-edp/jenkins-operator/blob/master/documentation/add-job-provision.md#add-job-provision) instruction and become familiar with the main steps to add an additional job provisioner.

10. Select Jenkins slave that will be used to handle a codebase. For details, refer to the [Add Jenkins Slave](https://github.com/epmd-edp/jenkins-operator/blob/master/documentation/add-jenkins-slave.md#add-jenkins-slave) instruction and inspect the steps that should be done to add a new Jenkins slave.

11. Select the necessary codebase versioning type:
         
    * **default** - the previous versioning logic that is realized in EDP Admin Console 2.2.0 and lower versions. Using the default versioning type, in order to specify the version of the current artifacts, images, and tags in the Version Control System, a developer should navigate to the corresponding file and change the version **manually**.
          
    * **edp** - the new versioning logic that is available in EDP Admin Console 2.3.0 and subsequent versions. Using the edp versioning type, a developer indicates the version number from which all the artifacts will be versioned and, as a result, **automatically** registered in the corresponding file (e.g. pom.xml). 
         
      When selecting the edp versioning type, the extra field will appear:
             
      ![add-app3_2](../readme-resource/addapp3_2.png "add-app3_2")
         
      a. Type the version number from which you want the artifacts to be versioned.
         
    _**NOTE**: The Start Version From field should be filled out in compliance with the semantic versioning rules, e.g. 1.2.3 or 10.10.10._    
12. In the **Select CI Tool** field, choose the necessary tool: Jenkins or GitLab CI, where Jenkins is the default tool and
    the GitLab CI tool can be additionally adjusted. For details, please refer to the [Adjust GitLab CI Tool](../documentation/ci-tool.md) page.
       >_**NOTE**: The GitLab CI tool is available only with the Import strategy and makes the **Jira integration** feature unavailable._   
13. Select the **Integrate with Jira Server** checkbox in case it is required to connect Jira tickets with the commits and have a respective label in the Fix Version field.
       >_**NOTE**: To adjust the Jira integration functionality, first apply the necessary changes described on the [Adjust Integration With Jira Server](../documentation/jira-server.md) page, and setup the [VCS Integration With Jira Server](../documentation/jira_vcs_integration.md). Pay attention that the Jira integration feature is not available when using the GitLab CI tool._ 
       
       ![add-app3_2](../readme-resource/add_lib3_ji2.png "add-app3_2")
14. As soon as the Jira server is set, select it in the **Select Jira Server** field.
15. Indicate the pattern using any character, which is followed on the project, to validate a commit message.
16. Indicate the pattern using any character, which is followed on the project, to find a Jira ticket number in a commit message.
17. Click the Create button to create a library or click the Proceed button to be switched to the next VCS menu that can be predefined.

    ## The Version Control System Info Menu

    ![addlib4](../readme-resource/addlib4.png "addlib4")

18. Enter the login credentials into the **VCS Login** field.
19. Enter the password into the **VCS Password (or API Token)** field OR add the API Token.
20. Click the Create button, check the CONFIRMATION summary, click Continue to add the library to the Libraries list.

> _**NOTE**: After the complete adding of the library, inspect the [Inspect Library](../documentation/inspect_library.md) part._

### Related Articles

* [Inspect Library](../documentation/inspect_library.md)
* [Delivery Dashboard Diagram](../documentation/d_d_diagram.md)
---
* [Add CD Pipelines](../documentation/add_CD_pipelines.md)
* [Adjust Integration With Jira Server](../documentation/jira-server.md)
* [Adjust VCS Integration With Jira Server](../documentation/jira_vcs_integration.md)