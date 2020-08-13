# Adjust Import Strategy

In order to use the import strategy, you should add Secret with ssh key, GitServer CR and Jenkins credentials by following the steps:

1. Create Secret in the OpenShift/K8S namespace for the Git account with the **id_rsa**, **id_rsa.pub**, and **username** fields:

    ![secret](../readme-resource/add-secret.png "secret")

   You can use the following example command:

   ```bash
   kubectl create secret generic gitlab -n edp --from-file=id_rsa=id_rsa --from-file id_rsa.pub=id_rsa.pub --from-literal=username=user@gitlab.com
   ``` 
2. Create GitServer CR in the OpenShift/K8S namespace with the **gitHost**, **gitUser**, **httpsPort**, **sshPort**, **nameSshKeySecret**, and **createCodeReviewPipeline** fields:

    ![git-server](../readme-resource/add-git-server.png "git-server")

   You can use the following example template:
   
   ```yaml
   apiVersion: v2.edp.epam.com/v1alpha1
   kind: GitServer
   metadata:
     name: git-epam
   spec:
     createCodeReviewPipeline: false
     gitHost: git.epam.com
     gitUser: git
     httpsPort: 443
     nameSshKeySecret: gitlab-sshkey
     sshPort: 22

   ```
    >*Note: The value of the **nameSshKeySecret** property is the name of the Secret that is indicated in the first point above.*

3. Create a Credential in Jenkins with the same ID as in the **nameSshKeySecret** property, and with the private key. Navigate to **Jenkins -> Credentials -> System -> Global credentials -> Add Credentials**:

    ![credential](../readme-resource/add-credentials.png "credential")
    
4. Change the Deployment Config of the Admin Console by adding the **Import** strategy to the **INTEGRATION_STRATEGIES** variable:

    ![integration-strategy](../readme-resource/add-integretion-strategies.png "integration-strategy")
    
5. As soon as the Admin Console is redeployed, the **Import** strategy will be added to the Create Application page. For details, please refer to the [Add Applications](../documentation/add_applications.md) page.