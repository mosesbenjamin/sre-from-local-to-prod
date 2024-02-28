# SRE From local to production

## This repository houses my implementation of [this SRE Bootcamp Exercises](https://playbook.one2n.in/sre-bootcamp/sre-bootcamp-exercises#:~:text=%F0%9F%92%AA-,SRE%20bootcamp%20exercises,-Part%20One%20%2D%20From)

### Milestones completed so far

- [🏅Create a simple REST API Webserver](#create-a-simple-rest-api-webserver)

## Create a simple REST API Webserver
[Problem statement, Functional requirement and API expectations](https://playbook.one2n.in/sre-bootcamp/sre-bootcamp-exercises/1-create-a-simple-rest-api) are available in the preceeding link.

The Go programming language has been selected to complete this milestone, alongside [chi](https://github.com/go-chi/chi) to provide multiplexing support.

### Usage
#### Prerequisites
- Docker 
- Sqlc (Optional)
- goose (Optional)

#### Quickstart
Run `make help` for a list of commands useful to run the REST API Webserver.

`make run` starts the webserver and you can start interacting with the API on http://localhost:8080/api/v1/. Please see [the postman documentation](#postman-api-documentation) for a list of supportted API endpoints.

Adminer is available at http://localhost:3333 to manage the postgres database created as part of the application setup.

#### Postman API Documentation
The API Documemtation is available via [this Postman Link](https://www.postman.com/mavewrick/workspace/public/collection/11230844-3adc49ff-fd4c-47a1-a2c7-895c642f33f6?action=share&source=copy-link&creator=0)

### What was covered in this implementation
#### Functional Requirement
Using the API we can.
- [x] Add a new student.
- [x] Get all students.
- [x] Get a student with an ID.
- [x] Update existing student information.
- [x] Delete a student record.

#### Expectations
The following expectations have been met as per the milestone requirements.
- [x] Create a public repository on GitHub.
- [x] The repository should contain the following
README.md file explaining the purpose of the repo, along with local setup instructions.
- [x] Explicitly maintaining dependencies in a file ex (pom.xml, build.gradle, go.mod, requirements.txt, etc).
- [x] Makefile to build and run the REST API locally.
- [x] Ability to run DB schema migrations to create the student table.
- [x] Config (such as database URL) should not be hard-coded in the code and should be passed through environment variables.
- [x] Postman collection for the APIs.
#### API expectations
- [x] Support API versioning (e.g., api/v1/<resource>).
- [x] Using proper HTTP verbs for different operations.
- [x] API should emit meaningful logs with appropriate log levels.
- [x] API should have a /healthcheck endpoint.
- [ ] Unit tests for different endpoints.

### Additional features
- [ ] [Fully conforms with the twelve-factor app methodology](https://12factor.net/)
- [x] [Follows Best Practices for REST API design](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
- [x] Password + ID Authentication
- [ ] Third Party Authentication with Auth0 
- [ ] API Keys Authentication for Webhooks
- [ ] Caching with redis