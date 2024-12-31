# AWS Devops/Kubernetes microservice CI/CD
This DevOps project employs a comprehensive CI pipeline to automate the development and deployment process. The architecture emphasizes security, performance, and reliability, integrating industry-leading tools and practices.

!!! note "TODO"
    - Merge the AWS docs with the microservices (AKS Readme)
    - Annotate that infra-steps is used to deploy instead of main

<!--
!!! warning "Warning"
    Warning Text

!!! error "Error"
    Error

-->

!!! hint This project is a simple java web app that allows users to post their thoughts and blog digitally. Its mostly used to prove that the pipeline is working.

This DevOps project employs a comprehensive CI/CD pipeline to automate the development and deployment process. The architecture emphasizes security, performance, and reliability, integrating industry-leading tools and practices.

## My DevOps Scripting Examples

- My current examples of this project are located here:
  - This included Terraform, Bash scripts and other scripting
  - https://github.com/T-Py-T/devops-install-scripts

## Best Practices Followed

1. **Automation**: The build, test, and deployment process is automated, reducing the risk of human error and speeding up the cycle times. Automation ensures that every code change is tested and validated before deployment.
2. **Security First**: Integrating Aqua Trivy and SonarQube ensures that security vulnerabilities and code quality issues are detected and addressed early in the pipeline, fostering a secure development lifecycle.
3. **Scalability**: Kubernetes provides a scalable infrastructure that can handle fluctuating loads, ensuring consistent performance during peak traffic.
4. **Observability**: Using Grafana and Prometheus allows real-time monitoring, enabling proactive identification and resolution of potential issues before they impact users.
5. **Version Control and Code Review**: GitHub serves as the foundation for collaboration and quality control, ensuring that only well-reviewed, high-quality code reaches production.

## Architecture

The CI/CD pipeline is depicted in the diagram below, which mirrors the "as-built" system, showcasing the tools and workflows utilized.

![Architecture Diagram](docs/img/CICD-Architechture.png)

### Key Components

#### A. **Source Code Management**

- **GitHub**: 
  - Serves as the backbone of version control, ensuring seamless collaboration among team members.
  - Pull requests and branch strategies help enforce coding standards and encourage peer reviews.
  - Integrated with Jenkins to trigger automated builds and tests upon code commits, ensuring continuous integration.

#### B. **Build and Test Automation**

## Branch Structure and CI/CD Workflow

Each microservice in this project has its own branch in the repository. This ensures that changes to one service do not affect others and allows for independent development and deployment.

### Branches

- `frontend`: Contains the code and Jenkinsfile for the `frontend` service.
- `cartservice`: Contains the code and Jenkinsfile for the `cartservice` service.
- `productcatalogservice`: Contains the code and Jenkinsfile for the `productcatalogservice` service.
- `currencyservice`: Contains the code and Jenkinsfile for the `currencyservice` service.
- `paymentservice`: Contains the code and Jenkinsfile for the `paymentservice` service.
- `shippingservice`: Contains the code and Jenkinsfile for the `shippingservice` service.
- `emailservice`: Contains the code and Jenkinsfile for the `emailservice` service.
- `checkoutservice`: Contains the code and Jenkinsfile for the `checkoutservice` service.
- `recommendationservice`: Contains the code and Jenkinsfile for the `recommendationservice` service.
- `adservice`: Contains the code and Jenkinsfile for the `adservice` service.
- `loadgenerator`: Contains the code and Jenkinsfile for the `loadgenerator` service.
- `infra-steps`: Contains the  deployment.yaml and Jenkinsfile. This is where we can run updates for ArgoCD to identify changes the deployment (image tag/version)

### Example Jenkinsfile

Each branch has a `Jenkinsfile` that defines the CI/CD pipeline for that specific service. This allows language specific testing automation pipelines.

Below is an example of what the adservice Jenkinsfile looks like (Java app):

```groovy
pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git branch: 'adservice', url: 'https://github.com/T-Py-T/eks-jenkins-microservices-cicd.git'
            }
        }
        stage('Build') {
            steps {
                // Build the application
                sh './gradlew build'
            }
        }
        stage('Test') {
            steps {
                // Run unit tests
                sh './gradlew test'
            }
        }
        stage('Build & Tag Docker Image') {
            steps {
                script {
                    withDockerRegistry(credentialsId: 'docker-cred', toolName: 'docker') {
                        sh "docker build -t tnt850910/adservice:latest ."
                    } } } }
        
        stage('Push Docker Image') {
            steps {
                script {
                    withDockerRegistry(credentialsId: 'docker-cred', toolName: 'docker') {
                        sh "docker push tnt850910/adservice:latest "
                    } } } }
    }
    post {
        always {
            // Clean up workspace
            cleanWs()
        }
    }
}


- **Jenkins**:
  - Orchestrates the CI/CD pipeline, ensuring that builds, tests, and deployments are fully automated.
  - Integrates with tools like GitHub and Docker to create a streamlined process from code commit to deployment.
  - Provides real-time feedback to developers about build status and test results.

  ![Jenkins Multibranch](docs/img/jenkins-multibranch.png)

- **Jenkins-TODO**:
  - Relies on external Terraform setup for environment to work
    - When cluster is built API changes
    - Permissions for Jenkins need to be created with cluster (kubectl)
  - Call Terraform in pipeline
    - Allow for tear down of resources after load testing is completed.
    - Pipeline created cluster can be used for load or AB testing.

- **Without Terraform :**
- Kubernetes deployment (apply and get pods) fails
![Jenkins Pipeline](docs/img/Jenkins-Pipeline.png)
- **With Terraform :**
![Completed Kube Deployment](docs/img/Completed-Kube-Deployment.png)

- **Maven**:
  - Simplifies dependency management and builds process for Java projects.
  - Ensures that all dependenciess ns are resolved before building the application, reducing errors and inconsistencies.

#### C. **Security Scanning**

- **Aqua Trivy**:
  - Scans Docker images and source code for vulnerabilities, ensuring that potential security issues are caught before deployment.
  - Generates detailed reports that can be used to address vulnerabilities promptly.

  ![Trivy Scan](docs/img/defaultImage.png)

- **SonarQube**:
  - Conducts comprehensive code analysis to identify bugs, code smells, and security vulnerabilities.
  - Provides actionable insights to improve code quality and enforce compliance with coding standards.

  ![Sonar Report](docs/img/defaultImage.png)

#### D. **Containerization**

- **Docker**:
  - Packages the Java application into lightweight, portable containers, ensuring consistent environments across development, testing, and production stages.
  - Simplifies deployment by abstracting underlying infrastructure differences.

  *Callout Area*: Showcase Dockerfile and container registry management, explaining how this enables rapid, reliable deployments.

#### E. **Container Orchestration**

- **Elastic Kubernetes Service (EKS)**:
  - Manages the deployment and scaling of containerized applications in a highly available environment.
  - Ensures zero downtime by automatically scaling and redistributing workloads as needed.
  - Namespace configurations (e.g., `webapps` and `namespace 2`) isolate different parts of the system for better organization and security. The second namespace is not currently used, but is planned for a similar python web app
  - The configuration for EKS was update from the **terraform.tf** listed in the linked repo and shown implemented below in a later section.

- **EKS Nodes**
  ![EKS Nodes Image](docs/img/EKS-Nodes.png)

- **EKS CLuster**
  ![EKS Cluster Image](docs/img/EKS-Cluster.png)

- **EKS Networking**
  ![EKS Cluster Image](docs/img/EKS-Networking.png)

#### F. **Monitoring and Observability**

- **Prometheus**:
  - Collects metrics from various components of the application and infrastructure, providing deep insights into system health and performance.
  - Supports custom queries to detect anomalies and trigger alerts proactively.

  ![Prometheus Image](docs/img/Prometheus.png)

- **Grafana**:
  - Provides user-friendly dashboards for visualizing Prometheus metrics.
  - Enables stakeholders to monitor key performance indicators (KPIs) in real-time, ensuring system reliability.

  ![Grafana Image](docs/img/Grafana.png)

  *Callout Area*: Include snapshots of Grafana dashboards and Prometheus query outputs, demonstrating the observability aspect of the pipeline.

#### G. **Infrastructure as Code (IaC)**

- **Terraform**:
  - Automates the provisioning and management of infrastructure required for the Kubernetes stack that hosts the Java application.
  - Ensures infrastructure consistency and repeatability by defining it as code.
  - The following key AWS resources are provisioned:
    - **VPC**: Creates a virtual private cloud for network isolation.
    - **Subnets**: Two public subnets in `us-east-1a` and `us-east-1b` availability zones.
    - **Internet Gateway**: Provides internet access to the resources within the VPC.
    - **Route Tables and Associations**: Configures routing for the subnets to allow public internet access.
    - **Security Groups**: Defines rules for cluster and node communication, ensuring controlled ingress and egress.
    - **EKS Cluster**: Deploys an Elastic Kubernetes Service cluster for managing the application containers.
    - **EKS Node Group**: Provisions a scalable worker node group with `t2.large` instances to support container workloads.
    - **IAM Roles and Policies**: Configures roles and permissions for both the EKS cluster and node group to interact with AWS services.
  - Facilitates rapid updates and scaling of infrastructure to match application requirements.

``` bash
terraform plan
```

![Terraform Plan](docs/img/TerraformPlan.png)

``` bash
terraform apply --auto-approve
```

![Terraform Apply](docs/img/TerraformApply.png)
![Terraform Output](docs/img/Terraform-Output.png)

#### H. **AWS Integration**

- **EC2**: Used to house the servers that control the CI/CD process and handles the actions.
  ![EC2 Instances](docs/img/EC2_Instances.png)

- **VPC**:
  - Ensures a secure and isolated environment for hosting applications and infrastructure.
- **S3 Bucket**:
  - Stores artifacts, logs, and backups, ensuring durability and availability.
- **Route 53**:
  - Manages domain names and routes traffic efficiently to ensure seamless user experiences.
- **CloudFront CDN**:
  - Enhances performance by caching content close to end users.
- **Network Load Balancer**:
  - Distributes incoming traffic to backend services, improving fault tolerance and scalability.
