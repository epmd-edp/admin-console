# Adjust Integration VCS with Jira Server
### Requirements
* VCS Server
* Jira
* Crucible

In order to adjust the Jira integration with VCS, follow the next steps:

1. Configure integration [Jira](https://jiraeu.epam.com) with [Crucible](https://crucible.epam.com/) by creating a request in [Support](https://support.epam.com/).

2. Each project in Gerrit should be integrated with project in [Crucible](https://crucible.epam.com/), by creating a request in [Support](https://support.epam.com/), example of request:

    ![request_example](../readme-resource/сrucible_integration_request.png "request_example")  

3. To link commits with Jira ticket, in commit message set ticket ID following a specific format:

    ![commit_pattern](../readme-resource/commit_pattern.png "commit_pattern")  

4. After that all commits from Gerrit will be displayed on [Crucible](https://crucible.epam.com/).