# Adjust Integration with Jira Server

In order to adjust the Jira server integration, first add JiraServer CR by performing the following:

1. Create Secret in the OpenShift/K8S namespace for Jira Server account with the **username** and **password** fields:

    ![jira-server-secret](../readme-resource/add-jira-server-secret.png "jira-server-secret")
 
2. Create JiraServer CR in the OpenShift/K8S namespace with the **apiUrl**, **credentialName** and **rootUrl** fields:

    ![jira-server](../readme-resource/jira-server.png "jira-server")
    
    >_**NOTE**: The value of the **credentialName** property is the name of the Secret, which is indicated in the first point above._
                                                                                                                                                                                                    
3. Being in Admin Console, navigate to the Advanced Settings menu to check that the Integrate with Jira Server check box became available:  

    ![jira-server-integration](../readme-resource/jira_integration_ac.png "jira-server-integration")

