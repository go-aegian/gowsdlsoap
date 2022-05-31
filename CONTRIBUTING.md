### Thanks for your interest in contributing to this project and for taking the time to read this guide.

## Development workflow
Go is unlike any other language in that it forces a specific development workflow and project structure. Trying to fight it is useless, frustrating and time consuming. So, you better be prepare to adapt your workflow when contributing to Go projects.

### Prerequisites
1. [Install Go][go-install].
2. Download the sources and switch the working directory:

    ```bash
    go get -u -d github.com/go-aegian/gowsdlsoap
    cd $GOPATH/src/github.com/go-aegian/gowsdlsoap
    ```

### Pull Requests
* Please be generous describing your changes.
* Although it is highly suggested to include tests, they are not a hard requirement in order to get your contributions accepted.
* Keep pull requests small so core developers can review them quickly.
* Make sure you run `go fmt` to format your code before submitting your pull request.

### Workflow for third-party code contributions
* In Github, fork `https://github.com/go-aegian/gowsdlsoap` to your own account
* Get the package using "go get": `go get github.com/go-aegian/gowsdlsoap`
* Move to where the package was cloned: `cd $GOPATH/src/github.com/go-aegian/gowsdlsoap/`
* Add a git remote pointing to your own fork: `git remote add downstream git@github.com:<your_account>/gowsdlsoap.git`
* Create a branch for making your changes, then commit them.
* Push changes to downstream repository, this is your own fork: `git push <mybranch> downstream`
* In Github, from your fork, create the Pull Request and send it upstream.
* You are done.

#### A typical workflow is:

1. [Fork the repository.][fork] [This tip maybe also helpful.][go-fork-tip]
2. [Create a topic branch.][branch]
3. Add tests for your change.
4. Run `go test`. If your tests pass, return to the step 3.
5. Implement the change and ensure the steps from the previous step pass.
6. Run `goimports -w .`, to ensure the new code conforms to Go formatting guideline.
7. [Add, commit and push your changes.][git-help]
8. [Submit a pull request.][pull-req]

### Workflow for core developers
Since core developers usually have access to the upstream repository, there is no need for having a workflow like the one for thid-party contributors.

* Get the package using "go get": `go get github.com/go-aegian/gowsdlsoap`
* Create a branch for making your changes, then commit them.
* Push changes to the repository: `git push origin <mybranch>`
* In Github, create the Pull Request from your branch to master.
* Before merging into master, wait for at least two developers to code review your contribution.


## Issues
* Before reporting an issue make sure you search first if anybody has already reported a similar issue and whether or not it has been fixed.
* Make sure your issue report sufficiently details the problem.
* Include code samples reproducing the issue.
* Please do not derail or troll issues. Keep the discussion on topic and respect the Code of conduct.
* Please do not open issues for personal support requests, use the mailing list instead.
* If you want to tackle any open issue, make sure you let people know you are working on it.


## Resources
* **W3C WSDL spec:** http://www.w3.org/TR/wsdl
* **W3C SOAP 1.2 spec:** http://www.w3.org/TR/soap/

[go-install]: https://golang.org/doc/install

[go-fork-tip]: http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html

[fork]: https://help.github.com/articles/fork-a-repo

[branch]: http://learn.github.com/p/branching.html

[git-help]: https://guides.github.com

[pull-req]: https://help.github.com/articles/using-pull-requests
