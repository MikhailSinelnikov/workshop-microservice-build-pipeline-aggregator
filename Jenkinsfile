def srcVersion = ''
def branchQual = ''
def dockerImageBase = 'nexus.build.svc.cluster.local:5000/workshop/aggregator'
def dockerImageTag = ''
def dockerImageFullTag = ''

timestamps {
    podTemplate(
        label: 'demoapp',
        containers: [
            containerTemplate(name: 'docker', image: 'docker', ttyEnabled: true, command: 'cat', args: '-v'),
            containerTemplate(name: 'helm', image: 'dtzar/helm-kubectl:2.9.1', ttyEnabled: true, command: 'cat')
        ],
        volumes: [
            // this is necessary for docker build to work
            hostPathVolume(hostPath: '/run/docker.sock', mountPath: '/run/docker.sock')
        ]) {

        node ('demoapp') {
            try {
                stage('prepare') {
                    // checkout the source code
                    checkout scm

                    // get current version from the source code
                    srcVersion = sh(returnStdout: true, script: 'echo -n "$(cat version)"').trim().toLowerCase()

                    // get git commit hash
                    def gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

                    // git branch name
                    def gitBranch = sh(returnStdout: true, script: 'git rev-parse --abbrev-ref HEAD').trim()

                    // branch qualifier for use in tag
                    branchQual = gitBranch.replaceAll('[^a-zA-Z0-9]', '')
                    if (branchQual.length() > 32) {
                        branchQual = branchQual.substring(branchQual.length() - 32, branchQual.length())
                    }

                    if (branchQual != 'master') {
                        dockerImageTag = "${srcVersion}-${branchQual}.${BUILD_NUMBER}"
                    } else {
                        dockerImageTag = "${srcVersion}"
                    }

                    dockerImageFullTag = "${dockerImageBase}:${dockerImageTag}"

                    echo "srcVersion: ${srcVersion}"

                    // Gen necessary tools
                    sh """
                        wget https://github.com/istio/istio/releases/download/1.0.3/istio-1.0.3-linux.tar.gz
                        tar zxvf istio-1.0.3-linux.tar.gz
                        mkdir -p ~/.local/bin
                        cp istio-1.0.3/bin/istioctl ~/.local/bin/
                        chmod a+x ~/.local/bin/istioctl
                    """
                }

                stage('build') {
                    container ('docker') {
                        sh "docker build --force-rm --tag '${dockerImageFullTag}' ."
                    }
                }

                stage('unit-test') {
                    container ('docker') {
                        echo 'e.g. run mocha unit tests'
                    }
                }

                stage('publish') {
                    container ('docker') {
                        withCredentials([
                                usernamePassword(credentialsId: 'docker-registry', usernameVariable: 'dockerUsername', passwordVariable: 'dockerPassword'),
                        ]) {
                            sh "docker login ${dockerImageFullTag} -u '${dockerUsername}' -p '${dockerPassword}'"
                            sh "docker push ${dockerImageFullTag}"
                        }
                    }
                }

                stage('deploy') {
                    container ('helm') {
                        withCredentials([
                                usernamePassword(credentialsId: 'docker-registry', usernameVariable: 'dockerUsername', passwordVariable: 'dockerPassword'),
                                file(credentialsId: 'k8s-config', variable: 'k8sConfig')
                        ]) {
                            def namespace = branchQual

                            // create namespace for the branch if it does not exist
                            sh "kubectl --kubeconfig='${k8sConfig}' create namespace '${namespace}' || true"

                            // make sure that docker registry secret is configured in the namespace
                            sh """
                                kubectl delete secrets registry-aggregator --kubeconfig='${k8sConfig}' --namespace '${namespace}' || true

                                kubectl create secret docker-registry registry-aggregator \\
                                    --kubeconfig='${k8sConfig}' \\
                                    --namespace '${namespace}' \\
                                    --docker-server='${dockerImageBase}' \\
                                    --docker-username='${dockerUsername}' \\
                                    --docker-password='${dockerPassword}' \\
                                    --docker-email=email@example.com
                            """

                            sh """
                                kubectl --kubeconfig='${k8sConfig}' --namespace '${namespace}' \\
                                    apply -f kubernetes/service.yaml

                                cat kubernetes/deployment.yaml | \\
                                    sed -e 's^{{image}}^${dockerImageFullTag}^g' | \\
                                    sed -e 's^{{version}}^jenkins^g' | \\
                                    ~/.local/bin/istioctl --kubeconfig='${k8sConfig}' kube-inject -f - | \\
                                    kubectl --kubeconfig='${k8sConfig}' --namespace '${namespace}' delete -f - || true

                                cat kubernetes/deployment.yaml | \\
                                    sed -e 's^{{image}}^${dockerImageFullTag}^g' | \\
                                    sed -e 's^{{version}}^jenkins^g' | \\
                                    ~/.local/bin/istioctl --kubeconfig='${k8sConfig}' kube-inject -f - | \\
                                    kubectl --kubeconfig='${k8sConfig}' --namespace '${namespace}' apply -f -
                            """
                        }
                    }
                }

                stage('integration-test') {
                    container ('docker') {
                        echo 'e.g. run integration tests'
                    }
                }
            } finally {
                stage('cleanup') {
                    container ('docker') {
                        sh "docker rmi --force ${dockerImageFullTag} || true"
                        sh "docker logout ${dockerImageFullTag} || true"
                    }
                }
            }
        }
    }
}
