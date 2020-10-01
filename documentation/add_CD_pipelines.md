# Add CD Pipelines

Admin Console provides the ability to deploy an environment on your own and specify the essential components as well. 

Navigate to the **Continuous Delivery** section on the left-side navigation bar and click the Create button.
Once clicked, the four-step menu will appear: 

* The Pipeline Menu
* The Applications Menu
* The Stages Menu
* The Services Menu

The creation of the CD pipeline becomes available as soon as an application is created including its provisioning in a branch and the necessary entities for the environment.

After the complete adding of the CD pipeline, inspect the [Check CD Pipeline Availability](#Check_CD_Pipeline_Availability) part.

## The Pipeline Menu

![addcdpipe1](../readme-resource/addcdpipe1.png "addcdpipe1")

1. Type the name of the pipeline in the **Pipeline Name** field by entering at least two characters and by using the lower-case letters, numbers and inner dashes.
2. Click the Proceed button to be switched to the next menu.

    ## The Applications Menu
    
    ![addcdpipe2](../readme-resource/addcdpipe2.png "addcdpipe2")

3. Select the check box of the necessary application in the **Applications** menu.
4. Specify the necessary codebase Docker stream (the output for the branch and other stages from other CD pipelines) from the drop-down menu.
5. Select the **Promote in pipeline** check box in order to transfer the application from one to another stage by the specified codebase Docker stream. If the Promote in pipeline check box is not selected, the same codebase Docker stream will be deployed regardless of the stage, i.e. the codebase Docker stream input, which was selected for the pipeline, will be always used.

    _**NOTE**: The newly created CD pipeline has the following pattern combination: [pipeline name]-[branch name]. If there is another deployed CD pipeline stage with the respective codebase Docker stream (= image stream as an OpenShift term), the pattern combination will be as follows: [pipeline name]-[stage name]-[application name]-[verified];_
    
6. Click the Proceed button to be switched to the next menu.

    ## The Stages Menu

    ![addcdpipe3](../readme-resource/addcdpipe3.png "addcdpipe3")

7. Click the plus sign icon in the **Stages** menu and fill in the necessary fields in the Adding Stage window:

    ![addcdpipe4](../readme-resource/addcdpipe4.png "addcdpipe4")

    a. Type the stage name;
    
    b. Enter the description for this stage;
    
    c. Select the quality gate type:

      - Manual - means that the promoting process should be confirmed in Jenkins manually; 

      - Autotests - means that the promoting process should be confirmed by the successful passing of the autotests. In the additional fields, select the previously created autotest name and specify its branch for the autotest that will be launched on the current stage.

    d. Type the step name, which will be displayed in Jenkins, for every quality gate type;

    e. Select the trigger type that allows promoting images to the next environment. Now the default manual is available; 

    _**INFO**: Add an unlimited number of quality gates by clicking a corresponding plus sign icon and remove them as well by clicking the recycle bin icon._

    _**NOTE**:Execution sequence. The image promotion and execution of the pipelines depend on the sequence in which the environments are added._

8. Click the Add button to display it in the Stages menu:

    _**INFO**: Perform the same steps as described above if there is a necessity to add one more stage._
    
    ![addcdpipe5](../readme-resource/addcdpipe5.png "addcdpipe5")

9. Edit the stage by clicking its name and applying changes, and remove the added stage by clicking the recycle bin icon next to its name.

    ## The Services Menu
    
    ![addcdpip6](../readme-resource/addcdpipe6.png "addcdpipe6")

10. Choose the necessary services by selecting the related check box(ex).

11. Click the Create button to start the provisioning of the pipeline. After the CD pipeline is added, the new project with the stage name will be created in OpenShift.

## <a name="Check_CD_Pipeline_Availability"></a> Check CD Pipeline Availability

As soon as the CD pipeline is provisioned and added to the CD Pipelines list, there is an ability to:

1. Remove the added CD pipeline:

    ![removepipe](../readme-resource/removecdpipe.png "removepipe")

2. Edit the CD pipeline by clicking the pen icon next to its name in the CD Pipelines list:

    ![addcdpip7](../readme-resource/addcdpipe7.png "addcdpipe7")

    - apply the necessary changes and click the Proceed button to confirm the editions:
    
    ![addcdpip8](../readme-resource/addcdpipe8.png "addcdpipe8")

3. Check the CD pipeline data and details by clicking the CD pipeline name in the CD Pipelines list:

    - the main link on the top of the details page refers to Jenkins;

    - the pen icon refers to the same **Edit CD Pipeline** page as mentioned above and allows to apply the necessary changes;

    - the Applications menu has the main information about the applications with the respective codebase Docker streams and links to Jenkins and Gerrit as well as the signification of the promotion in CD pipeline; 

    - the Stages menu includes the stages data that was previously mentioned, the direct links to the respective to every stage OpenShift page, and the link to the Autotest details page in case there are added autotests.
    
      _**NOTE**: In order to **remove a stage**, click the corresponding delete icon, type the CD pipeline name and confirm the deletion by clicking the Delete button. If you remove the last stage, the whole CD pipeline will be removed as the CD pipeline does not exist without stages._
      
    ![addcdpip9](../readme-resource/addcdpipe9.png "addcdpipe9")

    - the Third-Party Services menu displays the name, description, and version of the added service;

    - the Deployed Version menu indicates the applications and stages with the appropriate status. The status will be changed after stage deployment. 

    ![addcdpip10](../readme-resource/addcdpipe10.png "addcdpipe10")

    - the Status Info menu displays all the actions that were performed during the deployment process:
    
    ![addcdpip11](../readme-resource/addcdpipe11.png "addcdpipe11")
   
>_**INFO**: If there is a necessity to create another CD pipeline, navigate to the Continuous Delivery section, click the Create button and perform the same steps as described above._    

### Related Articles

* [EDP Admin Console](https://github.com/epmd-edp/admin-console#edp-admin-console)
* [Delivery Dashboard Diagram](../documentation/d_d_diagram.md)