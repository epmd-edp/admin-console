# Adjust Integration with Jira Server

In order to integrate with Jira server, make sure to add JiraServer CR by following steps:

1. Create Secret in the OpenShift/K8S namespace for Jira Server account with **username** and **password** fields:

    ![jira-server-secret](../readme-resource/add-jira-server-secret.png "jira-server-secret")
 
2. Create JiraServer CR in the OpenShift/K8S namespace with the **apiUrl**, **credentialName** and **rootUrl** fields:

    ![jira-server](../readme-resource/jira-server.png "jira-server")
    
    >*Note: The value of the **credentialName** property is the name of the Secret that is indicated in the first point above.*
                                                                                                                                                                                                    >
3. As soon as everything is configured refer 

    ![jira-server-integration](../readme-resource/addapp3_3.png "jira-server-integration")

