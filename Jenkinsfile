pipeline {
    agent any

    stages {
        stage('Deploy To Kubernetes') { steps {
            withKubeConfig(caCertificate: '', clusterName: 'devopstesting-cluster', contextName: '', credentialsId: 'k8-cred', namespace: 'webapps', restrictKubeConfigAccess: false, serverUrl: 'https://9A41A4B4FD1074711430D76C5221A946.gr7.us-east-1.eks.amazonaws.com') {
                sh "kubectl apply -f deployment-service.yml"
            } } }
        
        stage('Verify the Deployment') { steps {
            withKubeConfig(caCertificate: '', clusterName: 'devopstesting-cluster', contextName: '', credentialsId: 'k8-cred', namespace: 'webapps', restrictKubeConfigAccess: false, serverUrl: 'https://9A41A4B4FD1074711430D76C5221A946.gr7.us-east-1.eks.amazonaws.com') {
                sh "kubectl get pods -n webapps"
                sh "kubectl get svc -n webapps"
            } } }
    }
}
